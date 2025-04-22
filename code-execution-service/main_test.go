package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IMMZEK/AggieCode/code-execution-service/executor"
)

// TestExecutor implements the CodeExecutionService interface for testing
type TestExecutor struct{}

// Execute mocks the code execution process without actually using Docker
func (m *TestExecutor) Execute(_ context.Context, req executor.ExecutionRequest) (executor.ExecutionResult, error) {
	// Simulate different responses based on language and code content
	result := executor.ExecutionResult{
		ExecTimeMs: 10, // Fixed execution time for predictability in tests
	}

	// Simulate a timeout error if requested in the code
	if strings.Contains(req.Code, "timeout") {
		return result, executor.ExecutionError{
			Type:    "timeout",
			Message: "execution timed out after 10s",
		}
	}

	// Simulate a memory limit error if requested in the code
	if strings.Contains(req.Code, "memory_limit") {
		return result, executor.ExecutionError{
			Type:    "memory_limit",
			Message: "execution exceeded memory limit",
		}
	}

	// Simulate a concurrency limit error if requested in the code
	if strings.Contains(req.Code, "too_many") {
		return result, executor.ExecutionError{
			Type:    "limit_exceeded",
			Message: "too many concurrent executions, try again later",
		}
	}

	// For testing, return a predictable output based on the language
	switch req.Language {
	case "python":
		if strings.Contains(req.Code, "print") {
			result.Stdout = "Hello from Python!\n"
		} else if strings.Contains(req.Code, "error") {
			result.Stderr = "Error: Python error occurred\n"
			result.Error = "Process exited with code 1"
		}
	case "javascript":
		if strings.Contains(req.Code, "console.log") {
			result.Stdout = "Hello from JavaScript!\n"
		}
	case "cpp":
		if strings.Contains(req.Code, "cout") {
			result.Stdout = "Hello from C++!\n"
		}
	case "java":
		if strings.Contains(req.Code, "System.out.println") {
			result.Stdout = "Hello from Java!\n"
		}
	case "go":
		if strings.Contains(req.Code, "fmt.Println") {
			result.Stdout = "Hello from Go!\n"
		}
	default:
		return result, executor.ExecutionError{
			Type:    "unsupported_language",
			Message: "unsupported language: " + req.Language,
		}
	}

	// If stdin was provided, echo it back
	if req.Stdin != "" {
		result.Stdout += "Input: " + req.Stdin + "\n"
	}

	return result, nil
}

func setupTestServer() *httptest.Server {
	// Use the mock executor instead of the real one for tests
	codeExecutor = &TestExecutor{}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/execute", executeHandler)
	mux.HandleFunc("/health", healthCheckHandler)
	return httptest.NewServer(mux)
}

func TestExecuteHandler_ValidRequest(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Valid Python code
	requestBody := `{"language":"python","code":"print('Hello, world!')"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var result ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !strings.Contains(result.Stdout, "Hello from Python!") {
		t.Errorf("Expected stdout to contain 'Hello from Python!', got: %s", result.Stdout)
	}
}

func TestExecuteHandler_WithStdin(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Python code with stdin
	requestBody := `{"language":"python","code":"input = input()\\nprint(input)","stdin":"Hello, input!"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var result ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !strings.Contains(result.Stdout, "Input: Hello, input!") {
		t.Errorf("Expected stdout to contain input echo, got: %s", result.Stdout)
	}
}

func TestExecuteHandler_InvalidMethod(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/api/execute", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status MethodNotAllowed, got %v", resp.Status)
	}
}

func TestExecuteHandler_InvalidContentType(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewBufferString("test"))
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("Expected status UnsupportedMediaType, got %v", resp.Status)
	}
}

func TestExecuteHandler_InvalidJson(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %v", resp.Status)
	}
}

func TestExecuteHandler_MissingFields(t *testing.T) {
	testCases := []struct {
		name        string
		requestBody string
	}{
		{
			name:        "Missing Code",
			requestBody: `{"language":"python"}`,
		},
		{
			name:        "Missing Language",
			requestBody: `{"code":"print('hello')"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := setupTestServer()
			defer server.Close()

			req, _ := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("Expected status BadRequest, got %v", resp.Status)
			}
		})
	}
}

func TestExecuteHandler_Timeout(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Request that will trigger a timeout error
	requestBody := `{"language":"python","code":"# timeout"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusRequestTimeout {
		t.Errorf("Expected status RequestTimeout, got %v", resp.Status)
	}

	var result ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.ErrorType != "timeout" {
		t.Errorf("Expected error_type 'timeout', got: %s", result.ErrorType)
	}
}

func TestExecuteHandler_MemoryLimit(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Request that will trigger a memory limit error
	requestBody := `{"language":"python","code":"# memory_limit"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusRequestEntityTooLarge { // Using StatusRequestEntityTooLarge instead of StatusPayloadTooLarge
		t.Errorf("Expected status RequestEntityTooLarge, got %v", resp.Status)
	}

	var result ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.ErrorType != "memory_limit" {
		t.Errorf("Expected error_type 'memory_limit', got: %s", result.ErrorType)
	}
}

func TestExecuteHandler_TooManyConcurrentExecutions(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Request that will trigger a too many concurrent executions error
	requestBody := `{"language":"python","code":"# too_many"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status TooManyRequests, got %v", resp.Status)
	}

	var result ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.ErrorType != "limit_exceeded" {
		t.Errorf("Expected error_type 'limit_exceeded', got: %s", result.ErrorType)
	}
}

func TestExecuteHandler_UnsupportedLanguage(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Request with an unsupported language
	requestBody := `{"language":"unsupported","code":"test"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %v", resp.Status)
	}

	var result ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.ErrorType != "unsupported_language" {
		t.Errorf("Expected error_type 'unsupported_language', got: %s", result.ErrorType)
	}
}

func TestHealthCheck(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/health", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got: %s", result["status"])
	}
}
