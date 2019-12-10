package adaptiveip

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/logger"
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

	adaptiveip := GetAdaptiveIPNetwork()

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

	adaptiveip := GetAdaptiveIPNetwork()

	defaultrouteString := "defaultrouter=\"" + adaptiveip.GatewayAddress + "\"\n"
	err = fileutil.WriteFileAppend("/etc/rc.conf", defaultrouteString)
	if err != nil {
		return err
	}

	return nil
}

func settingExternalNetwork() error {
	logger.Logger.Println("Setting externel network...")

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

func restartNetif() error {
	logger.Logger.Println("Restarting netif service...")

	cmd := exec.Command("service", "netif", "restart")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func restartRouting() error {
	logger.Logger.Println("Restarting routing service...")

	cmd := exec.Command("service", "routing", "restart")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func restartNetwork() error {
	logger.Logger.Println("Restarting network services...")

	err := restartNetif()
	if err != nil {
		logger.Logger.Println(err)
	}

	err = restartRouting()
	if err != nil {
		logger.Logger.Println(err)
	}

	return nil
}
