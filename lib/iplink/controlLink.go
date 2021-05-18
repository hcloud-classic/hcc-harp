package iplink

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/iplinkext"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/vnstat"
	"net"
	"strconv"
	"strings"
)

func isHarpInternalDeviceExist(ip string) bool {
	err := runIP("link show " +
		iplinkext.HarpInternalPrefix + strconv.Itoa(iplinkext.GetIfaceVNUM(ip)))
	if err != nil {
		return false
	}

	return true
}

func addHarpInternalDevice(ip string) error {
	err := runIP("link add link " + config.AdaptiveIP.InternalIfaceName + " " +
		iplinkext.HarpInternalPrefix + strconv.Itoa(iplinkext.GetIfaceVNUM(ip)) +
		" address " + generateMACAddress(ip) +
		" type macvlan")
	if err != nil {
		return err
	}

	return nil
}

func deleteHarpInternalDevice(ip string) error {
	err := runIP("link delete " +
		iplinkext.HarpInternalPrefix + strconv.Itoa(iplinkext.GetIfaceVNUM(ip)))
	if err != nil {
		return err
	}

	return nil
}

func upHarpInternalDevice(ip string) error {
	err := runIP("link set dev " +
		iplinkext.HarpInternalPrefix + strconv.Itoa(iplinkext.GetIfaceVNUM(ip)) + " up")
	if err != nil {
		return err
	}

	return nil
}

func downHarpInternalDevice(ip string) error {
	err := runIP("link set dev " +
		iplinkext.HarpInternalPrefix + strconv.Itoa(iplinkext.GetIfaceVNUM(ip)) + " down")
	if err != nil {
		return err
	}

	return nil
}

func setIPtoHarpInternalDevice(ip string, cidr int) error {
	err := runIP("address add " +
		ip + "/" + strconv.Itoa(cidr) +
		" dev " +
		iplinkext.HarpInternalPrefix + strconv.Itoa(iplinkext.GetIfaceVNUM(ip)))
	if err != nil {
		return err
	}

	return nil
}

// SetHarpInternalDevice : Setting up harp internal gateway device
func SetHarpInternalDevice(ip string, netmask string) error {
	var err error
	var mask net.IPMask
	var maskLen int

	if isHarpInternalDeviceExist(ip) {
		goto SCHEDULE
	}

	err = addHarpInternalDevice(ip)
	if err != nil {
		return err
	}

	mask, err = iputil.CheckNetmask(netmask)
	if err != nil {
		return err
	}
	maskLen, _ = mask.Size()
	err = setIPtoHarpInternalDevice(ip, maskLen)
	if err != nil {
		return err
	}

	err = upHarpInternalDevice(ip)
	if err != nil {
		return err
	}

SCHEDULE:
	vnstat.ScheduleUpdateVnStat(iplinkext.HarpInternalPrefix+strconv.Itoa(iplinkext.GetIfaceVNUM(ip)), true)

	return nil
}

// UnsetHarpInternalDevice : Delete harp internal gateway device
func UnsetHarpInternalDevice(ip string) error {
	if !isHarpInternalDeviceExist(ip) {
		return nil
	}

	err := downHarpInternalDevice(ip)
	if err != nil {
		return err
	}

	err = deleteHarpInternalDevice(ip)
	if err != nil {
		return err
	}

	vnstat.RemoveUpdateVnStat(iplinkext.HarpInternalPrefix + strconv.Itoa(iplinkext.GetIfaceVNUM(ip)))

	return nil
}

// AddOrDeleteIPToHarpExternalDevice : Add/Delete AdaptiveIP address to/from external interface
func AddOrDeleteIPToHarpExternalDevice(ip string, netmask string, isAdd bool) error {
	mask, err := iputil.CheckNetmask(netmask)
	if err != nil {
		return err
	}
	maskLen, _ := mask.Size()

	var addArg string
	if isAdd {
		addArg = "add"
	} else {
		addArg = "del"
	}
	err = runIP("address " + addArg + " " +
		ip + "/" + strconv.Itoa(maskLen) +
		" dev " +
		config.AdaptiveIP.ExternalIfaceName)
	if err != nil {
		if strings.Contains(err.Error(), "RTNETLINK") &&
			strings.Contains(err.Error(), "exist") {
			return nil
		}
		return err
	}

	return nil
}
