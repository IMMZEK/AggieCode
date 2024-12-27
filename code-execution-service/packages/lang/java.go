package lang

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func ExecuteJavaCode(containerName, code string) (string, error) {
	// 1. Write the code to a temporary Main.java file.
	err := os.WriteFile("/tmp/Main.java", []byte(code), 0644)
	if err != nil {
		return "", err
	}

	// 2. Compile the code using javac inside the container.
	compileCmd := exec.Command(
		"docker", "exec", containerName,
		"javac", "/tmp/Main.java",
	)
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	if err := compileCmd.Run(); err != nil {
		return "", errors.New("Compilation error: " + compileErr.String())
	}

	// 3. Execute the compiled code using java inside the container.
	execCmd := exec.Command(
		"docker", "exec", containerName,
		"java", "-cp", "/tmp", "Main",
	)
	var out bytes.Buffer
	var execErr bytes.Buffer
	execCmd.Stdout = &out
	execCmd.Stderr = &execErr

	if err := execCmd.Run(); err != nil {
		return "", errors.New("Execution error: " + execErr.String())
	}

	return out.String(), nil
}
