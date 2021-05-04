package ifconfig

import (
	"os/exec"
	"strconv"
	"strings"
)

func runCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getifaceVNUM(ip string) (vnum int) {
	var ifaceVNUM = 0

	ipSplit := strings.Split(ip, ".")
	for _, ipSplited := range ipSplit {
		ipSplitedInt, _ := strconv.Atoi(ipSplited)
		ifaceVNUM += ipSplitedInt
	}

	return ifaceVNUM
}

// IfconfigAddVirtualIface : Add virtual iface by running ifconfig command
func IfconfigAddVirtualIface(ifaceName string, ip string, netmask string) error {
	ifconfigCommand := strings.Replace(ifconfigReplaceString, "IFCONFIG_IFACE_NAME", ifaceName, -1)

	ifconfigCommand = strings.Replace(ifconfigCommand, "IFCONFIG_IFACE_VNUM", strconv.Itoa(getifaceVNUM(ip)), -1)

	ifconfigCommand = strings.Replace(ifconfigCommand, "IFCONFIG_IP", ip, -1)
	ifconfigCommand = strings.Replace(ifconfigCommand, "IFCONFIG_NETMASK", netmask, -1)

	err := runCommand(ifconfigCommand)
	if err != nil {
		return err
	}

	return nil
}

// IfconfigDeleteVirtualIface : Delete and unload ifconfig scripts for external network
func IfconfigDeleteVirtualIface(ifaceName string, ip string) error {
	ifconfigCommand := strings.Replace(ifconfigDownString, "IFCONFIG_IFACE_NAME", ifaceName, -1)
	ifconfigCommand = strings.Replace(ifconfigCommand, "IFCONFIG_IFACE_VNUM", strconv.Itoa(getifaceVNUM(ip)), -1)

	err := runCommand(ifconfigCommand)
	if err != nil {
		return err
	}

	return nil
}
