package adaptiveip

import (
	"github.com/apparentlymart/go-cidr/cidr"
	"hcc/harp/lib/arping"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
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

	adaptiveip := configext.GetAdaptiveIPNetwork()
	netStartIP := iputil.CheckValidIP(adaptiveip.StartIPAddress)
	netEndIP := iputil.CheckValidIP(adaptiveip.EndIPAddress)
	ipRangeCount, _ := iputil.GetIPRangeCount(netStartIP, netEndIP)

	extIface, _ := syscheck.CheckIfaceExist(config.AdaptiveIP.ExternalIfaceName)
	extIPaddrs, err := extIface.Addrs()
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	ipMap, _ := arping.GetAvailableIPsStatusMap(netStartIP, netEndIP)

	for i := 0; i < ipRangeCount; i++ {
		ip := netStartIP.String()
		var ipUsed = false
		if ipMap[ip] {
			for _, addr := range extIPaddrs {
				var extIP net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					extIP = v.IP
				case *net.IPAddr:
					extIP = v.IP
				}

				if extIP.String() == ip {
					logger.Logger.Println("GetAvailableIPList(): " + ip + " is already used in external interface.")
					ipUsed = true
					break
				}
			}

			if !ipUsed {
				availableIPs = append(availableIPs, ip)
			}
		}

		netStartIP = cidr.Inc(netStartIP)
	}

	availableIPList.AvailableIp = availableIPs

	return &availableIPList, nil
}
