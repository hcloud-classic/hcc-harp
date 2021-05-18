package iplink

import (
	"errors"
	"os/exec"
)

func runIP(args string) error {
	cmd := exec.Command("sh", "-c", "ip "+args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	return nil
}
