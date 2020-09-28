package syscheck

import (
	"bufio"
	"errors"
	"github.com/opencontainers/selinux/go-selinux"
	"os"
	"os/exec"
	"strings"
)

func checkPFEnabled() error {
	file, err := os.Open("/etc/rc.conf")
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	var pfEnabled = false

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ToLower(line)
		if strings.Contains(line, "pf_enable=\"yes\"") {
			pfEnabled = true
		}
	}

	if !pfEnabled {
		return errors.New("PF firewall is not enabled")
	}

	return nil
}

func checkIPTABLES() error {
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

// CheckFirewall : Check firewall for Linux and FreeBSD
func CheckFirewall() error {
	if OS == "freebsd" {
		return checkPFEnabled()
	} else if OS == "linux" {
		return checkIPTABLES()
	}

	return errors.New("failed to check firewall")
}
