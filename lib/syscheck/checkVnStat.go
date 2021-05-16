package syscheck

import (
	"errors"
	"os/exec"
)

// CheckVnStat : Check if we are ready to use vnStat
func CheckVnStat() error {
	cmd := exec.Command("vnstat", "--help")
	err := cmd.Run()
	if err != nil {
		return errors.New("please install vnstat")
	}

	return nil

}
