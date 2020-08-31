package pf

import (
	"bytes"
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	pb "hcc/harp/action/grpc/pb/rpcharp"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/syscheck"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// checkDuplicatedIPAddress : Check duplicated ip address by sending arping.
func checkDuplicatedIPAddress(IP string) error {
	cmd := exec.Command("arping", "-i", config.AdaptiveIP.ExternalIfaceName, "-c",
		strconv.Itoa(int(config.AdaptiveIP.ArpingRetryCount)), IP)

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	cmdOutputStr := string(cmdOutput.Bytes())

	if strings.Contains(cmdOutputStr, "Timeout") ||
		strings.Contains(cmdOutputStr, "timeout") {
		return nil
	}

	if err != nil {
		logger.Logger.Println("arping: " + err.Error())
	}

	if strings.Contains(cmdOutputStr, "from") {
		err := errors.New("checkDuplicatedIPAddress(): Found duplicated IP address for " + IP)
		logger.Logger.Println(err.Error())
		return err
	}

	return nil
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
func GetAvailableIPList() *pb.AdaptiveIPAvailableIPList {
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
		return &availableIPList
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

	return &availableIPList
}
