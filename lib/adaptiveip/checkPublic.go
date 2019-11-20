package adaptiveip

import (
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/iputil"
	"net"
)

// CheckPublicNetwork : Check if configured public network address is available by provided interface
func CheckPublicNetwork(iface net.Interface) error {
	var publicNetworkOk = false

	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	netIPnetworkIP, mask, err := iputil.CheckNetwork(config.AdaptiveIP.PublicNetworkAddress,
		config.AdaptiveIP.PublicNetworkNetmask)
	if err != nil {
		return err
	}

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
		}
	}

	if !publicNetworkOk {
		return errors.New("configured public network address is not available for provided iface")
	}

	return nil
}
