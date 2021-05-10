package iptablesext

import (
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/logger"
	"os/exec"
	"strconv"
	"strings"
)

// PortForwarding : Add or delete iptables rules for port forwarding
func PortForwarding(isAdd bool, forwardTCP bool, forwardUDP bool, publicIP string, privateIP string,
	externalPort int, internalPort int) error {
	var addMsg string
	var addErrMsg string
	var addFlag string
	var protocol []string

	if isAdd {
		addMsg = "Adding"
		addErrMsg = "add"
		addFlag = "-A"
	} else {
		addMsg = "Deleting"
		addErrMsg = "delete"
		addFlag = "-D"
	}

	if forwardTCP && forwardUDP {
		protocol = []string{"tcp", "udp"}
	} else if forwardTCP {
		protocol = []string{"tcp"}
	} else if forwardUDP {
		protocol = []string{"udp"}
	} else {
		return errors.New("protocol is not selected")
	}

	adaptiveIP := configext.GetAdaptiveIPNetwork()

	for i := range protocol {
		logger.Logger.Println(addMsg + " " + strings.ToUpper(protocol[i]) + " forwarding iptables rules for " +
			publicIP + ":" + strconv.Itoa(externalPort) + " (private: " + privateIP + ":" + strconv.Itoa(internalPort) + ")")

		cmd := exec.Command("iptables", "-t", "nat",
			"-C", HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.InternalIfaceName,
			"-p", protocol[i], "--dport", strconv.Itoa(externalPort),
			"-d", publicIP,
			"-j", "SNAT",
			"--to-source", adaptiveIP.ExtIfaceIPAddress)
		err := cmd.Run()
		isExist := err == nil

		if (isAdd && !isExist) || (!isAdd && isExist) {
			cmd = exec.Command("iptables", "-t", "nat",
				addFlag, HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.InternalIfaceName,
				"-p", protocol[i], "--dport", strconv.Itoa(externalPort),
				"-d", publicIP,
				"-j", "SNAT",
				"--to-source", adaptiveIP.ExtIfaceIPAddress)
			err = cmd.Run()
			if err != nil {
				return errors.New("failed to " + addErrMsg + " " + strings.ToUpper(protocol[i]) +
					" POSTROUTING rule of " + publicIP + ":" + strconv.Itoa(externalPort))
			}
		}

		cmd = exec.Command("iptables", "-t", "nat",
			"-C", HarpChainNamePrefix+"PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
			"-p", protocol[i], "--dport", strconv.Itoa(externalPort),
			"-d", publicIP,
			"-j", "DNAT",
			"--to-destination", privateIP+":"+strconv.Itoa(internalPort))
		err = cmd.Run()
		isExist = err == nil

		if (isAdd && !isExist) || (!isAdd && isExist) {
			cmd = exec.Command("iptables", "-t", "nat",
				addFlag, HarpChainNamePrefix+"PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
				"-p", protocol[i], "--dport", strconv.Itoa(externalPort),
				"-d", publicIP,
				"-j", "DNAT",
				"--to-destination", privateIP+":"+strconv.Itoa(internalPort))
			err = cmd.Run()
			if err != nil {
				return errors.New("failed to " + addErrMsg + " " + strings.ToUpper(protocol[i]) +
					" PREROUTING rule of " + publicIP + ":" + strconv.Itoa(externalPort))
			}
		}

		cmd = exec.Command("iptables", "-t", "filter",
			"-C", HarpChainNamePrefix+"INPUT",
			"-p", protocol[i], "--dport", strconv.Itoa(externalPort),
			"-d", publicIP,
			"-j", "ACCEPT")
		err = cmd.Run()
		isExist = err == nil

		if (isAdd && !isExist) || (!isAdd && isExist) {
			if isAdd {
				cmd = exec.Command("iptables", "-t", "filter",
					"-I", HarpChainNamePrefix+"INPUT", "1",
					"-p", protocol[i], "--dport", strconv.Itoa(externalPort),
					"-d", publicIP,
					"-j", "ACCEPT")
				err = cmd.Run()
				if err != nil {
					return errors.New("failed to " + addErrMsg + " " + strings.ToUpper(protocol[i]) +
						" INPUT rule of " + publicIP + ":" + strconv.Itoa(externalPort))
				}
			} else {
				cmd = exec.Command("iptables", "-t", "filter",
					addFlag, HarpChainNamePrefix+"INPUT",
					"-p", protocol[i], "--dport", strconv.Itoa(externalPort),
					"-d", publicIP,
					"-j", "ACCEPT")
				err = cmd.Run()
				if err != nil {
					return errors.New("failed to " + addErrMsg + " " + strings.ToUpper(protocol[i]) +
						" INPUT rule of " + publicIP + ":" + strconv.Itoa(externalPort))
				}
			}
		}
	}

	return nil
}
