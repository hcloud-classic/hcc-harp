package syscheck

import (
	"fmt"
	"net"
)

// CheckIfaceExist : Check if given network interface name is exist in local system.
func CheckIfaceExist(ifaceName string) bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, iface := range interfaces {
		if iface.Name == ifaceName {
			fmt.Println("checkIfaceExist: '" + ifaceName + "' interface found.")
			return true
		}
	}

	fmt.Println("checkIfaceExist: '" + ifaceName + "' interface not found. Please check the configuration file.")
	return false
}
