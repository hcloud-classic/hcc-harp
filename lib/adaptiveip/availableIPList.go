package adaptiveip

import (
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"hcc/harp/lib/arping"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configadapriveipnetwork"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/syscheck"
	"innogrid.com/hcloud-classic/pb"
	"net"
)

// GetAvailableIPList : Get available IP lists by checking config files and sending arping.
func GetAvailableIPList() (*pb.AdaptiveIPAvailableIPList, error) {
	var availableIPList pb.AdaptiveIPAvailableIPList
	var availableIPs []string

	adaptiveip := configadapriveipnetwork.GetAdaptiveIPNetwork()

	internalNetStartIP := iputil.CheckValidIP(adaptiveip.InternalStartIPAddress)
	internalNetEndIP := iputil.CheckValidIP(adaptiveip.InternalEndIPAddress)
	internalIPRangeCount, _ := iputil.GetIPRangeCount(internalNetStartIP, internalNetEndIP)

	externalNetStartIP := iputil.CheckValidIP(adaptiveip.ExternalStartIPAddress)
	externalNetEndIP := iputil.CheckValidIP(adaptiveip.ExternalEndIPAddress)
	externalIPRangeCount, _ := iputil.GetIPRangeCount(externalNetStartIP, externalNetEndIP)

	if internalIPRangeCount != externalIPRangeCount {
		return nil, errors.New("external IP range is not match with internal IP range")
	}

	extIface, _ := syscheck.CheckIfaceExist(config.AdaptiveIP.ExternalIfaceName)
	extIPaddrs, err := extIface.Addrs()
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	ipMap, _ := arping.GetAvailableIPsStatusMap(internalNetStartIP, internalNetEndIP)

	for i := 0; i < internalIPRangeCount; i++ {
		internalIP := internalNetStartIP.To4().String()
		externalIP := externalNetStartIP.To4().String()

		var ipUsed = false
		if ipMap[internalIP] {
			for _, addr := range extIPaddrs {
				var extIP net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					extIP = v.IP
				case *net.IPAddr:
					extIP = v.IP
				}

				if extIP.String() == internalIP {
					logger.Logger.Println("GetAvailableIPList(): Internal IP address (" +
						internalIP + ") is already used in external interface.")
					ipUsed = true
					break
				}
			}

			if !ipUsed {
				availableIPs = append(availableIPs, externalIP)
			}
		}

		internalNetStartIP = cidr.Inc(internalNetStartIP)
		externalNetStartIP = cidr.Inc(externalNetStartIP)
	}

	availableIPList.AvailableIp = availableIPs

	return &availableIPList, nil
}
