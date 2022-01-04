package cmd

import (
	"errors"
	"os/exec"
)

func RunCMD(args string) error {
	cmd := exec.Command("sh", "-c", args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	return nil
}
