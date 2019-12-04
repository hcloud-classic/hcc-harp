package syscheck

import (
	"errors"
	"fmt"
	"net"
)

// CheckIfaceExist : Check if given network interface name is exist in local system.
func CheckIfaceExist(ifaceName string) (net.Interface, error) {
	var iface net.Interface

	interfaces, err := net.Interfaces()
	if err != nil {
		return iface, err
	}

	for _, i := range interfaces {
		if i.Name == ifaceName {
			fmt.Println("checkIfaceExist: '" + ifaceName + "' interface found.")
			return i, nil
		}
	}

	fmt.Println("checkIfaceExist: '" + ifaceName + "' interface not found. Please check the configuration file.")
	return iface, errors.New("interface not found")
}
