package iptablesext

import (
	"errors"
	"hcc/harp/lib/arping"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/ifconfig"
	"hcc/harp/lib/logger"
	"os/exec"
)

func adaptiveIPServerForwarding(isAdd bool, publicIP string, privateIP string) error {
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

	logger.Logger.Println(addMsg + " AdaptiveIP Server forwarding iptables rules for " + publicIP + " (privateIP: " + privateIP + ")")

	cmd := exec.Command("iptables", "-t", "filter",
		"-C", HarpAdaptiveIPInputDropChainName,
		"-d", publicIP,
		"-j", "DROP")
	err := cmd.Run()
	isExist := err == nil

	if (isAdd && !isExist) || (!isAdd && isExist) {
		cmd = exec.Command("iptables", "-t", "filter",
			addFlag, HarpAdaptiveIPInputDropChainName,
			"-d", publicIP,
			"-j", "DROP")
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to " + addErrMsg + " ADAPTIVE_IP_INPUT_DROP rule of " + publicIP)
		}
	}

	cmd = exec.Command("iptables", "-t", "filter",
		"-C", HarpChainNamePrefix+"FORWARD",
		"-s", publicIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	isExist = err == nil

	if (isAdd && !isExist) || (!isAdd && isExist) {
		cmd = exec.Command("iptables", "-t", "filter",
			addFlag, HarpChainNamePrefix+"FORWARD",
			"-s", publicIP,
			"-j", "ACCEPT")
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to " + addErrMsg + " external FORWARD rule of " + publicIP)
		}
	}

	cmd = exec.Command("iptables", "-t", "filter",
		"-C", HarpChainNamePrefix+"FORWARD",
		"-d", privateIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	isExist = err == nil

	if (isAdd && !isExist) || (!isAdd && isExist) {
		cmd = exec.Command("iptables", "-t", "filter",
			addFlag, HarpChainNamePrefix+"FORWARD",
			"-d", privateIP,
			"-j", "ACCEPT")
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to " + addErrMsg + " internal FORWARD rule of " + publicIP)
		}
	}

	return nil
}

// ControlIfconfigAndIPTABLES : Add or delete iptables rules and virtual interface for AdaptiveIPServer
func ControlIfconfigAndIPTABLES(isAdd bool, publicIP string, privateIP string) error {
	var err error

	adaptiveIP := configext.GetAdaptiveIPNetwork()

	if isAdd {
		err = arping.CheckDuplicatedIPAddress(publicIP)
		if err != nil {
			return err
		}
	}

	err = adaptiveIPServerForwarding(isAdd, publicIP, privateIP)
	if err != nil {
		return err
	}

	err = ICMPForwarding(isAdd, publicIP, privateIP)
	if err != nil {
		return err
	}

	if isAdd {
		err = ifconfig.AddVirtualIface(config.AdaptiveIP.ExternalIfaceName, publicIP, adaptiveIP.Netmask)
		if err != nil {
			return errors.New("failed to add virtual interface for " + publicIP)
		}
	} else {
		err = ifconfig.DeleteVirtualIface(config.AdaptiveIP.ExternalIfaceName, publicIP)
		if err != nil {
			return errors.New("failed to delete virtual interface for " + publicIP)
		}
	}

	return nil
}
