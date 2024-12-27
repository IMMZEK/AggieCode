package lang

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func ExecuteCppCode(containerName, code string) (string, error) {
	// 1. Write the code to a temporary .cpp file.
	err := os.WriteFile("/tmp/main.cpp", []byte(code), 0644)
	if err != nil {
		return "", err
	}

	// 2. Compile the code using g++ inside the container.
	//    We mount the temporary directory to share files with the container.
	compileCmd := exec.Command(
		"docker", "exec", containerName,
		"g++", "-o", "/tmp/main", "/tmp/main.cpp",
	)
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	if err := compileCmd.Run(); err != nil {
		return "", errors.New("Compilation error: " + compileErr.String())
	}

	// 3. Execute the compiled binary inside the container.
	execCmd := exec.Command("docker", "exec", containerName, "/tmp/main")
	var out bytes.Buffer
	var execErr bytes.Buffer
	execCmd.Stdout = &out
	execCmd.Stderr = &execErr

	if err := execCmd.Run(); err != nil {
		return "", errors.New("Execution error: " + execErr.String())
	}

	return out.String(), nil
}
