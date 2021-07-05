package iptables

import (
	"errors"
	"hcc/harp/lib/configadapriveipnetwork"
	"hcc/harp/lib/iptablesext"
	"hcc/harp/lib/logger"
	"os/exec"
	"strings"
)

func allowInputPrivateClass(class string) error {
	privateNetwork := ""

	class = strings.ToUpper(class)

	switch class {
	case "A":
		privateNetwork = "10.0.0.0/8"
	case "B":
		privateNetwork = "172.16.0.0/12"
	case "C":
		privateNetwork = "192.168.0.0/16"
	}

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()

	cmd := exec.Command("iptables", "-t", "filter",
		"-C", iptablesext.HarpChainNamePrefix+"INPUT",
		"-s", privateNetwork,
		"-d", adaptiveIP.ExtIfaceIPAddress,
		"-j", "ACCEPT")
	err := cmd.Run()
	isExist := err == nil

	if !isExist {
		cmd = exec.Command("iptables", "-t", "filter",
			"-I", iptablesext.HarpChainNamePrefix+"INPUT", "1",
			"-s", privateNetwork,
			"-d", adaptiveIP.ExtIfaceIPAddress,
			"-j", "ACCEPT")
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to add " + class + " class private network accept rule of the Master Node")
		}
	}

	return nil
}

func allowInputPrivateNetworks() error {
	logger.Logger.Println("Adding allow rules of private networks for the Master Node...")

	err := allowInputPrivateClass("A"); if err != nil {
		return err
	}

	err = allowInputPrivateClass("B"); if err != nil {
		return err
	}

	err = allowInputPrivateClass("C"); if err != nil {
		return err
	}

	return nil
}

func allowInputPingReply() error {
	logger.Logger.Println("Adding allow rules of ping reply for the Master Node...")

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()

	cmd := exec.Command("iptables", "-t", "filter",
		"-C", iptablesext.HarpChainNamePrefix+"INPUT",
		"-p", "icmp", "--icmp-type", "echo-reply",
		"-d", adaptiveIP.ExtIfaceIPAddress,
		"-j", "ACCEPT")
	err := cmd.Run()
	isExist := err == nil

	if !isExist {
		cmd = exec.Command("iptables", "-t", "filter",
			"-I", iptablesext.HarpChainNamePrefix+"INPUT", "1",
			"-p", "icmp", "--icmp-type", "echo-reply",
			"-d", adaptiveIP.ExtIfaceIPAddress,
			"-j", "ACCEPT")
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to add ping request accept rule of the Master Node")
		}
	}

	return nil
}

func allowInputEstablishedRelated() error {
	logger.Logger.Println("Adding allow rules of state for the Master Node...")

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()
	protocol := []string{"tcp", "udp"}

	for i := range protocol {
		cmd := exec.Command("iptables", "-t", "filter",
			"-C", iptablesext.HarpChainNamePrefix+"INPUT",
			"-p", protocol[i], "-m", "state", "--state", "ESTABLISHED,RELATED",
			"-d", adaptiveIP.ExtIfaceIPAddress,
			"-j", "ACCEPT")
		err := cmd.Run()
		isExist := err == nil

		if !isExist {
			cmd = exec.Command("iptables", "-t", "filter",
				"-I", iptablesext.HarpChainNamePrefix+"INPUT", "1",
				"-p", protocol[i], "-m", "state", "--state", "ESTABLISHED,RELATED",
				"-d", adaptiveIP.ExtIfaceIPAddress,
				"-j", "ACCEPT")
			err = cmd.Run()
			if err != nil {
				return errors.New("failed to add " + strings.ToUpper(protocol[i]) + " state rule of the Master Node")
			}
		}
	}

	return nil
}

func addMasterDrop() error {
	logger.Logger.Println("Adding Master Node's DROP rule...")

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()

	cmd := exec.Command("iptables", "-t", "filter",
		"-C", iptablesext.HarpAdaptiveIPInputDropChainName,
		"-d", adaptiveIP.ExtIfaceIPAddress,
		"-j", "DROP")
	err := cmd.Run()
	isExist := err == nil

	if !isExist {
		cmd = exec.Command("iptables", "-t", "filter",
			"-A", iptablesext.HarpAdaptiveIPInputDropChainName,
			"-d", adaptiveIP.ExtIfaceIPAddress,
			"-j", "DROP")
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to add ADAPTIVE_IP_INPUT_DROP rule of the Master Node")
		}
	}

	return nil
}

func masterInputControl() error {
	err := allowInputPrivateNetworks(); if err != nil {
		return err
	}

	err = allowInputPingReply(); if err != nil {
		return err
	}

	err = allowInputEstablishedRelated(); if err != nil {
		return err
	}

	err = addMasterDrop(); if err != nil {
		return err
	}

	return nil
}
