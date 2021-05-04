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
		addFlag, HarpAdaptiveIPInputDropChainName,
		"-d", publicIP,
		"-j", "DROP")
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to " + addErrMsg + " ADAPTIVE_IP_INPUT_DROP rule of " + publicIP)
	}

	cmd = exec.Command("iptables", "-t", "filter",
		addFlag, HarpChainNamePrefix+"FORWARD",
		"-s", publicIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to " + addErrMsg + " external FORWARD rule of " + publicIP)
	}

	cmd = exec.Command("iptables", "-t", "filter",
		addFlag, HarpChainNamePrefix+"FORWARD",
		"-d", privateIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to " + addErrMsg + " internal FORWARD rule of " + publicIP)
	}

	return nil
}

// CreateIPTABLESRulesAndExtIface : Check if public IP is duplicated then create
// ifconfig script file and virtual external interface, iptables rules.
func CreateIPTABLESRulesAndExtIface(publicIP string, privateIP string) error {
	adaptiveip := configext.GetAdaptiveIPNetwork()

	err := arping.CheckDuplicatedIPAddress(publicIP)
	if err != nil {
		return err
	}

	err = adaptiveIPServerForwarding(true, publicIP, privateIP)
	if err != nil {
		return err
	}

	err = ICMPForwarding(true, publicIP, privateIP)
	if err != nil {
		return err
	}

	err = ifconfig.IfconfigAddVirtualIface(config.AdaptiveIP.ExternalIfaceName, publicIP, adaptiveip.Netmask)
	if err != nil {
		return errors.New("failed to run ifconfig command of " + publicIP)
	}

	return nil
}

// DeleteIPTABLESRulesAndExtIface : Delete ifconfig script file and virtual interface, iptables rules
// match with public IP address.
func DeleteIPTABLESRulesAndExtIface(publicIP string, privateIP string) error {
	err := ICMPForwarding(false, publicIP, privateIP)
	if err != nil {
		return err
	}

	err = adaptiveIPServerForwarding(false, publicIP, privateIP)
	if err != nil {
		goto Error
	}

	err = ifconfig.IfconfigDeleteVirtualIface(config.AdaptiveIP.ExternalIfaceName, publicIP)
	if err != nil {
		goto Error
	}

	return nil
Error:
	return err
}
