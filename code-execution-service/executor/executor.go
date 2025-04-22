// Package executor provides functionality for securely executing code using Docker containers.
package executor

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// Default resource limits
const (
	DefaultMemoryLimit     = 256 * 1024 * 1024 // 256MB
	DefaultCPULimit        = 1.0               // 1 CPU core
	DefaultExecutionTime   = 10 * time.Second  // 10 seconds
	MaxExecutionTime       = 30 * time.Second  // 30 seconds max
	DefaultPidsLimit       = int64(50)         // Max number of processes
	DefaultConcurrentLimit = 10                // Max concurrent executions
	DefaultNetworkPolicy   = "none"            // Disable networking
)

// CodeExecutionService defines the interface for code execution
type CodeExecutionService interface {
	Execute(ctx context.Context, req ExecutionRequest) (ExecutionResult, error)
}

// SupportedLanguages is a map of languages that the executor supports
var SupportedLanguages = map[string]string{
	"python":     "python-executor",
	"javascript": "js-executor",
	"cpp":        "cpp-executor",
	"java":       "java-executor",
	"go":         "go-executor",
}

// LanguageCompilers defines which languages need compilation
var LanguageCompilers = map[string]bool{
	"cpp":  true,
	"java": true,
	"go":   true, // go run does compilation internally
}

// ExecutionRequest encapsulates all information needed to execute code
type ExecutionRequest struct {
	Language string
	Code     string
	Stdin    string
	Timeout  time.Duration // Maximum execution time
}

// ExecutionResult contains the results of code execution
type ExecutionResult struct {
	Stdout     string
	Stderr     string
	Error      string
	ExecTimeMs int64
}

// ExecutionError represents specific error types that can occur during execution
type ExecutionError struct {
	Type    string // "timeout", "memory", "runtime", "compilation", etc.
	Message string
}

func (e ExecutionError) Error() string {
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

// MockExecutor provides a fallback execution mode when Docker is not available
type MockExecutor struct{}

// Execute simulates code execution without Docker
func (m *MockExecutor) Execute(_ context.Context, req ExecutionRequest) (ExecutionResult, error) {
	startTime := time.Now()
	result := ExecutionResult{}

	// Generate predictable output based on the language and code
	switch req.Language {
	case "python":
		if strings.Contains(req.Code, "print") {
			result.Stdout = "Python output: " + extractPrintContent(req.Code, "python")
			if strings.Contains(req.Code, "input") && req.Stdin != "" {
				result.Stdout += "\nInput was: " + req.Stdin
			}
		} else if strings.Contains(req.Code, "error") || strings.Contains(req.Code, "raise") {
			result.Stderr = "Python error: Simulated exception"
			result.Error = "Process exited with code 1"
		}
	case "javascript":
		if strings.Contains(req.Code, "console.log") {
			result.Stdout = "JavaScript output: " + extractPrintContent(req.Code, "javascript")
		}
	case "cpp":
		if strings.Contains(req.Code, "cout") {
			result.Stdout = "C++ output: " + extractPrintContent(req.Code, "cpp")
		}
	case "java":
		if strings.Contains(req.Code, "System.out.println") {
			result.Stdout = "Java output: " + extractPrintContent(req.Code, "java")
		}
	case "go":
		if strings.Contains(req.Code, "fmt.Println") {
			result.Stdout = "Go output: " + extractPrintContent(req.Code, "go")
		}
	default:
		return result, fmt.Errorf("unsupported language in mock mode: %s", req.Language)
	}

	// Calculate execution time
	result.ExecTimeMs = time.Since(startTime).Milliseconds()
	return result, nil
}

// Helper function to extract content from print statements for the mock executor
func extractPrintContent(code, language string) string {
	var printStart, printEnd string
	switch language {
	case "python":
		printStart = "print("
		printEnd = ")"
	case "javascript":
		printStart = "console.log("
		printEnd = ")"
	case "cpp":
		printStart = "cout <<"
		printEnd = ";"
	case "java":
		printStart = "System.out.println("
		printEnd = ")"
	case "go":
		printStart = "fmt.Println("
		printEnd = ")"
	}

	if idx := strings.Index(code, printStart); idx >= 0 {
		code = code[idx+len(printStart):]
		if idx = strings.Index(code, printEnd); idx >= 0 {
			content := code[:idx]
			// Clean up quotes if present
			content = strings.Trim(content, "'\"")
			return content
		}
	}
	return "[Content could not be extracted]"
}

// CodeExecutor handles code execution in Docker containers
type CodeExecutor struct {
	dockerClient       *client.Client
	imagePrefix        string               // Prefix for Docker images, e.g. "aggiecode/"
	fallbackMode       bool                 // Use fallback mode when Docker is not available
	mockExecutor       CodeExecutionService // Mock executor for fallback mode
	concurrentLimit    int                  // Maximum number of concurrent executions
	executionSemaphore *chan struct{}       // Semaphore to limit concurrent executions
	executionLock      sync.Mutex           // Lock to protect concurrent access to the semaphore
}

// ExecutorConfig provides configuration options for the CodeExecutor
type ExecutorConfig struct {
	ImagePrefix     string        // Prefix for Docker images
	ConcurrentLimit int           // Maximum number of concurrent executions
	DefaultTimeout  time.Duration // Default timeout for code execution
}

// NewExecutor creates a new CodeExecutor instance with default configuration
func NewExecutor(imagePrefix string) (*CodeExecutor, error) {
	return NewExecutorWithConfig(ExecutorConfig{
		ImagePrefix:     imagePrefix,
		ConcurrentLimit: DefaultConcurrentLimit,
		DefaultTimeout:  DefaultExecutionTime,
	})
}

// NewExecutorWithConfig creates a new CodeExecutor instance with the given configuration
func NewExecutorWithConfig(config ExecutorConfig) (*CodeExecutor, error) {
	// Set default values if not provided
	if config.ConcurrentLimit <= 0 {
		config.ConcurrentLimit = DefaultConcurrentLimit
	}

	if config.DefaultTimeout <= 0 {
		config.DefaultTimeout = DefaultExecutionTime
	}

	// Try to create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	// Create semaphore for limiting concurrent executions
	semaphore := make(chan struct{}, config.ConcurrentLimit)

	if err != nil {
		// Docker client creation failed, use fallback mode
		fmt.Println("WARNING: Could not create Docker client, using fallback mode")
		return &CodeExecutor{
			dockerClient:       nil,
			imagePrefix:        config.ImagePrefix,
			fallbackMode:       true,
			mockExecutor:       &MockExecutor{},
			concurrentLimit:    config.ConcurrentLimit,
			executionSemaphore: &semaphore,
		}, nil
	}

	// Test Docker connection
	_, err = cli.Ping(context.Background())
	if err != nil {
		// Docker daemon is not running, use fallback mode
		fmt.Println("WARNING: Docker daemon is not running, using fallback mode")
		return &CodeExecutor{
			dockerClient:       nil,
			imagePrefix:        config.ImagePrefix,
			fallbackMode:       true,
			mockExecutor:       &MockExecutor{},
			concurrentLimit:    config.ConcurrentLimit,
			executionSemaphore: &semaphore,
		}, nil
	}

	return &CodeExecutor{
		dockerClient:       cli,
		imagePrefix:        config.ImagePrefix,
		fallbackMode:       false,
		mockExecutor:       nil,
		concurrentLimit:    config.ConcurrentLimit,
		executionSemaphore: &semaphore,
	}, nil
}

// Execute runs the provided code in a Docker container or falls back to mock execution
func (e *CodeExecutor) Execute(ctx context.Context, req ExecutionRequest) (ExecutionResult, error) {
	// If in fallback mode, use the mock executor
	if e.fallbackMode {
		return e.mockExecutor.Execute(ctx, req)
	}

	// Validate the timeout
	if req.Timeout <= 0 {
		req.Timeout = DefaultExecutionTime
	} else if req.Timeout > MaxExecutionTime {
		req.Timeout = MaxExecutionTime
	}

	// Acquire semaphore to limit concurrent executions
	e.executionLock.Lock()
	select {
	case *e.executionSemaphore <- struct{}{}:
		// Successfully acquired semaphore
		e.executionLock.Unlock()
		defer func() {
			// Release semaphore when done
			<-*e.executionSemaphore
		}()
	case <-ctx.Done():
		// Context canceled while waiting for semaphore
		e.executionLock.Unlock()
		return ExecutionResult{}, ExecutionError{
			Type:    "timeout",
			Message: "execution queue is full, try again later",
		}
	default:
		// Semaphore channel is full
		e.executionLock.Unlock()
		return ExecutionResult{}, ExecutionError{
			Type:    "limit_exceeded",
			Message: "too many concurrent executions, try again later",
		}
	}

	// Create a timeout context for the execution
	execCtx, cancel := context.WithTimeout(ctx, req.Timeout)
	defer cancel()

	startTime := time.Now()
	result := ExecutionResult{}

	// Check if language is supported
	imageName, supported := SupportedLanguages[req.Language]
	if !supported {
		return result, ExecutionError{
			Type:    "unsupported_language",
			Message: fmt.Sprintf("unsupported language: %s", req.Language),
		}
	}

	// Add image prefix if provided
	if e.imagePrefix != "" {
		imageName = e.imagePrefix + imageName
	}

	// Create temporary directory for code files
	tempDir, err := ioutil.TempDir("", fmt.Sprintf("aggiecode-%s-", req.Language))
	if err != nil {
		return result, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up when done

	// Write code to appropriate file based on language
	filename, err := e.writeCodeFile(tempDir, req.Language, req.Code)
	if err != nil {
		return result, fmt.Errorf("failed to write code file: %w", err)
	}

	// Write stdin to file if provided
	var stdinFile string
	if req.Stdin != "" {
		stdinFile = filepath.Join(tempDir, "input.txt")
		if err := ioutil.WriteFile(stdinFile, []byte(req.Stdin), 0644); err != nil {
			return result, fmt.Errorf("failed to write stdin file: %w", err)
		}
	}

	// Create and run the container
	containerID, err := e.createAndStartContainer(execCtx, imageName, tempDir, filename, stdinFile, req.Language)
	if err != nil {
		if execCtx.Err() == context.DeadlineExceeded {
			return result, ExecutionError{
				Type:    "timeout",
				Message: fmt.Sprintf("container creation timed out after %v", req.Timeout),
			}
		}
		return result, fmt.Errorf("container execution failed: %w", err)
	}
	defer e.cleanupContainer(context.Background(), containerID)

	// Wait for the container to finish with timeout
	statusCh, errCh := e.dockerClient.ContainerWait(execCtx, containerID, container.WaitConditionNotRunning)
	var statusCode int64

	select {
	case err := <-errCh:
		if execCtx.Err() == context.DeadlineExceeded {
			// Context deadline exceeded - execution timed out
			return result, ExecutionError{
				Type:    "timeout",
				Message: fmt.Sprintf("execution timed out after %v", req.Timeout),
			}
		}
		if err != nil {
			return result, fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-statusCh:
		statusCode = status.StatusCode
	case <-execCtx.Done():
		// Context canceled or timed out
		if execCtx.Err() == context.DeadlineExceeded {
			return result, ExecutionError{
				Type:    "timeout",
				Message: fmt.Sprintf("execution timed out after %v", req.Timeout),
			}
		}
		return result, fmt.Errorf("execution canceled: %v", execCtx.Err())
	}

	// Check if the container was killed due to OOM (out of memory)
	containerJSON, err := e.dockerClient.ContainerInspect(context.Background(), containerID)
	if err == nil && containerJSON.State != nil && containerJSON.State.OOMKilled {
		return result, ExecutionError{
			Type:    "memory_limit",
			Message: "execution exceeded memory limit",
		}
	}

	// Get container logs
	stdout, stderr, err := e.getContainerLogs(context.Background(), containerID)
	if err != nil {
		return result, fmt.Errorf("failed to get container logs: %w", err)
	}

	// Process the results
	result.Stdout = stdout
	result.Stderr = stderr
	result.ExecTimeMs = time.Since(startTime).Milliseconds()

	// Handle non-zero exit codes
	if statusCode != 0 {
		// Check if this is a compilation error (for compiled languages)
		if needsCompilation, ok := LanguageCompilers[req.Language]; ok && needsCompilation {
			if strings.Contains(stderr, "error") || strings.Contains(stderr, "Error") {
				result.Error = fmt.Sprintf("Compilation error (exit code %d)", statusCode)
			} else {
				result.Error = fmt.Sprintf("Runtime error (exit code %d)", statusCode)
			}
		} else {
			result.Error = fmt.Sprintf("Process exited with code %d", statusCode)
		}
	}

	return result, nil
}

// createAndStartContainer creates and starts a Docker container for code execution
func (e *CodeExecutor) createAndStartContainer(ctx context.Context, imageName, tempDir, codeFile, stdinFile, language string) (string, error) {
	// Prepare mount for code directory
	mounts := []mount.Mount{
		{
			Type:     mount.TypeBind,
			Source:   tempDir,
			Target:   "/code",
			ReadOnly: false, // Enable writing for compilation outputs
		},
	}

	// Set up command based on language
	cmd := e.buildCommand(filepath.Base(codeFile), filepath.Base(stdinFile), language)

	// Create container configuration
	config := &container.Config{
		Image:      imageName,
		Cmd:        cmd,
		Tty:        false,
		WorkingDir: "/code", // Set working directory
	}

	// Convert CPU limit from core count to nano-CPUs
	nanoCPUs := int64(DefaultCPULimit * 1e9)

	// Store our pids limit
	pidsLimit := DefaultPidsLimit

	// Create host configuration with security settings
	hostConfig := &container.HostConfig{
		Mounts:         mounts,
		NetworkMode:    container.NetworkMode(DefaultNetworkPolicy), // Disable networking
		ReadonlyRootfs: true,                                        // Read-only filesystem for security
		Resources: container.Resources{
			Memory:    DefaultMemoryLimit, // Memory limit
			NanoCPUs:  nanoCPUs,           // CPU limit
			PidsLimit: &pidsLimit,         // Process limit
		},
	}

	// Create the container with updated API
	resp, err := e.dockerClient.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}

	// Start the container
	if err := e.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return resp.ID, err // Return ID for cleanup, even though start failed
	}

	return resp.ID, nil
}

// buildCommand constructs the command to run based on the language and files
func (e *CodeExecutor) buildCommand(codeFile, stdinFile, language string) []string {
	var cmd strslice.StrSlice

	switch language {
	case "python":
		if stdinFile != "" {
			// Use shell to handle redirection
			cmd = strslice.StrSlice{"/bin/sh", "-c", fmt.Sprintf("python3 %s < %s", codeFile, stdinFile)}
		} else {
			cmd = strslice.StrSlice{"python3", codeFile}
		}
	case "javascript":
		if stdinFile != "" {
			cmd = strslice.StrSlice{"/bin/sh", "-c", fmt.Sprintf("node %s < %s", codeFile, stdinFile)}
		} else {
			cmd = strslice.StrSlice{"node", codeFile}
		}
	case "cpp":
		// For C++, we use the entrypoint script in the container
		if stdinFile != "" {
			cmd = strslice.StrSlice{"/bin/sh", "-c", fmt.Sprintf("./%s < %s", codeFile, stdinFile)}
		} else {
			cmd = strslice.StrSlice{codeFile}
		}
	case "java":
		// For Java, we use the entrypoint script in the container
		if stdinFile != "" {
			cmd = strslice.StrSlice{"/bin/sh", "-c", fmt.Sprintf("./%s < %s", codeFile, stdinFile)}
		} else {
			cmd = strslice.StrSlice{codeFile}
		}
	case "go":
		// For Go, we use 'go run'
		if stdinFile != "" {
			cmd = strslice.StrSlice{"/bin/sh", "-c", fmt.Sprintf("go run %s < %s", codeFile, stdinFile)}
		} else {
			cmd = strslice.StrSlice{"go", "run", codeFile}
		}
	default:
		// Default to executing the file directly
		if stdinFile != "" {
			cmd = strslice.StrSlice{"/bin/sh", "-c", fmt.Sprintf("%s < %s", codeFile, stdinFile)}
		} else {
			cmd = strslice.StrSlice{codeFile}
		}
	}

	return cmd
}

// writeCodeFile writes the code to an appropriate file based on the language
func (e *CodeExecutor) writeCodeFile(tempDir, language, code string) (string, error) {
	var filename string

	switch language {
	case "python":
		filename = filepath.Join(tempDir, "main.py")
	case "javascript":
		filename = filepath.Join(tempDir, "main.js")
	case "cpp":
		filename = filepath.Join(tempDir, "main.cpp")
	case "java":
		filename = filepath.Join(tempDir, "Main.java")
	case "go":
		filename = filepath.Join(tempDir, "main.go")
	default:
		// Default to a generic name
		filename = filepath.Join(tempDir, "main.txt")
	}

	return filename, ioutil.WriteFile(filename, []byte(code), 0644)
}

// getContainerLogs retrieves the stdout and stderr from the container
func (e *CodeExecutor) getContainerLogs(ctx context.Context, containerID string) (string, string, error) {
	// Get logs from the container
	reader, err := e.dockerClient.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", "", err
	}
	defer reader.Close()

	// Separate stdout and stderr
	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, reader)
	if err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}

// cleanupContainer removes the container after execution
func (e *CodeExecutor) cleanupContainer(ctx context.Context, containerID string) {
	// First try to stop the container gracefully
	stopTimeout := 1 // 1 second timeout for stopping
	err := e.dockerClient.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &stopTimeout})
	if err != nil {
		// If stopping fails, try to kill it
		e.dockerClient.ContainerKill(ctx, containerID, "SIGKILL")
	}

	// Remove the container
	e.dockerClient.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force: true,
	})
}
