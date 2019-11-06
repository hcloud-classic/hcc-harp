package syscheck

import (
	"fmt"
	"os/exec"
)

// CheckArpingCommand : Check if 'arping' command is available from local system.
func CheckArpingCommand() bool {
	cmd := exec.Command("arping", "--help")
	err := cmd.Run()
	if err != nil {
		fmt.Println("arping command not found! Please install arping first.")
		return false
	}

	return true
}
