package ifconfig

import "hcc/harp/lib/syscheck"

func ifconfigShell() string {
	if syscheck.OS == "freebsd" {
		return "#!/bin/csh/\n"
	}

	// Linux
	return "#!/bin/bash/\n"
}

var ifconfigFilenamePrefix = "ifconfig_"

func ifconfigReplaceString() string {
	if syscheck.OS == "freebsd" {
		return "ifconfig IFCONFIG_IFACE_NAME IFCONFIG_IP netmask IFCONFIG_NETMASK ALIAS_STATE\n"
	}

	// Linux
	return "ifconfig IFCONFIG_IFACE_NAME:IFCONFIG_IFACE_VNUM IFCONFIG_IP netmask IFCONFIG_NETMASK\n"
}

func ifconfigDownStringLinux() string {
	return "ifconfig IFCONFIG_IFACE_NAME:IFCONFIG_IFACE_VNUM down\n"
}
