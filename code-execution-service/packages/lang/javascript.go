package lang

import (
	"bytes"
	"errors"
	"os/exec"
)

func ExecuteJsCode(containerName, code string) (string, error) {
	// Execute the code using node inside the container.
	execCmd := exec.Command("docker", "exec", containerName, "node", "-e", code)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	execCmd.Stdout = &out
	execCmd.Stderr = &errBuf

	if err := execCmd.Run(); err != nil {
		return "", errors.New("Execution error: " + errBuf.String())
	}

	return out.String(), nil
}
