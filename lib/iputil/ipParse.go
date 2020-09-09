package iputil

import (
	"errors"
	"net"
)

// CheckNetwork : Get network IP address and netmask as string value then check if valid.
// Return network address as net.IP and netmask as net.IPMask if valid.
func CheckNetwork(networkIP string, networkNetmask string) (net.IP, net.IPMask, error) {
	netIPnetworkIP := CheckValidIP(networkIP)
	if netIPnetworkIP == nil {
		return nil, nil, errors.New("wrong network IP address")
	}

	mask, err := CheckNetmask(networkNetmask)
	if err != nil {
		return nil, nil, err
	}

	return netIPnetworkIP, mask, nil
}

