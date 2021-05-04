package iptablesext

import (
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/logger"
	"os/exec"
)

// ICMPForwarding : Add or delete iptables rules for ICMP forwarding
func ICMPForwarding(isAdd bool, publicIP string, privateIP string) error {
	var addMsg string
	var addErrMsg string
	var addFlag string

	if isAdd {
		addMsg = "Adding"
		addErrMsg = "add"
		addFlag = "-A"
	} else {
		addMsg = "Deleting"
		addErrMsg = "delete"
		addFlag = "-D"
	}

	adaptiveIP := configext.GetAdaptiveIPNetwork()

	logger.Logger.Println(addMsg + " ICMP forwarding iptables rules for " + publicIP + " (privateIP: " + privateIP + ")")

	cmd := exec.Command("iptables", "-t", "nat",
		"-C", HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.InternalIfaceName,
		"-p", "icmp", "--icmp-type", "echo-request",
		"-d", publicIP,
		"-j", "SNAT",
		"--to-source", adaptiveIP.ExtIfaceIPAddress)
	err := cmd.Run()
	isExist := err == nil

	if (isAdd && !isExist) || (!isAdd && isExist) {
		cmd = exec.Command("iptables", "-t", "nat",
			addFlag, HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.InternalIfaceName,
			"-p", "icmp", "--icmp-type", "echo-request",
			"-d", publicIP,
			"-j", "SNAT",
			"--to-source", adaptiveIP.ExtIfaceIPAddress)
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to " + addErrMsg + " ICMP POSTROUTING rule of " + publicIP)
		}
	}

	cmd = exec.Command("iptables", "-t", "nat",
		"-C", HarpChainNamePrefix+"PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
		"-p", "icmp", "--icmp-type", "echo-request",
		"-d", publicIP,
		"-j", "DNAT",
		"--to-destination", privateIP)
	err = cmd.Run()
	isExist = err == nil

	if (isAdd && !isExist) || (!isAdd && isExist) {
		cmd = exec.Command("iptables", "-t", "nat",
			addFlag, HarpChainNamePrefix+"PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
			"-p", "icmp", "--icmp-type", "echo-request",
			"-d", publicIP,
			"-j", "DNAT",
			"--to-destination", privateIP)
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to " + addErrMsg + " ICMP PREROUTING rule of " + publicIP)
		}
	}

	cmd = exec.Command("iptables", "-t", "filter",
		"-C", HarpChainNamePrefix+"INPUT",
		"-p", "icmp", "--icmp-type", "echo-request",
		"-d", publicIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	isExist = err == nil

	if (isAdd && !isExist) || (!isAdd && isExist) {
		cmd = exec.Command("iptables", "-t", "filter",
			addFlag, HarpChainNamePrefix+"INPUT",
			"-p", "icmp", "--icmp-type", "echo-request",
			"-d", publicIP,
			"-j", "ACCEPT")
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to " + addErrMsg + " ICMP INPUT rule of " + publicIP)
		}
	}

	return nil
}
