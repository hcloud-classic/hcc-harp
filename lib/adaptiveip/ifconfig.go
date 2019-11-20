package adaptiveip

import (
	"hcc/harp/lib/logger"
	"os/exec"
)

func ifconfigAlias(ifaceName string, ip string, netmask string, isAdd bool) error {
	logger.Logger.Println("Adding IP alias for IP="+ip, ", Netmask="+netmask)

	var aliasStr string
	if isAdd {
		aliasStr = "alias"
	} else {
		aliasStr = "-alias"
	}

	cmd := exec.Command("ifconfig", ifaceName, ip, "netmask", netmask, aliasStr)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
