package configadaptiveip

import (
	"errors"
	"hcc/harp/lib/iputil"
	"innogrid.com/hcloud-classic/pb"
	"net"
)

// CheckAdaptiveIPConfig : Check configuration of Adaptive IP
func CheckAdaptiveIPConfig(adaptiveIP *pb.AdaptiveIPSetting) error {
	internalNetNetwork, err := iputil.CheckNetwork(adaptiveIP.ExtIfaceIPAddress,
		adaptiveIP.Netmask)
	if err != nil {
		return err
	}

	err = iputil.CheckIPisInSubnet(*internalNetNetwork, adaptiveIP.GatewayAddress)
	if err != nil {
		return err
	}

	internalNetStartIP := iputil.CheckValidIP(adaptiveIP.InternalStartIPAddress)
	if internalNetStartIP == nil {
		return errors.New("wrong internal start IP address")
	}

	internalNetEndIP := iputil.CheckValidIP(adaptiveIP.InternalEndIPAddress)
	if internalNetEndIP == nil {
		return errors.New("wrong public end IP address")
	}

	isStartIPContainedInNetwork := internalNetNetwork.Contains(internalNetStartIP)
	if isStartIPContainedInNetwork == false {
		return errors.New("start IP address is not in the public network address")
	}

	isEndIPContainedInNetwork := internalNetNetwork.Contains(internalNetEndIP)
	if isEndIPContainedInNetwork == false {
		return errors.New("end IP address is not in the public network address")
	}

	externalNetStartIP := iputil.CheckValidIP(adaptiveIP.ExternalStartIPAddress)
	if internalNetStartIP == nil {
		return errors.New("wrong internal start IP address")
	}

	externalNetEndIP := iputil.CheckValidIP(adaptiveIP.ExternalEndIPAddress)
	if internalNetEndIP == nil {
		return errors.New("wrong public end IP address")
	}

	totalInternalAvailableIPs, err := iputil.GetTotalAvailableIPs(internalNetNetwork.IP.String(), net.IP(internalNetNetwork.Mask).String())
	if err != nil {
		return err
	}

	ipRangeCountInternal, err := iputil.GetIPRangeCount(internalNetStartIP, internalNetEndIP)
	if err != nil {
		return err
	}

	ipRangeCountExternal, err := iputil.GetIPRangeCount(externalNetStartIP, externalNetEndIP)
	if err != nil {
		return err
	}

	if ipRangeCountInternal > totalInternalAvailableIPs {
		return errors.New("internal IP range is bigger than total available IPs")
	}

	if ipRangeCountInternal != ipRangeCountExternal {
		return errors.New("external IP range is not match with internal IP range")
	}

	return nil
}
