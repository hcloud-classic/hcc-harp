package iptables

import (
	"hcc/harp/lib/logger"
	"os/exec"
)

// FlushIPTABLESRules : Remove all of configured iptables rules
func FlushIPTABLESRules() error {
	logger.Logger.Println("Flushing iptables rules...")

	cmd := exec.Command("iptables", "-F")
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("iptables", "-X")
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("iptables", "-t", "nat", "-F")
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("iptables", "-t", "nat", "-X")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
