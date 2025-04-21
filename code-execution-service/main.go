package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/IMMZEK/AggieCode/code-execution-service/executor"
)

// ExecuteRequest defines the structure for code execution requests.
type ExecuteRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Stdin    string `json:"stdin,omitempty"` // Optional standard input
}

// ExecuteResponse defines the structure for code execution responses.
type ExecuteResponse struct {
	Stdout          string `json:"stdout"`
	Stderr          string `json:"stderr"`
	Error           string `json:"error,omitempty"` // For execution or setup errors
	ExecutionTimeMs int64  `json:"execution_time_ms"`
}

// Global executor instance - now using the interface type
var codeExecutor executor.CodeExecutionService

func main() {
	// Basic structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize the executor with image prefix from environment variable
	// Default to empty string (use local images without prefix)
	imagePrefix := os.Getenv("IMAGE_PREFIX")
	var err error

	// Create the concrete executor implementation
	dockerExecutor, err := executor.NewExecutor(imagePrefix)
	if err != nil {
		slog.Error("Failed to initialize code executor", "error", err)
		os.Exit(1)
	}

	// Assign to the interface variable
	codeExecutor = dockerExecutor

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
		WriteTimeout: 30 * time.Second, // Longer timeout for code execution
		IdleTimeout:  15 * time.Second,
	}

	slog.Info("Starting Code Execution Service", "address", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

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

	slog.Info("Received execution request", "language", req.Language, "code_length", len(req.Code))

	// Create a context with timeout for the execution
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	// Create an execution request for the executor
	execReq := executor.ExecutionRequest{
		Language: req.Language,
		Code:     req.Code,
		Stdin:    req.Stdin,
		Timeout:  15 * time.Second,
	}

	// Execute the code
	result, err := codeExecutor.Execute(ctx, execReq)
	if err != nil {
		slog.Error("Code execution failed", "error", err, "language", req.Language)
		// Provide a sanitized error message to the client
		resp := ExecuteResponse{
			Error:           fmt.Sprintf("Execution error: %v", err),
			ExecutionTimeMs: time.Since(startTime).Milliseconds(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Map the execution result to our response
	resp := ExecuteResponse{
		Stdout:          result.Stdout,
		Stderr:          result.Stderr,
		Error:           result.Error,
		ExecutionTimeMs: result.ExecTimeMs,
	}

	slog.Info("Code execution completed",
		"language", req.Language,
		"execution_time_ms", resp.ExecutionTimeMs,
		"has_error", resp.Error != "")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		// Log the error, but the header is likely already sent.
		slog.Error("Failed to encode response", "error", err)
	}
}

// healthCheckHandler returns a basic health check response
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
