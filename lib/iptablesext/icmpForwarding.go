package iptablesext

import (
	"errors"
	"hcc/harp/lib/adaptiveipext"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configadapriveipnetwork"
	"hcc/harp/lib/logger"
	"os/exec"
)

// ICMPForwarding : Add or delete iptables rules for ICMP forwarding
func ICMPForwarding(isAdd bool, internalIP string, privateIP string) error {
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

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()
	externalIP, _ := adaptiveipext.InternalIPtoExternalIP(internalIP)
	if len(externalIP) != 0 {
		externalIP = " (" + externalIP + ")"
	}

	logger.Logger.Println(addMsg + " ICMP forwarding iptables rules for " +
		internalIP + externalIP + " (privateIP: " + privateIP + ")")

	cmd := exec.Command("iptables", "-t", "nat",
		"-C", HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.InternalIfaceName,
		"-p", "icmp", "--icmp-type", "echo-request",
		"-d", internalIP,
		"-j", "SNAT",
		"--to-source", adaptiveIP.ExtIfaceIPAddress)
	err := cmd.Run()
	isExist := err == nil

	if (isAdd && !isExist) || (!isAdd && isExist) {
		cmd = exec.Command("iptables", "-t", "nat",
			addFlag, HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.InternalIfaceName,
			"-p", "icmp", "--icmp-type", "echo-request",
			"-d", internalIP,
			"-j", "SNAT",
			"--to-source", adaptiveIP.ExtIfaceIPAddress)
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to " + addErrMsg + " ICMP POSTROUTING rule of " +
				internalIP + externalIP)
		}
	}

	cmd = exec.Command("iptables", "-t", "nat",
		"-C", HarpChainNamePrefix+"PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
		"-p", "icmp", "--icmp-type", "echo-request",
		"-d", internalIP,
		"-j", "DNAT",
		"--to-destination", privateIP)
	err = cmd.Run()
	isExist = err == nil

	if (isAdd && !isExist) || (!isAdd && isExist) {
		cmd = exec.Command("iptables", "-t", "nat",
			addFlag, HarpChainNamePrefix+"PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
			"-p", "icmp", "--icmp-type", "echo-request",
			"-d", internalIP,
			"-j", "DNAT",
			"--to-destination", privateIP)
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to " + addErrMsg + " ICMP PREROUTING rule of " +
				internalIP + externalIP)
		}
	}

	cmd = exec.Command("iptables", "-t", "filter",
		"-C", HarpChainNamePrefix+"INPUT",
		"-p", "icmp", "--icmp-type", "echo-request",
		"-d", internalIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	isExist = err == nil

	if (isAdd && !isExist) || (!isAdd && isExist) {
		if isAdd {
			cmd = exec.Command("iptables", "-t", "filter",
				"-I", HarpChainNamePrefix+"INPUT", "1",
				"-p", "icmp", "--icmp-type", "echo-request",
				"-d", internalIP,
				"-j", "ACCEPT")
			err = cmd.Run()
			if err != nil {
				return errors.New("failed to " + addErrMsg + " ICMP INPUT rule of " +
					internalIP + externalIP)
			}
		} else {
			cmd = exec.Command("iptables", "-t", "filter",
				addFlag, HarpChainNamePrefix+"INPUT",
				"-p", "icmp", "--icmp-type", "echo-request",
				"-d", internalIP,
				"-j", "ACCEPT")
			err = cmd.Run()
			if err != nil {
				return errors.New("failed to " + addErrMsg + " ICMP INPUT rule of " +
					internalIP + externalIP)
			}
		}
	}

	return nil
}
