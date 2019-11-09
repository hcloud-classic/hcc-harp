package syscheck

import (
	"errors"
	"fmt"
	"os/exec"
)

// CheckArpingCommand : Check if 'arping' command is available from local system.
func CheckArpingCommand() error {
	cmd := exec.Command("arping", "--help")
	err := cmd.Run()
	if err != nil {
		fmt.Println("arping command not found! Please install arping first.")
		return errors.New("arping command not found")
	}

	return nil
}
