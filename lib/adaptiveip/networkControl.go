package adaptiveip

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/configAdapriveIPNetwork"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/iptables"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"net"
	"os/exec"
	"sync"
)

func checkIPConfigured(ifaceName string, ip string) (bool, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return false, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return false, err
	}

	netIP := iputil.CheckValidIP(ip)

	for _, addr := range addrs {
		var ifaceIP net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ifaceIP = v.IP
		case *net.IPAddr:
			ifaceIP = v.IP
		}

		if ifaceIP != nil && ifaceIP.Equal(netIP) {
			return true, nil
		}
	}

	return false, nil
}

func checkGatewayConfigured(gateway string) (bool, error) {
	systemGateway, err := iputil.GetDefaultRoute()
	if err != nil {
		return false, err
	}

	gatewayNetIP := iputil.CheckValidIP(gateway)
	if systemGateway.Equal(gatewayNetIP) {
		return true, nil
	}

	return false, nil
}

func settingExternalInterface() error {
	logger.Logger.Println("Setting external interface...")

	adaptiveip := configAdapriveIPNetwork.GetAdaptiveIPNetwork()

	cmd := exec.Command("ifconfig", config.AdaptiveIP.ExternalIfaceName, adaptiveip.ExtIfaceIPAddress, "netmask",
		adaptiveip.Netmask)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func flushDefaultGateways() error {
	logger.Logger.Println("Flushing default gateways...")

	cmd := exec.Command("ip", "route", "flush", "0/0")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func settingDefaultGateway() error {
	err := flushDefaultGateways()
	if err != nil {
		return err
	}

	logger.Logger.Println("Setting default gateway...")

	adaptiveip := configAdapriveIPNetwork.GetAdaptiveIPNetwork()

	cmd := exec.Command("route", "add", "default", "gw", adaptiveip.GatewayAddress)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func settingExternalNetwork() error {
	logger.Logger.Println("Setting external network...")

	ifaceName := config.AdaptiveIP.ExternalIfaceName
	adaptiveip := configAdapriveIPNetwork.GetAdaptiveIPNetwork()

	isIPConfigured, err := checkIPConfigured(ifaceName, adaptiveip.ExtIfaceIPAddress)
	if err != nil {
		logger.Logger.Println(err)
	}

	if !isIPConfigured {
		err = settingExternalInterface()
		if err != nil {
			logger.Logger.Println(err)
		}
	}

	isGatewayConfigured, err := checkGatewayConfigured(adaptiveip.GatewayAddress)
	if err != nil {
		logger.Logger.Println(err)
	}

	if !isGatewayConfigured {
		err = settingDefaultGateway()
		if err != nil {
			logger.Logger.Println(err)
		}
	}

	return nil
}

// LoadHarpIPTABLESRules : Load iptables rules for harp module
func LoadHarpIPTABLESRules() error {
	err := settingExternalNetwork()
	if err != nil {
		return err
	}

	err = dhcpd.CheckDatabaseAndPrepareDHCPD()
	if err != nil {
		return err
	}

	err = iptables.InitIPTABLES()
	if err != nil {
		return err
	}

	err = iptables.LoadAdaptiveIPNetDevAndIPTABLESRules()
	if err != nil {
		return err
	}

	err = iptables.EnableIPForwardV4()
	if err != nil {
		return err
	}

	return nil
}

var firewallLoadLock sync.Mutex

// LoadFirewall : Load firewall rules for harp module
func LoadFirewall() error {
	var err error = nil

	firewallLoadLock.Lock()

	err = LoadHarpIPTABLESRules()
	if err != nil {
		goto ERROR
	}

ERROR:
	firewallLoadLock.Unlock()

	return err
}
