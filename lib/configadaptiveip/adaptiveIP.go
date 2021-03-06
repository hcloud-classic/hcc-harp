package configadaptiveip

import (
	"errors"
	"hcc/harp/lib/iputil"
	"innogrid.com/hcloud-classic/pb"
	"net"
)

// CheckAdaptiveIPConfig : Check configuration of Adaptive IP
func CheckAdaptiveIPConfig(adaptiveIP *pb.AdaptiveIPSetting) error {
	netNetwork, err := iputil.CheckNetwork(adaptiveIP.ExtIfaceIPAddress,
		adaptiveIP.Netmask)
	if err != nil {
		return err
	}

	err = iputil.CheckIPisInSubnet(*netNetwork, adaptiveIP.GatewayAddress)
	if err != nil {
		return err
	}

	netStartIP := iputil.CheckValidIP(adaptiveIP.StartIPAddress)
	if netStartIP == nil {
		return errors.New("wrong public start IP address")
	}

	netEndIP := iputil.CheckValidIP(adaptiveIP.EndIPAddress)
	if netEndIP == nil {
		return errors.New("wrong public end IP address")
	}

	isStartIPContainedInNetwork := netNetwork.Contains(netStartIP)
	if isStartIPContainedInNetwork == false {
		return errors.New("start IP address is not in the public network address")
	}

	isEndIPContainedInNetwork := netNetwork.Contains(netEndIP)
	if isEndIPContainedInNetwork == false {
		return errors.New("end IP address is not in the public network address")
	}

	totalAvailableIPs, err := iputil.GetTotalAvailableIPs(netNetwork.IP.String(), net.IP(netNetwork.Mask).String())
	if err != nil {
		return err
	}

	ipRangeCount, err := iputil.GetIPRangeCount(netStartIP, netEndIP)
	if err != nil {
		return err
	}

	if ipRangeCount > totalAvailableIPs {
		return errors.New("IP range count is bigger than total available IPs")
	}

	return nil
}
