package iputil

import (
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
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

// GetFirstAndLastIPs : Return first and last IP addresses from given network IP address and netmask.
// Return as net.IP for both first and last IP addresses.
func GetFirstAndLastIPs(networkIP string, networkNetmask string) (net.IP, net.IP, error) {
	netIPnetworkIP, mask, err := CheckNetwork(networkIP, networkNetmask)
	if err != nil {
		return nil, nil, err
	}

	ipNet := net.IPNet{
		IP:   netIPnetworkIP,
		Mask: mask,
	}

	firstIP, lastIP := cidr.AddressRange(&ipNet)
	firstIP = cidr.Inc(firstIP)
	lastIP = cidr.Dec(lastIP)

	return firstIP, lastIP, nil
}

