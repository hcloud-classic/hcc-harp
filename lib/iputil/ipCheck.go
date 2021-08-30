package iputil

import (
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"net"
	"strconv"
	"strings"
)

// CheckValidIP : Parses string value of IPv4 address then return as net.IP.
// If given wrong IP address, it wil return nil.
func CheckValidIP(ip string) net.IP {
	netIP := net.ParseIP(ip).To4()
	return netIP
}

// CheckNetmask : Check string value of IPv4 netmask then return as net.IPMask.
// If given wrong netmask, it will return nil and error.
func CheckNetmask(netmask string) (net.IPMask, error) {
	var err error

	var maskPartsStr = strings.Split(netmask, ".")
	if len(maskPartsStr) != 4 {
		return nil, errors.New("netmask should be X.X.X.X form")
	}

	var maskParts [4]int
	for i := range maskPartsStr {
		maskParts[i], err = strconv.Atoi(maskPartsStr[i])
		if err != nil {
			return nil, errors.New("netmask contained non-integer value")
		}
	}

	var mask = net.IPv4Mask(
		byte(maskParts[0]),
		byte(maskParts[1]),
		byte(maskParts[2]),
		byte(maskParts[3]))

	maskSizeOne, maskSizeBit := mask.Size()
	if maskSizeOne == 0 && maskSizeBit == 0 {
		return nil, errors.New("invalid netmask")
	}

	if maskSizeOne > 30 {
		return nil, errors.New("netmask bit should be equal or smaller than 30")
	}

	return mask, err
}

// CheckIPisInSubnet : Check if provided IP address is in the given subnet.
// Subnet must be given as net.IPNet and IP must be given as string value.
// It will return error if given invalid IP address or IP is not in the subnet.
func CheckIPisInSubnet(subnet net.IPNet, IP string) error {
	netIP := CheckValidIP(IP)
	if netIP == nil {
		return errors.New("wrong IP address")
	}
	IPisInSubnet := subnet.Contains(netIP)
	if IPisInSubnet == false {
		return errors.New("given IP address is not in the subnet")
	}

	maskLen, _ := subnet.Mask.Size()
	_, netNetwork, _ := net.ParseCIDR(IP + "/" + strconv.Itoa(maskLen))
	firstIP, lastIP := cidr.AddressRange(netNetwork)
	if IP == firstIP.To4().String() {
		return errors.New("you can't use network address for IP address")
	}

	if IP == lastIP.To4().String() {
		return errors.New("you can't use broadcast address for IP address")
	}

	return nil
}

// CheckSubnetIsUsedByIface : Check if given subnet is used in one of iface
func CheckSubnetIsUsedByIface(subnet net.IPNet) error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			var netIP net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				netIP = v.IP
			case *net.IPAddr:
				netIP = v.IP
			}

			IPisInSubnet := subnet.Contains(netIP)
			if IPisInSubnet {
				return errors.New("Given subnet is conflicted with a network interface (" + iface.Name + ")")
			}
		}
	}

	return nil
}
