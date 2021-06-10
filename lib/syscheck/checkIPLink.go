package syscheck

import (
	"errors"
	"os/exec"
)

// CheckIPLink : Check if we are ready to use ip command
func CheckIPLink() error {
	cmd := exec.Command("ip", "link", "show")
	err := cmd.Run()
	if err != nil {
		return errors.New("ip command not found")
	}

	return nil
}
