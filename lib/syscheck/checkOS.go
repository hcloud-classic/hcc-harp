package syscheck

import (
	"errors"
	"fmt"
	"runtime"
)

// OS : Contain currently running OS
var OS string

// CheckOS : Check OS then return error if not Linux or FreeBSD
func CheckOS() error {
	OS = runtime.GOOS

	switch OS {
	case "linux":
		fmt.Println("Running harp module on Linux machine")
	case "freebsd":
		fmt.Println("Running harp module on FreeBSD machine")
	default:
		return errors.New("this machine is not compatible with harp module")
	}

	return nil
}
