package iptables

import (
	"hcc/harp/lib/logger"
	"os/exec"
)

// EnableIPForwardV4 : Enable IPv4 IP forward to use NAT with iptables.
func EnableIPForwardV4() error {
	logger.Logger.Println("Enabling ip_forward for IPv4...")

	cmd := exec.Command("echo", "1", ">", "/proc/sys/net/ipv4/ip_forward")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
