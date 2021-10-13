package adaptiveipext

import (
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"hcc/harp/lib/configadapriveipnetwork"
	"hcc/harp/lib/iputil"
	"strconv"
	"strings"
)

// ExternalIPtoInternalIP : Change AdaptiveIP external IP address to internal IP address
func ExternalIPtoInternalIP(externalIP string) (string, error) {
	adaptiveip := configadapriveipnetwork.GetAdaptiveIPNetwork()

	internalNetStartIP := iputil.CheckValidIP(adaptiveip.InternalStartIPAddress)
	internalNetEndIP := iputil.CheckValidIP(adaptiveip.InternalEndIPAddress)
	internalIPRangeCount, _ := iputil.GetIPRangeCount(internalNetStartIP, internalNetEndIP)

	externalNetStartIP := iputil.CheckValidIP(adaptiveip.ExternalStartIPAddress)
	externalNetEndIP := iputil.CheckValidIP(adaptiveip.ExternalEndIPAddress)
	externalIPRangeCount, _ := iputil.GetIPRangeCount(externalNetStartIP, externalNetEndIP)

	if internalIPRangeCount != externalIPRangeCount {
		return "", errors.New("external IP range is not match with internal IP range")
	}

	var startIPSum = 0
	var endIPSsum = 0
	var externalIPSum = 0

	startIPSplit := strings.Split(externalNetStartIP.To4().String(), ".")
	endIPSplit := strings.Split(externalNetEndIP.To4().String(), ".")
	externalIPSplit := strings.Split(externalIP, ".")

	for _, startIPSplited := range startIPSplit {
		num, _ := strconv.Atoi(startIPSplited)
		startIPSum += num
	}
	for _, endIPSplited := range endIPSplit {
		num, _ := strconv.Atoi(endIPSplited)
		endIPSsum += num
	}
	for _, externalIPSplited := range externalIPSplit {
		num, _ := strconv.Atoi(externalIPSplited)
		externalIPSum += num
	}

	if externalIPSum < startIPSum || externalIPSum > endIPSsum {
		return "", errors.New("external IP address is out of range")
	}

	for i := 0; i < externalIPRangeCount; i++ {
		if externalNetStartIP.To4().String() == externalIP {
			break
		}

		internalNetStartIP = cidr.Inc(internalNetStartIP)
		externalNetStartIP = cidr.Inc(externalNetStartIP)
	}

	return internalNetStartIP.To4().String(), nil
}

// InternalIPtoExternalIP : Change AdaptiveIP internal IP address to external IP address
func InternalIPtoExternalIP(internalIP string) (string, error) {
	adaptiveip := configadapriveipnetwork.GetAdaptiveIPNetwork()

	internalNetStartIP := iputil.CheckValidIP(adaptiveip.InternalStartIPAddress)
	internalNetEndIP := iputil.CheckValidIP(adaptiveip.InternalEndIPAddress)
	internalIPRangeCount, _ := iputil.GetIPRangeCount(internalNetStartIP, internalNetEndIP)

	externalNetStartIP := iputil.CheckValidIP(adaptiveip.ExternalStartIPAddress)
	externalNetEndIP := iputil.CheckValidIP(adaptiveip.ExternalEndIPAddress)
	externalIPRangeCount, _ := iputil.GetIPRangeCount(externalNetStartIP, externalNetEndIP)

	if internalIPRangeCount != externalIPRangeCount {
		return "", errors.New("external IP range is not match with internal IP range")
	}

	var startIPSum = 0
	var endIPSsum = 0
	var internalIPSum = 0

	startIPSplit := strings.Split(internalNetStartIP.To4().String(), ".")
	endIPSplit := strings.Split(internalNetEndIP.To4().String(), ".")
	internalIPSplit := strings.Split(internalIP, ".")

	for _, startIPSplited := range startIPSplit {
		num, _ := strconv.Atoi(startIPSplited)
		startIPSum += num
	}
	for _, endIPSplited := range endIPSplit {
		num, _ := strconv.Atoi(endIPSplited)
		endIPSsum += num
	}
	for _, internalIPSplited := range internalIPSplit {
		num, _ := strconv.Atoi(internalIPSplited)
		internalIPSum += num
	}

	if internalIPSum < startIPSum || internalIPSum > endIPSsum {
		return "", errors.New("internal IP address is out of range")
	}

	for i := 0; i < internalIPRangeCount; i++ {
		if internalNetStartIP.To4().String() == internalIP {
			break
		}

		internalNetStartIP = cidr.Inc(internalNetStartIP)
		externalNetStartIP = cidr.Inc(externalNetStartIP)
	}

	return externalNetStartIP.To4().String(), nil
}
