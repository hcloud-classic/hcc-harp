package ifconfig

var ifconfigShell = "#!/bin/bash/\n"
var ifconfigFilenamePrefix = "ifconfig_"
var ifconfigReplaceString = "ifconfig IFCONFIG_IFACE_NAME:IFCONFIG_IFACE_VNUM IFCONFIG_IP netmask IFCONFIG_NETMASK\n"
var ifconfigDownString = "ifconfig IFCONFIG_IFACE_NAME:IFCONFIG_IFACE_VNUM down\n"
