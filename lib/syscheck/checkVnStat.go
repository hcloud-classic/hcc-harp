package syscheck

import (
	"errors"
	"os/exec"
	"strings"
)

// CheckVnStat : Check if we are ready to use vnStat
func CheckVnStat() error {
	cmd := exec.Command("vnstat", "--help")
	err := cmd.Run()
	if err != nil {
		return errors.New("please install VnStat v1.18")
	}

	cmd = exec.Command("sh", "-c", "vnstat --version | awk '{print $2}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(err.Error())
	}

	version := strings.TrimSpace(string(output))
	if version != "1.18" {
		return errors.New("unsupported VnStat version is installed. Please install VnStat v1.18")
	}

	return nil
}
