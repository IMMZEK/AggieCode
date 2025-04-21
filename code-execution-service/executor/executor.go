// Package executor provides functionality for securely executing code using Docker containers.
package executor

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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

// CodeExecutor handles code execution in Docker containers
type CodeExecutor struct {
	dockerClient *client.Client
	imagePrefix  string // Prefix for Docker images, e.g. "aggiecode/"
}

// NewExecutor creates a new CodeExecutor instance
func NewExecutor(imagePrefix string) (*CodeExecutor, error) {
	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &CodeExecutor{
		dockerClient: cli,
		imagePrefix:  imagePrefix,
	}, nil
}

// Execute runs the provided code in a Docker container
func (e *CodeExecutor) Execute(ctx context.Context, req ExecutionRequest) (ExecutionResult, error) {
	startTime := time.Now()
	result := ExecutionResult{}

	// Check if language is supported
	imageName, supported := SupportedLanguages[req.Language]
	if !supported {
		return result, fmt.Errorf("unsupported language: %s", req.Language)
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
	containerID, err := e.createAndStartContainer(ctx, imageName, tempDir, filename, stdinFile)
	if err != nil {
		return result, fmt.Errorf("container execution failed: %w", err)
	}
	defer e.cleanupContainer(ctx, containerID)

	// Wait for the container to finish with timeout
	statusCh, errCh := e.dockerClient.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	var statusCode int64
	select {
	case err := <-errCh:
		if err != nil {
			return result, fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-statusCh:
		statusCode = status.StatusCode
	case <-time.After(req.Timeout):
		return result, fmt.Errorf("execution timed out after %v", req.Timeout)
	}

	// Get container logs
	stdout, stderr, err := e.getContainerLogs(ctx, containerID)
	if err != nil {
		return result, fmt.Errorf("failed to get container logs: %w", err)
	}

	// Process the results
	result.Stdout = stdout
	result.Stderr = stderr
	result.ExecTimeMs = time.Since(startTime).Milliseconds()

	// Handle non-zero exit codes
	if statusCode != 0 {
		result.Error = fmt.Sprintf("Process exited with code %d", statusCode)
	}

	return result, nil
}

// writeCodeFile writes the code to the appropriate file based on language
func (e *CodeExecutor) writeCodeFile(dir, language, code string) (string, error) {
	var filename string

	switch language {
	case "python":
		filename = filepath.Join(dir, "main.py")
	case "javascript":
		filename = filepath.Join(dir, "main.js")
	case "cpp":
		filename = filepath.Join(dir, "main.cpp")
	case "java":
		// For Java, we need to use a class name of "Main"
		filename = filepath.Join(dir, "Main.java")
		// Check if the code contains a public class that's not named Main
		// This is a simplified check and may not catch all cases
		if !bytes.Contains([]byte(code), []byte("class Main")) {
			// Wrap the code in a Main class if it doesn't define one
			code = fmt.Sprintf("public class Main {\n    %s\n}", code)
		}
	case "go":
		filename = filepath.Join(dir, "main.go")
	default:
		return "", fmt.Errorf("unsupported language for file creation: %s", language)
	}

	return filename, ioutil.WriteFile(filename, []byte(code), 0644)
}

// createAndStartContainer creates and starts a Docker container for code execution
func (e *CodeExecutor) createAndStartContainer(ctx context.Context, imageName, tempDir, codeFile, stdinFile string) (string, error) {
	// Prepare mount for code directory
	mounts := []mount.Mount{
		{
			Type:     mount.TypeBind,
			Source:   tempDir,
			Target:   "/code",
			ReadOnly: true, // Mount as read-only for security
		},
	}

	// Create container configuration
	config := &container.Config{
		Image:      imageName,
		Cmd:        e.buildCommand(filepath.Base(codeFile), filepath.Base(stdinFile)),
		Tty:        false,
		WorkingDir: "/code", // Set working directory
	}

	// Create host configuration with security settings
	hostConfig := &container.HostConfig{
		Mounts:      mounts,
		NetworkMode: container.NetworkMode("none"), // Disable networking
		Resources: container.Resources{
			Memory:    256 * 1024 * 1024, // 256MB
			CPUShares: 512,
		},
		ReadonlyRootfs: true, // Read-only filesystem for security
	}

	// Create the container
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
func (e *CodeExecutor) buildCommand(codeFile, stdinFile string) []string {
	// Extract language from filename
	ext := filepath.Ext(codeFile)
	var cmd []string

	switch ext {
	case ".py":
		cmd = []string{"python3", codeFile}
	case ".js":
		cmd = []string{"node", codeFile}
	case ".cpp":
		// For C++, we need to compile first, then run
		// This is handled in the container's entrypoint script
		cmd = []string{codeFile}
	case ".java":
		// For Java, we need to compile first, then run
		// This is handled in the container's entrypoint script
		cmd = []string{codeFile}
	case ".go":
		// For Go, we use 'go run'
		cmd = []string{"go", "run", codeFile}
	}

	// If stdin file is provided, add it to the command
	if stdinFile != "" {
		cmd = append(cmd, "<", stdinFile)
	}

	return cmd
}

// getContainerLogs fetches stdout and stderr from a container
func (e *CodeExecutor) getContainerLogs(ctx context.Context, containerID string) (string, string, error) {
	// Get container logs
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}
	reader, err := e.dockerClient.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", "", err
	}
	defer reader.Close()

	// Docker multiplexes stdout and stderr, so we need to separate them
	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, reader)
	if err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}

// cleanupContainer stops and removes a container
func (e *CodeExecutor) cleanupContainer(ctx context.Context, containerID string) {
	// Stop the container (timeout after 5 seconds)
	stopTimeout := 5
	stopOptions := container.StopOptions{
		Timeout: &stopTimeout,
	}
	e.dockerClient.ContainerStop(ctx, containerID, stopOptions)

	// Remove the container
	removeOptions := container.RemoveOptions{
		Force: true,
	}
	e.dockerClient.ContainerRemove(ctx, containerID, removeOptions)
}
