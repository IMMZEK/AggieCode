/*
Golang Code Execution Service:
1.- Recieve a code snippet from the client
2.- Select the appropriate language runtime (Docker container)
3.- Run the docker container with the code snippet (sanitized and isolated)
4.- Return the output of the code snippet to the client
*/
//TODO: Implement the code execution logic
//TODO: Update readme with new tech stack and instrctions to run the service
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Placeholder types for request and response
type ExecuteCodeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type ExecuteCodeResponse struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

func executeCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ExecuteCodeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Placeholder response
	response := ExecuteCodeResponse{
		Output: fmt.Sprintf("Received code in %s: %s", req.Language, req.Code),
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/execute", executeCodeHandler).Methods("POST")

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
