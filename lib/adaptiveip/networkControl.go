package adaptiveip

import (
	"hcc/harp/action/grpc/client"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/ifconfig"
	"hcc/harp/lib/iptables"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/pf"
	"hcc/harp/lib/servicecontrol"
	"hcc/harp/lib/syscheck"
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

func settingExternalInterfaceFreeBSD() error {
	logger.Logger.Println("Setting external interface...")

	cmd := exec.Command("/usr/bin/sed", "-i", "", "/ifconfig_"+
		config.AdaptiveIP.ExternalIfaceName+"/d", "/etc/rc.conf")
	err := cmd.Run()
	if err != nil {
		return err
	}

	adaptiveip := configext.GetAdaptiveIPNetwork()

	externalInterfaceString := "ifconfig_" + config.AdaptiveIP.ExternalIfaceName +
		"=\"inet " + adaptiveip.ExtIfaceIPAddress + " netmask " + adaptiveip.Netmask + "\"\n"
	err = fileutil.WriteFileAppend("/etc/rc.conf", externalInterfaceString)
	if err != nil {
		return err
	}

	return nil
}

func settingExternalInterfaceLinux() error {
	logger.Logger.Println("Setting external interface...")

	adaptiveip := configext.GetAdaptiveIPNetwork()

	cmd := exec.Command("ifconfig", config.AdaptiveIP.ExternalIfaceName, adaptiveip.ExtIfaceIPAddress, "netmask",
		adaptiveip.Netmask)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func settingExternalInterface() error {
	if syscheck.OS == "freebsd" {
		return settingExternalInterfaceFreeBSD()
	}

	return settingExternalInterfaceLinux()
}

func settingDefaultGatewayFreeBSD() error {
	logger.Logger.Println("Setting default gateway...")

	cmd := exec.Command("/usr/bin/sed", "-i", "", "/defaultrouter/d", "/etc/rc.conf")
	err := cmd.Run()
	if err != nil {
		return err
	}

	adaptiveip := configext.GetAdaptiveIPNetwork()

	defaultrouteString := "defaultrouter=\"" + adaptiveip.GatewayAddress + "\"\n"
	err = fileutil.WriteFileAppend("/etc/rc.conf", defaultrouteString)
	if err != nil {
		return err
	}

	return nil
}

func flushDefaultGatewaysLinux() error {
	logger.Logger.Println("Flushing default gateways...")

	cmd := exec.Command("ip", "route", "flush", "0/0")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func settingDefaultGatewayLinux() error {
	err := flushDefaultGatewaysLinux()
	if err != nil {
		return err
	}

	logger.Logger.Println("Setting default gateway...")

	adaptiveip := configext.GetAdaptiveIPNetwork()

	cmd := exec.Command("route", "add", "default", "gw", adaptiveip.GatewayAddress)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func settingDefaultGateway() error {
	if syscheck.OS == "freebsd" {
		return settingDefaultGatewayFreeBSD()
	}

	return settingDefaultGatewayLinux()
}

func settingExternalNetwork() error {
	logger.Logger.Println("Setting external network...")

	var needNetworkRestart = false

	ifaceName := config.AdaptiveIP.ExternalIfaceName
	adaptiveip := configext.GetAdaptiveIPNetwork()

	isIPConfigured, err := checkIPConfigured(ifaceName, adaptiveip.ExtIfaceIPAddress)
	if err != nil {
		logger.Logger.Println(err)
	}

	if !isIPConfigured {
		err = settingExternalInterface()
		if err != nil {
			logger.Logger.Println(err)
		}
		needNetworkRestart = true
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
		needNetworkRestart = true
	}

	if syscheck.OS == "freebsd" && needNetworkRestart {
		client.End()

		err = servicecontrol.RestartNetwork()
		if err != nil {
			return err
		}

		err = client.Init()
		if err != nil {
			panic(err)
		}
	}

	return nil
}

// LoadHarpPFRules : Load pf rules for harp module
func LoadHarpPFRules() error {
	err := settingExternalNetwork()
	if err != nil {
		return err
	}

	err = dhcpd.CheckDatabaseAndGenerateDHCPDConfigs()
	if err != nil {
		return err
	}

	err = ifconfig.LoadExistingIfconfigScriptsInternal()
	if err != nil {
		return err
	}

	err = pf.FlushPFRules()
	if err != nil {
		return err
	}

	err = pf.LoadPFRules(config.AdaptiveIP.PFRulesFileLocation)
	if err != nil {
		return err
	}

	err = pf.LoadExstingBinatAndNATRules()
	if err != nil {
		return err
	}

	err = ifconfig.LoadExistingIfconfigScriptsExternal()
	if err != nil {
		return err
	}

	return nil
}

// LoadHarpIPTABLESRules : Load iptables rules for harp module
func LoadHarpIPTABLESRules() error {
	err := settingExternalNetwork()
	if err != nil {
		return err
	}

	err = dhcpd.CheckDatabaseAndGenerateDHCPDConfigs()
	if err != nil {
		return err
	}

	err = ifconfig.LoadExistingIfconfigScriptsInternal()
	if err != nil {
		return err
	}

	err = ifconfig.LoadExistingIfconfigScriptsExternal()
	if err != nil {
		return err
	}

	err = iptables.InitIPTABLES()
	if err != nil {
		return err
	}

	err = iptables.LoadAdaptiveIPIPTABLESRules()
	if err != nil {
		return err
	}

	err = iptables.EnableAllRouteLocal()
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

	if syscheck.OS == "freebsd" {
		err = pf.PreparePFConfigFiles()
		if err != nil {
			goto ERROR
		}

		err = LoadHarpPFRules()
		if err != nil {
			goto ERROR
		}
	} else {
		err = LoadHarpIPTABLESRules()
		if err != nil {
			goto ERROR
		}
	}

ERROR:
	firewallLoadLock.Unlock()

	return err
}
