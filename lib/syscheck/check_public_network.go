package syscheck

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"net"
)

var LocalPublicIP string

func CheckPublicNetwork(iface net.Interface) bool {
	var publicNetworkOk = false

	addrs, err := iface.Addrs()
	if err != nil {
		logger.Logger.Println(err)
		return false
	}

	netIPnetworkIP, mask, err := iputil.CheckNetwork(config.AdaptiveIP.PublicNetworkAddress,
		config.AdaptiveIP.PublicNetworkNetmask)
	ipNet := net.IPNet{
		IP:   netIPnetworkIP,
		Mask: mask,
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ipNet.Contains(ip) {
			publicNetworkOk = true
			if ip != nil {
				LocalPublicIP = ip.String()
			}
			break
		}
	}

	if !publicNetworkOk {
		logger.Logger.Println("Configured public network address is not available for provided iface!")
		return false
	}

	return true
}
