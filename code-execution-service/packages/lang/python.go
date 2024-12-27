package lang

import (
	"bytes"
	"errors"
	"os/exec"
)

func ExecutePythonCode(containerName, code string) (string, error) {
	// Execute the code using python3 inside the container.
	execCmd := exec.Command("docker", "exec", containerName, "python3", "-c", code)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	execCmd.Stdout = &out
	execCmd.Stderr = &errBuf

	if err := execCmd.Run(); err != nil {
		return "", errors.New("Execution error: " + errBuf.String())
	}

	return out.String(), nil
}
