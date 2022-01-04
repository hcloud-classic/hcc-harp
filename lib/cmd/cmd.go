package cmd

import (
	"errors"
	"os/exec"
)

func RunCMD(cmdWithArgs string) error {
	cmd := exec.Command("sh", "-c", cmdWithArgs)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	return nil
}

func RunScript(filename string) error {
	cmd := exec.Command("sh", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	return nil
}
