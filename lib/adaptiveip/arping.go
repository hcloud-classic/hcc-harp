package adaptiveip

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	"hcc/harp/lib/config"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/model"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

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
		return errors.New("Found duplicated IP address for " + IP)
	}

	return nil
}

// true : available, false : not available
func getAvailableIPsStatusMap() map[string]bool {
	logger.Logger.Println("Getting available IPs status... (This may take a while.)")
	ipMap := make(map[string]bool)

	adaptiveip := GetAdaptiveIPNetwork()
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

func GetAvailableIPList() model.AdaptiveIPAvailableIPList {
	var availableIPList model.AdaptiveIPAvailableIPList

	adaptiveip := GetAdaptiveIPNetwork()
	netStartIP := iputil.CheckValidIP(adaptiveip.StartIPAddress)
	netEndIP := iputil.CheckValidIP(adaptiveip.EndIPAddress)
	ipRangeCount, _ := iputil.GetIPRangeCount(netStartIP, netEndIP)
	fmt.Println("ipRangeCount", ipRangeCount)

	ipMap := getAvailableIPsStatusMap()

	for i := 0; i < ipRangeCount; i++ {
		ip := netStartIP.String()
		if checkBinatAnchorFileExist(ip) == nil && ipMap[ip] {
			availableIPList.AvailableIPList = append(availableIPList.AvailableIPList, ip)
		}

		netStartIP = cidr.Inc(netStartIP)
	}

	return availableIPList
}
