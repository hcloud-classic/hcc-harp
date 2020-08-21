package adaptiveip

import (
	"hcc/harp/action/grpc/client"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/ifconfig"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/pf"
	"hcc/harp/lib/servicecontrol"
	"net"
	"os/exec"
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

func settingDefaultGateway() error {
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

	if needNetworkRestart {
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
