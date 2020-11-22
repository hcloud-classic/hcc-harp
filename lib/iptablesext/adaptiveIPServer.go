package iptablesext

import (
	"hcc/harp/lib/arping"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/ifconfig"
	"hcc/harp/lib/logger"
	"os/exec"
)

func addAdaptiveIPServerIPTABLESRules(publicIP string, privateIP string) error {
	logger.Logger.Println("Adding AdaptiveIP Server iptables rules for " + publicIP + " (privateIP: " + privateIP + ")")

	cmd := exec.Command("iptables", "-t", "nat",
		"-A", "POSTROUTING", "-o", config.AdaptiveIP.ExternalIfaceName,
		"-s", privateIP,
		"-j", "SNAT",
		"--to-source", publicIP)
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("iptables", "-t", "nat",
		"-A", "PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
		"-d", publicIP,
		"-j", "DNAT",
		"--to-destination", privateIP)
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("iptables",
		"-A", "FORWARD",
		"-s", publicIP,
		"-j", "ACCEPT")
	err = cmd.Run()

	cmd = exec.Command("iptables",
		"-A", "FORWARD",
		"-d", privateIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func deleteAdaptiveIPServerIPTABLESRules(publicIP string, privateIP string) error {
	logger.Logger.Println("Deleting AdaptiveIP Server iptables rules for " + publicIP + " (privateIP: " + privateIP + ")")

	cmd := exec.Command("iptables", "-t", "nat",
		"-D", "POSTROUTING", "-o", config.AdaptiveIP.ExternalIfaceName,
		"-s", privateIP,
		"-j", "SNAT",
		"--to-source", publicIP)
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("iptables", "-t", "nat",
		"-D", "PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
		"-d", publicIP,
		"-j", "DNAT",
		"--to-destination", privateIP)
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("iptables",
		"-D", "FORWARD",
		"-s", publicIP,
		"-j", "ACCEPT")
	err = cmd.Run()

	cmd = exec.Command("iptables",
		"-D", "FORWARD",
		"-d", privateIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// CreateIPTABLESRulesAndExtIface : Check if public IP is duplicated then create
// ifconfig script file and virtual external interface, iptables rules.
func CreateIPTABLESRulesAndExtIface(publicIP string, privateIP string) error {
	adaptiveip := configext.GetAdaptiveIPNetwork()

	err := arping.CheckDuplicatedIPAddress(publicIP)
	if err != nil {
		goto Error
	}

	err = addAdaptiveIPServerIPTABLESRules(publicIP, privateIP)
	if err != nil {
		goto Error
	}

	err = ifconfig.CreateAndLoadIfconfigScriptExternal(config.AdaptiveIP.ExternalIfaceName, publicIP,
		adaptiveip.Netmask)
	if err != nil {
		goto Error
	}

	return nil
Error:
	return err
}

// DeleteIPTABLESRulesAndExtIface : Delete ifconfig script file and virtual interface, iptables rules
// match with public IP address.
func DeleteIPTABLESRulesAndExtIface(publicIP string, privateIP string) error {
	adaptiveip := configext.GetAdaptiveIPNetwork()

	err := deleteAdaptiveIPServerIPTABLESRules(publicIP, privateIP)
	if err != nil {
		goto Error
	}

	err = ifconfig.DeleteAndUnloadIfconfigScriptExternal(config.AdaptiveIP.ExternalIfaceName, publicIP,
		adaptiveip.Netmask)
	if err != nil {
		goto Error
	}

	return nil
Error:
	return err
}
