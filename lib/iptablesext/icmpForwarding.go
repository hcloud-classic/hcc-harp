package iptablesext

import (
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"os/exec"
)

// AddICMPForwarding : Add iptables rules for ICMP forwarding
func AddICMPForwarding(publicIP string, privateIP string) error {
	logger.Logger.Println("Adding ICMP forwarding iptables rules for " + publicIP + " (privateIP: " + privateIP + ")")

	cmd := exec.Command("iptables", "-t", "nat",
		"-A", HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.InternalIfaceName,
		"-p", "icmp", "--icmp-type", "echo-request",
		"-j", "SNAT",
		"--to-source", publicIP)
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to add ICMP POSTROUTING rule of " + publicIP)
	}

	cmd = exec.Command("iptables", "-t", "nat",
		"-A", HarpChainNamePrefix+"PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
		"-p", "icmp", "--icmp-type", "echo-request",
		"-d", publicIP,
		"-j", "DNAT",
		"--to-destination", privateIP)
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to add ICMP PREROUTING rule of " + publicIP)
	}

	cmd = exec.Command("iptables", "-t", "filter",
		"-A", HarpChainNamePrefix+"INPUT",
		"-p", "icmp", "--icmp-type", "echo-request",
		"-d", publicIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to add ICMP INPUT rule of " + publicIP)
	}

	return nil
}

// DeleteICMPForwarding : Delete iptables rules for ICMP forwarding
func DeleteICMPForwarding(publicIP string, privateIP string) error {
	logger.Logger.Println("Deleting ICMP forwarding iptables rules for " + publicIP + " (privateIP: " + privateIP + ")")

	cmd := exec.Command("iptables", "-t", "nat",
		"-D", HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.InternalIfaceName,
		"-p", "icmp", "--icmp-type", "echo-request",
		"-j", "SNAT",
		"--to-source", publicIP)
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to delete ICMP POSTROUTING rule of " + publicIP)
	}

	cmd = exec.Command("iptables", "-t", "nat",
		"-D", HarpChainNamePrefix+"PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
		"-p", "icmp", "--icmp-type", "echo-request",
		"-d", publicIP,
		"-j", "DNAT",
		"--to-destination", privateIP)
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to delete ICMP PREROUTING rule of " + publicIP)
	}

	cmd = exec.Command("iptables", "-t", "filter",
		"-D", HarpChainNamePrefix+"INPUT",
		"-p", "icmp", "--icmp-type", "echo-request",
		"-d", publicIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to delete ICMP INPUT rule of " + publicIP)
	}

	return nil
}
