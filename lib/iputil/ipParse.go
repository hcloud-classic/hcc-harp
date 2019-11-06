package iputil

import (
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"net"
)

func CheckNetwork(networkIP string, networkNetmask string) (error, net.IP, net.IPMask){
	netIPnetworkIP := CheckValidIP(networkIP)
	if netIPnetworkIP == nil {
		return errors.New("wrong network IP address"), nil, nil
	}

	mask, err := CheckNetmask(networkNetmask)
	if err != nil {
		return err, nil, nil
	}

	return nil, netIPnetworkIP, mask
}

func GetFirstAndLastIPs(networkIP string, networkNetmask string) (error, net.IP, net.IP) {
	err, netIPnetworkIP, mask := CheckNetwork(networkIP, networkNetmask)
	if err != nil {
		return err, nil, nil
	}

	ipNet := net.IPNet{
		IP:   netIPnetworkIP,
		Mask: mask,
	}

	firstIP, lastIP := cidr.AddressRange(&ipNet)
	firstIP = cidr.Inc(firstIP)
	lastIP = cidr.Dec(lastIP)

	return nil, firstIP, lastIP
}

func GetTotalAvailableIPs(networkIP string, networkNetmask string) (error, int) {
	err, firstIP, lastIP := GetFirstAndLastIPs(networkIP, networkNetmask)
	if err != nil {
		return err, 0
	}

	firstIPsum := int(firstIP[0]) + int(firstIP[1]) + int(firstIP[2]) + int(firstIP[3])
	lastIPsum := int(lastIP[0]) + int(lastIP[1]) + int(lastIP[2]) + int(lastIP[3])

	totalAvailableIPs := lastIPsum - firstIPsum + 1

	return nil, totalAvailableIPs
}

func GetIPRangeCount(startIP net.IP, endIP net.IP) (error, int) {
	startIPsum := int(startIP[0]) + int(startIP[1]) + int(startIP[2]) + int(startIP[3])
	endIPsum := int(endIP[0]) + int(endIP[1]) + int(endIP[2]) + int(endIP[3])

	if startIPsum > endIPsum {
		return errors.New("startIPsum is bigger than endIPsum"), 0
	}

	ipRangeCount := endIPsum - startIPsum + 1

	return nil, ipRangeCount
}
