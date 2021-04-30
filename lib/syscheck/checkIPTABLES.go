package syscheck

import (
	"errors"
	"github.com/opencontainers/selinux/go-selinux"
	"os/exec"
)

// CheckIPTABLES : Check if we are ready to use iptables
func CheckIPTABLES() error {
	if selinux.GetEnabled() {
		return errors.New("SELinux is Enabled. Disable SELinux and try again")
	}

	cmd := exec.Command("iptables", "--help")
	err := cmd.Run()
	if err != nil {
		return errors.New("iptables command not found")
	}

	return nil

}
