package lang

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExecuteGoCode(containerName, code string) (string, error) {
	// Create a temporary directory for the Go project
	tmpDir, err := os.MkdirTemp("", "go-project-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up the temporary directory when done

	// Create a temporary file for the Go code inside the temporary directory
	tmpFile := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(tmpFile, []byte(code), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write Go code to temporary file: %v", err)
	}

	// Initialize a Go module inside the temporary directory
	moduleName := "temp-go-module" // You can generate this dynamically if needed
	initCmd := exec.Command("docker", "exec", containerName, "go", "mod", "init", moduleName)
	initCmd.Dir = "/tmp/" + filepath.Base(tmpDir)
	if out, err := initCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to initialize Go module: %v, output: %s", err, out)
	}

	// Run the Go code inside the container using `go run`
	runCmd := exec.Command("docker", "exec", "-w", "/tmp/"+filepath.Base(tmpDir), containerName, "go", "run", "main.go")
	var out, errBuf bytes.Buffer
	runCmd.Stdout = &out
	runCmd.Stderr = &errBuf

	if err := runCmd.Run(); err != nil {
		// Check if the error is due to compilation or runtime issues
		if strings.Contains(errBuf.String(), "go: go.mod file not found") {
			return "", fmt.Errorf("execution error: go.mod file not found. Ensure your code includes a valid 'module' directive at the beginning. Output: %s", errBuf.String())
		}
		return "", fmt.Errorf("execution error: %v, output: %s", err, errBuf.String())
	}

	return out.String(), nil
}
