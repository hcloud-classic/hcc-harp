package syscheck

import (
	"hcc/harp/lib/logger"
	"net"
)

func CheckIfaceExist(ifaceName string) bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Logger.Println(err)
		return false
	}

	for _, iface := range interfaces {
		if iface.Name == ifaceName {
			logger.Logger.Println("checkIfaceExist: '" + ifaceName + "' interface found.")
			return true
		}
	}

	logger.Logger.Println("checkIfaceExist: '" + ifaceName + "' interface not found. Please check the configuration file.")
	return false
}
