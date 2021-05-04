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

func addAdaptiveIPServerIPTABLESRules(publicIP string, privateIP string) error {
	logger.Logger.Println("Adding AdaptiveIP Server iptables rules for " + publicIP + " (privateIP: " + privateIP + ")")
	cmd := exec.Command("iptables", "-t", "filter",
		"-A", HarpAdaptiveIPInputDropChainName,
		"-d", publicIP,
		"-j", "DROP")
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to add ADAPTIVE_IP_INPUT_DROP rule of " + publicIP)
	}

	cmd = exec.Command("iptables", "-t", "filter",
		"-A", HarpChainNamePrefix+"FORWARD",
		"-s", publicIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to add external FORWARD rule of " + publicIP)
	}

	cmd = exec.Command("iptables", "-t", "filter",
		"-A", HarpChainNamePrefix+"FORWARD",
		"-d", privateIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to add internal FORWARD rule of " + publicIP)
	}

	return nil
}

func deleteAdaptiveIPServerIPTABLESRules(publicIP string, privateIP string) error {
	logger.Logger.Println("Deleting AdaptiveIP Server iptables rules for " + publicIP + " (privateIP: " + privateIP + ")")
	cmd := exec.Command("iptables", "-t", "filter",
		"-D", HarpAdaptiveIPInputDropChainName,
		"-d", publicIP,
		"-j", "DROP")
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to delete ADAPTIVE_IP_INPUT_DROP rule of " + publicIP)
	}

	cmd = exec.Command("iptables", "-t", "filter",
		"-D", HarpChainNamePrefix+"FORWARD",
		"-s", publicIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to delete external FORWARD rule of " + publicIP)
	}

	cmd = exec.Command("iptables", "-t", "filter",
		"-D", HarpChainNamePrefix+"FORWARD",
		"-d", privateIP,
		"-j", "ACCEPT")
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to delete internal FORWARD rule of " + publicIP)
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

	err = addAdaptiveIPServerIPTABLESRules(publicIP, privateIP)
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

	err = deleteAdaptiveIPServerIPTABLESRules(publicIP, privateIP)
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
