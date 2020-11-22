package syscheck

import (
	"errors"
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
			return i, nil
		}
	}

	return iface, errors.New(ifaceName + "' interface not found. Please check the configuration file.")
}
