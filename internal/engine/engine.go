package engine

import (
	"fmt"
	"os"
	"os/exec"
)

func ExecuteScript(data []byte) error {
	tmp, err := os.CreateTemp("", "task-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		return fmt.Errorf("fauked ti write script: %w", err)
	}
	tmp.Close()

	if err := os.Chmod(tmp.Name(), 0700); err != nil {
		return fmt.Errorf("failed to chmod script: %w", err)
	}
	cmd := exec.Command(tmp.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("script execution failed: %w", err)
	}
	return nil
}
