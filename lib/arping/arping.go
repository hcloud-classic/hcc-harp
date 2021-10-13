package arping

import (
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/mdlayher/arp"
	"hcc/harp/lib/config"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"net"
	"sync"
	"time"
)

// CheckDuplicatedIPAddress : Check duplicated ip address by sending arping.
func CheckDuplicatedIPAddress(IP string) error {
	// Ensure valid network interface
	ifi, err := net.InterfaceByName(config.AdaptiveIP.ExternalIfaceName)
	if err != nil {
		return err
	}

	// Set up ARP client with socket
	c, err := arp.Dial(ifi)
	if err != nil {
		return err
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

	err = errors.New("checkDuplicatedIPAddress(): Found duplicated internal IP address for " + IP + " (MAC: " + mac.String() + ")")
	logger.Logger.Println(err.Error())
	return err
}

// GetAvailableIPsStatusMap : Check duplicated IPs by sending arping and get available IPs in map.
// (true : available, false : not available)
func GetAvailableIPsStatusMap(netStartIP net.IP, netEndIP net.IP) (map[string]bool, error) {
	logger.Logger.Println("Getting available IPs status... (This may take a while.)")
	ipMap := make(map[string]bool)

	ipRangeCount, err := iputil.GetIPRangeCount(netStartIP, netEndIP)
	if err != nil {
		return nil, err
	}

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
			go func(wait *sync.WaitGroup, ip string) {
				err := CheckDuplicatedIPAddress(ip)
				// Write to map need a lock
				mutex.Lock()
				if err != nil {
					ipMap[ip] = false
				} else {
					ipMap[ip] = true
				}
				mutex.Unlock()

				wait.Done()
			}(&wait, netStartIP.String())

			netStartIP = cidr.Inc(netStartIP)

			i++
			if i == ipRangeCount {
				break
			}
		}

		wait.Wait()
	}

	return ipMap, nil
}
