package adaptiveip

import (
	"errors"
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

	adaptiveIP := GetAdaptiveIPNetwork()

	netNetwork, err := iputil.CheckNetwork(adaptiveIP.ExtIfaceIPAddress, adaptiveIP.Netmask)
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if netNetwork.Contains(ip) {
			publicNetworkOk = true
		}
	}

	if !publicNetworkOk {
		return errors.New("configured public network address is not available for provided iface")
	}

	return nil
}
