package syscheck

import (
	"errors"
	"runtime"
)

// CheckOS : Check OS then return error if not Linux
func CheckOS() error {
	if runtime.GOOS != "linux" {
		return errors.New("this machine is not compatible with harp module")
	}

	return nil
}
