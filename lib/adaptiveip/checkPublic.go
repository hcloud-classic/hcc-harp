package adaptiveip

import (
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"net"
)

var localPublicIP string

func CheckPublicNetwork(iface net.Interface) error {
	var publicNetworkOk = false

	addrs, err := iface.Addrs()
	if err != nil {
		return err
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
				localPublicIP = ip.String()
			}
			break
		}
	}

	if !publicNetworkOk {
		return errors.New("configured public network address is not available for provided iface")
	}

	return nil
}
