package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/IMMZEK/AggieCode/code-execution-service/executor"
)

// ExecuteRequest defines the structure for code execution requests.
type ExecuteRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Stdin    string `json:"stdin,omitempty"`   // Optional standard input
	Timeout  int    `json:"timeout,omitempty"` // Optional timeout in seconds
}

// ExecuteResponse defines the structure for code execution responses.
type ExecuteResponse struct {
	Stdout          string `json:"stdout"`
	Stderr          string `json:"stderr"`
	Error           string `json:"error,omitempty"`      // For execution or setup errors
	ErrorType       string `json:"error_type,omitempty"` // Type of error (timeout, memory_limit, etc.)
	ExecutionTimeMs int64  `json:"execution_time_ms"`
}

// Global executor service
var codeExecutor executor.CodeExecutionService

func main() {
	// Basic structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize the executor with configuration from environment variables
	executorConfig := executor.ExecutorConfig{
		ImagePrefix:     os.Getenv("IMAGE_PREFIX"),
		ConcurrentLimit: getConcurrentLimitFromEnv(),
		DefaultTimeout:  getDefaultTimeoutFromEnv(),
	}

	var err error
	codeExecutor, err = executor.NewExecutorWithConfig(executorConfig)
	if err != nil {
		slog.Error("Failed to initialize code executor", "error", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default port for the CES
		slog.Info("Defaulting to port", "port", port)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/execute", executeHandler)
	mux.HandleFunc("/health", healthCheckHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 35 * time.Second, // Longer timeout to account for maximum execution time
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("Starting Code Execution Service",
		"address", server.Addr,
		"concurrent_limit", executorConfig.ConcurrentLimit,
		"default_timeout", executorConfig.DefaultTimeout)

	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

// Helper function to get concurrent limit from environment
func getConcurrentLimitFromEnv() int {
	limitStr := os.Getenv("CONCURRENT_LIMIT")
	if limitStr == "" {
		return executor.DefaultConcurrentLimit
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		slog.Warn("Invalid CONCURRENT_LIMIT value, using default",
			"value", limitStr,
			"default", executor.DefaultConcurrentLimit)
		return executor.DefaultConcurrentLimit
	}

	return limit
}

// Helper function to get default timeout from environment
func getDefaultTimeoutFromEnv() time.Duration {
	timeoutStr := os.Getenv("DEFAULT_TIMEOUT")
	if timeoutStr == "" {
		return executor.DefaultExecutionTime
	}

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil || timeout <= 0 || timeout > 30 {
		slog.Warn("Invalid DEFAULT_TIMEOUT value, using default",
			"value", timeoutStr,
			"default", executor.DefaultExecutionTime)
		return executor.DefaultExecutionTime
	}

	return time.Duration(timeout) * time.Second
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var req ExecuteRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent unexpected fields
	err := decoder.Decode(&req)
	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Basic Validation
	if req.Code == "" {
		http.Error(w, "Missing 'code' field in request", http.StatusBadRequest)
		return
	}
	if req.Language == "" {
		http.Error(w, "Missing 'language' field in request", http.StatusBadRequest)
		return
	}

	slog.Info("Received execution request",
		"language", req.Language,
		"code_length", len(req.Code),
		"timeout", req.Timeout)

	// Convert timeout to duration
	var timeout time.Duration
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
		// Cap at the maximum allowed timeout
		if timeout > executor.MaxExecutionTime {
			timeout = executor.MaxExecutionTime
		}
	}

	// Create an execution request for the executor
	execReq := executor.ExecutionRequest{
		Language: req.Language,
		Code:     req.Code,
		Stdin:    req.Stdin,
		Timeout:  timeout,
	}

	// Create a context with request-scoped cancellation
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Handle client disconnection
	go func() {
		<-r.Context().Done() // Wait for client to disconnect
		cancel()             // Cancel our execution context
	}()

	// Execute the code
	result, err := codeExecutor.Execute(ctx, execReq)

	// Create the response
	resp := ExecuteResponse{
		Stdout:          result.Stdout,
		Stderr:          result.Stderr,
		Error:           result.Error,
		ExecutionTimeMs: result.ExecTimeMs,
	}

	// Handle specific error types
	statusCode := http.StatusOK
	if err != nil {
		slog.Error("Code execution failed", "error", err, "language", req.Language)

		// Check if this is a specific execution error type
		if execErr, ok := err.(executor.ExecutionError); ok {
			resp.ErrorType = execErr.Type
			resp.Error = execErr.Message

			// Map error types to appropriate status codes
			switch execErr.Type {
			case "timeout":
				statusCode = http.StatusRequestTimeout
			case "memory_limit":
				statusCode = http.StatusRequestEntityTooLarge // Using StatusRequestEntityTooLarge instead of StatusPayloadTooLarge
			case "limit_exceeded":
				statusCode = http.StatusTooManyRequests
			case "unsupported_language":
				statusCode = http.StatusBadRequest
			default:
				statusCode = http.StatusInternalServerError
			}
		} else {
			// Generic error
			resp.Error = fmt.Sprintf("Execution error: %v", err)
			statusCode = http.StatusInternalServerError
		}
	}

	slog.Info("Code execution completed",
		"language", req.Language,
		"execution_time_ms", resp.ExecutionTimeMs,
		"error_type", resp.ErrorType,
		"has_error", resp.Error != "")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// healthCheckHandler returns a basic health check response
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
