// Package main provides a simple test client for the Code Execution Service.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Request structure matching the API
type ExecuteRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Stdin    string `json:"stdin,omitempty"`
}

// Response structure matching the API
type ExecuteResponse struct {
	Stdout         string `json:"stdout"`
	Stderr         string `json:"stderr"`
	Error          string `json:"error,omitempty"`
	ExecutionTimeMs int64  `json:"execution_time_ms"`
}

func main() {
	// Define test cases
	testCases := []ExecuteRequest{
		{
			Language: "python",
			Code:     "print('Hello, World!')",
		},
		{
			Language: "python",
			Code:     "print(input('Enter something: '))",
			Stdin:    "Test Input",
		},
		{
			Language: "javascript",
			Code:     "console.log('Hello from JavaScript');",
		},
		{
			Language: "go",
			Code:     "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello from Go\")\n}",
		},
	}

	// API endpoint
	url := "http://localhost:8081/api/execute"
	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	fmt.Printf("Testing Code Execution Service at %s\n\n", url)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Run tests
	for i, tc := range testCases {
		fmt.Printf("Test Case %d: %s\n", i+1, tc.Language)
		fmt.Printf("Code: %s\n", tc.Code)
		if tc.Stdin != "" {
			fmt.Printf("Stdin: %s\n", tc.Stdin)
		}

		// Create request body
		reqBody, err := json.Marshal(tc)
		if err != nil {
			fmt.Printf("Error marshaling request: %v\n", err)
			continue
		}

		// Create request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		// Check response status
		fmt.Printf("Status: %s\n", resp.Status)

		// Parse response
		var execResp ExecuteResponse
		if err := json.NewDecoder(resp.Body).Decode(&execResp); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			continue
		}

		// Print response
		fmt.Printf("Stdout: %s\n", execResp.Stdout)
		if execResp.Stderr != "" {
			fmt.Printf("Stderr: %s\n", execResp.Stderr)
		}
		if execResp.Error != "" {
			fmt.Printf("Error: %s\n", execResp.Error)
		}
		fmt.Printf("Execution Time: %d ms\n\n", execResp.ExecutionTimeMs)
	}
}
