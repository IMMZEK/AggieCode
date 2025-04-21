package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

	// Setup: Create a valid request body
	reqBody := ExecuteRequest{
		Language: "python",
		Code:     "print('hello')",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Create a request to pass to our handler
	req, err := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the status code
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var execResp ExecuteResponse
	err = json.NewDecoder(resp.Body).Decode(&execResp)
	if err != nil {
		t.Fatalf("Could not unmarshal response body: %v", err)
	}

	// Check specific fields
	expectedStdout := "Hello from Python!\n"
	if execResp.Stdout != expectedStdout {
		t.Errorf("handler returned unexpected stdout: got %q want %q", execResp.Stdout, expectedStdout)
	}
	if execResp.Stderr != "" {
		t.Errorf("handler returned unexpected stderr: got %q want %q", execResp.Stderr, "")
	}
	if execResp.Error != "" {
		t.Errorf("handler returned unexpected error: got %q want %q", execResp.Error, "")
	}
	if execResp.ExecutionTimeMs <= 0 {
		t.Errorf("handler returned non-positive execution time: got %v", execResp.ExecutionTimeMs)
	}
}

func TestExecuteHandler_WithStdin(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Setup: Create a request with stdin
	reqBody := ExecuteRequest{
		Language: "python",
		Code:     "print('Input provided')",
		Stdin:    "test input",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Create a request to our test server
	req, err := http.NewRequest("POST", server.URL+"/api/execute", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the status code
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var execResp ExecuteResponse
	err = json.NewDecoder(resp.Body).Decode(&execResp)
	if err != nil {
		t.Fatalf("Could not unmarshal response body: %v", err)
	}

	// Verify that stdin was processed
	if !strings.Contains(execResp.Stdout, "Input: test input") {
		t.Errorf("handler did not process stdin properly: got %q", execResp.Stdout)
	}
}

func TestExecuteHandler_InvalidMethod(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/api/execute", nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestExecuteHandler_InvalidContentType(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, err := http.NewRequest("POST", server.URL+"/api/execute", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusUnsupportedMediaType {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnsupportedMediaType)
	}
}

func TestExecuteHandler_InvalidJson(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, err := http.NewRequest("POST", server.URL+"/api/execute", strings.NewReader("{"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestExecuteHandler_MissingFields(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	testCases := []struct {
		name        string
		body        string
		expectedMsg string
	}{
		{
			name:        "Missing Code",
			body:        `{"language": "python"}`,
			expectedMsg: "Missing 'code' field in request",
		},
		{
			name:        "Missing Language",
			body:        `{"code": "print(1)"}`,
			expectedMsg: "Missing 'language' field in request",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", server.URL+"/api/execute", strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: time.Second * 5}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if status := resp.StatusCode; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
			}

			body := new(bytes.Buffer)
			body.ReadFrom(resp.Body)
			if !strings.Contains(body.String(), tc.expectedMsg) {
				t.Errorf("handler returned wrong error message: got %q want containing %q", body.String(), tc.expectedMsg)
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("health check handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var healthResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&healthResp)
	if err != nil {
		t.Fatalf("Could not unmarshal response body: %v", err)
	}

	if status, ok := healthResp["status"]; !ok || status != "ok" {
		t.Errorf("health check handler returned incorrect status: got %v want %v", status, "ok")
	}
}
