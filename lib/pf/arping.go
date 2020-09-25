package pf

import (
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/mdlayher/arp"
	pb "hcc/harp/action/grpc/pb/rpcharp"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/syscheck"
	"log"
	"net"
	"sync"
	"time"
)

// checkDuplicatedIPAddress : Check duplicated ip address by sending arping.
func checkDuplicatedIPAddress(IP string) error {
	// Ensure valid network interface
	ifi, err := net.InterfaceByName(config.AdaptiveIP.ExternalIfaceName)
	if err != nil {
		return err
	}

	// Set up ARP client with socket
	c, err := arp.Dial(ifi)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = c.Close()
	}()

	// Set request deadline from flag
	if err := c.SetDeadline(time.Now().Add(1 * time.Second)); err != nil {
		return err
	}

	// Request hardware address for IP address
	ip := net.ParseIP(IP).To4()
	mac, err := c.Resolve(ip)
	if err != nil {
		return nil
	}

	err = errors.New("checkDuplicatedIPAddress(): Found duplicated IP address for " + IP + "(MAC: " + mac.String() + ")")
	logger.Logger.Println(err.Error())
	return err
}

// true : available, false : not available
func getAvailableIPsStatusMap() map[string]bool {
	logger.Logger.Println("Getting available IPs status... (This may take a while.)")
	ipMap := make(map[string]bool)

	adaptiveip := configext.GetAdaptiveIPNetwork()
	netStartIP := iputil.CheckValidIP(adaptiveip.StartIPAddress)
	netEndIP := iputil.CheckValidIP(adaptiveip.EndIPAddress)
	ipRangeCount, _ := iputil.GetIPRangeCount(netStartIP, netEndIP)

	var RoutineMAX = int(config.AdaptiveIP.ArpingRoutineMaxNum)
	if RoutineMAX == 0 {
		RoutineMAX = 5
	}
	var routineMax = RoutineMAX
	var wait sync.WaitGroup
	var mutex = &sync.Mutex{}

	for i := 0; i < ipRangeCount; {
		if ipRangeCount-i < RoutineMAX {
			routineMax = ipRangeCount - i
		}

		wait.Add(routineMax)

		for j := 0; j < routineMax; j++ {
			go func(ip string) {
				err := checkDuplicatedIPAddress(ip)
				// Write to map need a lock
				mutex.Lock()
				if err != nil {
					ipMap[ip] = false
				} else {
					ipMap[ip] = true
				}
				mutex.Unlock()

				wait.Done()
			}(netStartIP.String())

			netStartIP = cidr.Inc(netStartIP)

			i++
			if i == ipRangeCount {
				break
			}
		}

		wait.Wait()
	}

	return ipMap
}

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

	ipMap := getAvailableIPsStatusMap()

	for i := 0; i < ipRangeCount; i++ {
		ip := netStartIP.String()
		var ipUsed = false
		if CheckBinatAnchorFileExist(ip) == nil && ipMap[ip] {
			for _, addr := range extIPaddrs {
				var extIP net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					extIP = v.IP
				case *net.IPAddr:
					extIP = v.IP
				}

				if extIP.String() == ip {
					logger.Logger.Println("GetAvailableIPList: " + ip + " is already used in external interface.")
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
