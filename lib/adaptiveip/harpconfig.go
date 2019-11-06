package adaptiveip

import (
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/iputil"
	"net"
)

func checkHarpConfigNetwork() error {
	err, netIPnetworkIP, mask := iputil.CheckNetwork(config.AdaptiveIP.PublicNetworkAddress,
		config.AdaptiveIP.PublicNetworkNetmask)
	if err != nil {
		return err
	}

	netStartIP := iputil.CheckValidIP(config.AdaptiveIP.PublicStartIP)
	if netStartIP == nil {
		return errors.New("wrong public start IP address")
	}

	netEndIP := iputil.CheckValidIP(config.AdaptiveIP.PublicEndIP)
	if netEndIP == nil {
		return errors.New("wrong public end IP address")
	}

	ipNet := net.IPNet{
		IP:   netIPnetworkIP,
		Mask: mask,
	}

	isStartIPContainedInNetwork := ipNet.Contains(netStartIP)
	if isStartIPContainedInNetwork == false {
		return errors.New("start IP address is not in the public network address")
	}

	isEndIPContainedInNetwork := ipNet.Contains(netEndIP)
	if isEndIPContainedInNetwork == false {
		return errors.New("end IP address is not in the public network address")
	}

	err, totalAvailableIPs := iputil.GetTotalAvailableIPs(config.AdaptiveIP.PublicNetworkAddress,
		config.AdaptiveIP.PublicNetworkNetmask)
	if err != nil {
		return err
	}

	err, ipRangeCount := iputil.GetIPRangeCount(netStartIP, netEndIP)
	if err != nil {
		return err
	}

	if ipRangeCount > totalAvailableIPs {
		return errors.New("IP range count is bigger than total available IPs")
	}

	return nil
}
