package iptablesext

import (
	"errors"
	"hcc/harp/lib/adaptiveipext"
	"hcc/harp/lib/arping"
	"hcc/harp/lib/configadapriveipnetwork"
	"hcc/harp/lib/iplink"
	"hcc/harp/lib/logger"
	"os/exec"
)

// AdaptiveIPServerForwarding : Forwarding internal IP address to private IP address
func AdaptiveIPServerForwarding(isAdd bool, preventInput bool, internalIP string, privateIP string) error {
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

	externalIP, _ := adaptiveipext.InternalIPtoExternalIP(internalIP)
	if len(externalIP) != 0 {
		externalIP = " (" + externalIP + ")"
	}

	logger.Logger.Println(addMsg + " AdaptiveIP Server forwarding iptables rules for " +
		internalIP + externalIP + " (privateIP: " + privateIP + ")")

	if preventInput {
		cmd := exec.Command("iptables", "-t", "filter",
			"-C", HarpAdaptiveIPInputDropChainName,
			"-d", internalIP,
			"-j", "DROP")
		err := cmd.Run()
		isExist := err == nil

		if (isAdd && !isExist) || (!isAdd && isExist) {
			cmd = exec.Command("iptables", "-t", "filter",
				addFlag, HarpAdaptiveIPInputDropChainName,
				"-d", internalIP,
				"-j", "DROP")
			err = cmd.Run()
			if err != nil {
				return errors.New("failed to " + addErrMsg + " ADAPTIVE_IP_INPUT_DROP rule of " +
					internalIP + externalIP)
			}
		}
	}

	cmd := exec.Command("iptables", "-t", "filter",
		"-C", HarpChainNamePrefix+"FORWARD",
		"-s", internalIP,
		"-j", "ACCEPT")
	err := cmd.Run()
	isExist := err == nil

	if (isAdd && !isExist) || (!isAdd && isExist) {
		cmd = exec.Command("iptables", "-t", "filter",
			addFlag, HarpChainNamePrefix+"FORWARD",
			"-s", internalIP,
			"-j", "ACCEPT")
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to " + addErrMsg + " external FORWARD rule of " +
				internalIP + externalIP)
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
			return errors.New("failed to " + addErrMsg + " internal FORWARD rule of " +
				internalIP + externalIP)
		}
	}

	return nil
}

// ControlNetDevAndIPTABLES : Add or delete iptables rules and virtual interface for AdaptiveIPServer
func ControlNetDevAndIPTABLES(isAdd bool, publicIP string, privateIP string) error {
	var err error

	internalIP, err := adaptiveipext.ExternalIPtoInternalIP(publicIP)
	if err != nil {
		return err
	}

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()

	if isAdd {
		err = arping.CheckDuplicatedIPAddress(internalIP)
		if err != nil {
			return err
		}
	}

	if isAdd {
		err = iplink.AddOrDeleteIPToHarpExternalDevice(internalIP, adaptiveIP.Netmask, true)
		if err != nil {
			logger.Logger.Println("AddOrDeleteIPToHarpExternalDevice(): " + err.Error())
			return errors.New("failed to add AdaptiveIP IP address " + publicIP)
		}
	} else {
		err = iplink.AddOrDeleteIPToHarpExternalDevice(internalIP, adaptiveIP.Netmask, false)
		if err != nil {
			logger.Logger.Println("AddOrDeleteIPToHarpExternalDevice(): " + err.Error())
			return errors.New("failed to delete AdaptiveIP IP address " + publicIP)
		}
	}

	err = AdaptiveIPServerForwarding(isAdd, true, internalIP, privateIP)
	if err != nil {
		return err
	}

	err = ICMPForwarding(isAdd, internalIP, privateIP)
	if err != nil {
		return err
	}

	return nil
}
