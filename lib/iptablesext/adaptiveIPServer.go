package iptablesext

import (
	"errors"
	"hcc/harp/lib/arping"
	"hcc/harp/lib/configadapriveipnetwork"
	"hcc/harp/lib/iplink"
	"hcc/harp/lib/logger"
	"os/exec"
)

// AdaptiveIPServerForwarding : Forwarding public IP address to private IP address
func AdaptiveIPServerForwarding(isAdd bool, preventInput bool, publicIP string, privateIP string) error {
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

	if preventInput {
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
	}

	cmd := exec.Command("iptables", "-t", "filter",
		"-C", HarpChainNamePrefix+"FORWARD",
		"-s", publicIP,
		"-j", "ACCEPT")
	err := cmd.Run()
	isExist := err == nil

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

// ControlNetDevAndIPTABLES : Add or delete iptables rules and virtual interface for AdaptiveIPServer
func ControlNetDevAndIPTABLES(isAdd bool, publicIP string, privateIP string) error {
	var err error

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()

	if isAdd {
		err = arping.CheckDuplicatedIPAddress(publicIP)
		if err != nil {
			return err
		}
	}

	if isAdd {
		err = iplink.AddOrDeleteIPToHarpExternalDevice(publicIP, adaptiveIP.Netmask, true)
		if err != nil {
			logger.Logger.Println("AddOrDeleteIPToHarpExternalDevice(): " + err.Error())
			return errors.New("failed to add AdaptiveIP IP address " + publicIP)
		}
	} else {
		err = iplink.AddOrDeleteIPToHarpExternalDevice(publicIP, adaptiveIP.Netmask, false)
		if err != nil {
			logger.Logger.Println("AddOrDeleteIPToHarpExternalDevice(): " + err.Error())
			return errors.New("failed to delete AdaptiveIP IP address " + publicIP)
		}
	}

	err = AdaptiveIPServerForwarding(isAdd, true, publicIP, privateIP)
	if err != nil {
		return err
	}

	err = ICMPForwarding(isAdd, publicIP, privateIP)
	if err != nil {
		return err
	}

	return nil
}
