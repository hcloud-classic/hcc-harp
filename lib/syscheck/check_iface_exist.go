package syscheck

import (
	"fmt"
	"net"
)

// CheckIfaceExist : Check if given network interface name is exist in local system.
func CheckIfaceExist(ifaceName string) (bool, net.Interface) {
	var iface net.Interface

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return false, iface
	}

	for _, i := range interfaces {
		if i.Name == ifaceName {
			fmt.Println("checkIfaceExist: '" + ifaceName + "' interface found.")
			return true, i
		}
	}

	fmt.Println("checkIfaceExist: '" + ifaceName + "' interface not found. Please check the configuration file.")
	return false, iface
}
