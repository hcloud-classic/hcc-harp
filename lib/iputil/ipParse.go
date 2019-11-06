package iputil

import (
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"net"
)

// CheckNetwork : Get network IP address and netmask as string value then check if valid.
// Return network address as net.IP and netmask as net.IPMask if valid.
func CheckNetwork(networkIP string, networkNetmask string) (net.IP, net.IPMask, error){
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

// GetTotalAvailableIPs : Return total available IPs count for given network IP address and netmask.
func GetTotalAvailableIPs(networkIP string, networkNetmask string) (int, error) {
	firstIP, lastIP, err := GetFirstAndLastIPs(networkIP, networkNetmask)
	if err != nil {
		return 0, err
	}

	firstIPsum := int(firstIP[0]) + int(firstIP[1]) + int(firstIP[2]) + int(firstIP[3])
	lastIPsum := int(lastIP[0]) + int(lastIP[1]) + int(lastIP[2]) + int(lastIP[3])

	totalAvailableIPs := lastIPsum - firstIPsum + 1

	return totalAvailableIPs, nil
}

// GetIPRangeCount : Calculate IPs count from given start IP address and end IP address.
func GetIPRangeCount(startIP net.IP, endIP net.IP) (int, error) {
	startIPsum := int(startIP[0]) + int(startIP[1]) + int(startIP[2]) + int(startIP[3])
	endIPsum := int(endIP[0]) + int(endIP[1]) + int(endIP[2]) + int(endIP[3])

	if startIPsum > endIPsum {
		return 0, errors.New("startIPsum is bigger than endIPsum")
	}

	ipRangeCount := endIPsum - startIPsum + 1

	return ipRangeCount, nil
}
