package adaptiveip

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/ifconfig"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/pf"
	"hcc/harp/lib/serviceControl"
	"os/exec"
)

func settingExternalInterface() error {
	logger.Logger.Println("Setting external interface...")

	cmd := exec.Command("/usr/bin/sed", "-i", "", "/ifconfig_"+
		config.AdaptiveIP.ExternalIfaceName+"/d", "/etc/rc.conf")
	err := cmd.Run()
	if err != nil {
		return err
	}

	adaptiveip := config.GetAdaptiveIPNetwork()

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

	adaptiveip := config.GetAdaptiveIPNetwork()

	defaultrouteString := "defaultrouter=\"" + adaptiveip.GatewayAddress + "\"\n"
	err = fileutil.WriteFileAppend("/etc/rc.conf", defaultrouteString)
	if err != nil {
		return err
	}

	return nil
}

func settingExternalNetwork() error {
	logger.Logger.Println("Setting external network...")

	err := settingExternalInterface()
	if err != nil {
		logger.Logger.Println(err)
	}

	err = settingDefaultGateway()
	if err != nil {
		logger.Logger.Println(err)
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

	err = serviceControl.RestartNetwork()
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

