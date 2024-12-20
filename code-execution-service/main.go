package main

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gorilla/mux"
)

// Define a struct to represent the request body
type ExecuteCodeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Input    string `json:"input"`
}

// Define a struct to represent the response body
type ExecuteCodeResponse struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

func executeCodeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Decode the request body
	var req ExecuteCodeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Create a Docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		http.Error(w, "Failed to create Docker client", http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	// 3. Pull the Docker image for the requested language (if not already present)
	imageName := fmt.Sprintf("your-docker-hub-username/%s-executor:latest", req.Language) // Example: "your-docker-hub-username/python-executor:latest"
	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to pull Docker image: %s", imageName), http.StatusInternalServerError)
		return
	}

	// 4. Create a Docker container
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageName,
			Cmd:          []string{"/bin/sh", "-c", getExecutionCommand(req.Language, "code")}, // Use a helper function to get the command
			Tty:          false,
			AttachStdout: true,
			AttachStderr: true,
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory: 128 * 1024 * 1024, // Example: Limit to 128MB of memory
				// Add other resource limits as needed (CPU, etc.)
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		http.Error(w, "Failed to create Docker container", http.StatusInternalServerError)
		return
	}

	// 5. Copy the code into the container
	codeFileName := getCodeFileName(req.Language) // Use a helper function to get the file name based on language
	err = copyToContainer(ctx, cli, resp.ID, "/app/"+codeFileName, req.Code)
	if err != nil {
		http.Error(w, "Failed to copy code to container", http.StatusInternalServerError)
		return
	}

	// 6. Start the container
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		http.Error(w, "Failed to start Docker container", http.StatusInternalServerError)
		return
	}

	// 7. Set a timeout for execution (e.g., 5 seconds)
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			http.Error(w, "Error waiting for container execution", http.StatusInternalServerError)
			return
		}
	case <-statusCh:
		// Container execution finished
	case <-ctx.Done():
		// Timeout occurred, stop the container
		log.Println("Timeout occurred, stopping container...")
		timeoutErr := cli.ContainerStop(context.Background(), resp.ID, container.StopOptions{})
		if timeoutErr != nil {
			log.Printf("Error stopping container: %v\n", timeoutErr)
		}
		http.Error(w, "Code execution timed out", http.StatusRequestTimeout)
		return
	}

	// 8. Get the container logs (stdout and stderr)
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		http.Error(w, "Failed to get container logs", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 9. Read the logs into a buffer
	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, out)
	if err != nil {
		http.Error(w, "Failed to read container logs", http.StatusInternalServerError)
		return
	}

	// 10. Remove the container (optional, you might want to keep it for debugging)
	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		log.Printf("Failed to remove container: %v\n", err)
	}

	// 11. Prepare the response
	response := ExecuteCodeResponse{
		Output: stdoutBuf.String(),
		Error:  stderrBuf.String(),
	}

	// 12. Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to get the execution command based on the language
func getExecutionCommand(language, codeFileName string) string {
	switch language {
	case "python":
		return fmt.Sprintf("python3 /app/%s", codeFileName)
	case "cpp":
		return fmt.Sprintf("g++ /app/%s -o /app/a.out && /app/a.out", codeFileName)
	case "java":
		return fmt.Sprintf("javac /app/%s && java /app/%s", codeFileName, codeFileName)
	case "javascript":
		return fmt.Sprintf("node /app/%s", codeFileName)
	default:
		return ""
	}
}

// Helper function to get the code file name based on the language
func getCodeFileName(language string) string {
	switch language {
	case "python":
		return "code.py"
	case "cpp":
		return "code.cpp"
	case "java":
		return "Main.java"
	case "javascript":
		return "code.js"
	default:
		return "code.txt"
	}
}

// Helper function to copy a file to a Docker container
func copyToContainer(ctx context.Context, cli *client.Client, containerID, dstPath, content string) error {
	// Create a reader for the content
	reader := bytes.NewReader([]byte(content))

	// Create a tar archive
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	header := &tar.Header{
		Name: filepath.Base(dstPath),
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := io.Copy(tw, reader); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}

	// Copy the tar archive to the container
	return cli.CopyToContainer(ctx, containerID, filepath.Dir(dstPath), &buf, types.CopyToContainerOptions{})
}

func main() {
	// Use gorilla/mux for routing (optional, you can use the standard net/http as well)
	r := mux.NewRouter()
	r.HandleFunc("/api/execute", executeCodeHandler).Methods("POST")

	// Start the HTTP server
	log.Println("Starting code execution service on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
