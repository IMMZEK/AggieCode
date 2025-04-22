// Package main provides a command-line interface for testing the code execution service
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// ExecuteRequest defines the structure for code execution requests.
type ExecuteRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Stdin    string `json:"stdin,omitempty"`
	Timeout  int    `json:"timeout,omitempty"`
}

// ExecuteResponse defines the structure for code execution responses.
type ExecuteResponse struct {
	Stdout          string `json:"stdout"`
	Stderr          string `json:"stderr"`
	Error           string `json:"error,omitempty"`
	ErrorType       string `json:"error_type,omitempty"`
	ExecutionTimeMs int64  `json:"execution_time_ms"`
}

// Test cases to execute
var testCases = []struct {
	Name     string
	Language string
	Code     string
	Stdin    string
	Timeout  int
}{
	{
		Name:     "python",
		Language: "python",
		Code:     "print('Hello, World!')",
	},
	{
		Name:     "python with input",
		Language: "python",
		Code:     "print(input('Enter something: '))",
		Stdin:    "Test Input",
	},
	{
		Name:     "javascript",
		Language: "javascript",
		Code:     "console.log('Hello from JavaScript');",
	},
	{
		Name:     "go",
		Language: "go",
		Code: `package main

import "fmt"

func main() {
	fmt.Println("Hello from Go")
}`,
	},
	{
		Name:     "timeout example",
		Language: "python",
		Code: `import time
print("Starting long operation...")
time.sleep(15)  # This should trigger a timeout
print("Finished")`,
		Timeout: 5, // Set a custom timeout of 5 seconds
	},
	{
		Name:     "memory limit example",
		Language: "python",
		Code: `# This will attempt to allocate a large list to exceed memory limits
data = [0] * 1000000000  # Try to allocate a very large list
print("Allocated large memory block")`,
	},
	{
		Name:     "error case",
		Language: "python",
		Code: `# This will generate a syntax error
if True
    print("Missing colon")`,
	},
}

func main() {
	// Get service URL from command line or use default
	serviceURL := "http://localhost:8081/api/execute"
	if len(os.Args) > 1 {
		serviceURL = os.Args[1]
	}

	fmt.Printf("Testing Code Execution Service at %s\n\n", serviceURL)

	// Run each test case
	for i, tc := range testCases {
		fmt.Printf("Test Case %d: %s\n", i+1, tc.Name)
		fmt.Printf("Code: %s\n", tc.Code)

		if tc.Stdin != "" {
			fmt.Printf("Stdin: %s\n", tc.Stdin)
		}

		// Create the request
		req := ExecuteRequest{
			Language: tc.Language,
			Code:     tc.Code,
			Stdin:    tc.Stdin,
			Timeout:  tc.Timeout,
		}

		// Convert to JSON
		jsonData, err := json.Marshal(req)
		if err != nil {
			fmt.Printf("Error creating request JSON: %v\n\n", err)
			continue
		}

		// Send the request
		client := &http.Client{
			Timeout: 60 * time.Second, // Client timeout for the entire request
		}

		resp, err := client.Post(serviceURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error sending request: %v\n\n", err)
			continue
		}

		// Decode the response
		var result ExecuteResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			resp.Body.Close()
			fmt.Println()
			continue
		}
		resp.Body.Close()

		// Print results
		fmt.Printf("Status: %s\n", resp.Status)
		if result.ErrorType != "" {
			fmt.Printf("Error Type: %s\n", result.ErrorType)
		}
		if result.Error != "" {
			fmt.Printf("Error: %s\n", result.Error)
		}
		if result.Stdout != "" {
			fmt.Printf("Stdout: %s\n", result.Stdout)
		}
		if result.Stderr != "" {
			fmt.Printf("Stderr: %s\n", result.Stderr)
		}
		fmt.Printf("Execution Time: %d ms\n", result.ExecutionTimeMs)
		fmt.Println()
	}
}
