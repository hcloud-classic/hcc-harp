package adaptiveip

var harpBinatanchorRelaceString = "HARP_BINAT_ANCHOR_REPLACE_STRING"
var binatanchorFilenamePrefix = "binat_"
var binatanchorStr = "binat-anchor " + binatanchorFilenamePrefix + "HARP_SERVER_IP\n"
var binatStr = "binat on HARP_EXTERNAL_IFACE_NAME from HARP_PF_PRIVATE_IP to any -> HARP_PF_PUBLIC_IP\n"

var harpNatanchorRelaceString = "HARP_NAT_ANCHOR_REPLACE_STRING"
var natanchorFilenamePrefix = "nat_"
var natanchorStr = "nat-anchor " + natanchorFilenamePrefix + "HARP_SERVER_IP\n"
var natStr = "nat on HARP_INTERNAL_IFACE_NAME from HARP_PF_PRIVATE_IP to any -> HARP_PF_PUBLIC_IP\n"

var ifconfigSHELL = "#!/bin/csh/\n"
var ifconfigFilenamePrefix = "ifconfig_"
var ifconfigReplaceString = "ifconfig IFCONFIG_IFACE_NAME IFCONFIG_IP netmask IFCONFIG_NETMASK ALIAS_STATE\n"
