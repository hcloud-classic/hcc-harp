package pf

var binatanchorFilenamePrefix = "binat_"
var natanchorFilenamePrefix = "nat_"

var harpBinatanchorRelaceString = "HARP_BINAT_ANCHOR_REPLACE_STRING"

var binatStr = "binat on HARP_EXTERNAL_IFACE_NAME from HARP_PF_PRIVATE_IP to any -> HARP_PF_PUBLIC_IP\n"
var binatanchorStr = "binat-anchor " + binatanchorFilenamePrefix + "HARP_SERVER_IP\n"
var natStr = "nat on HARP_INTERNAL_IFACE_NAME from HARP_PF_PRIVATE_IP to any -> HARP_PF_PUBLIC_IP\n"
var natanchorStr = "nat-anchor " + natanchorFilenamePrefix + "HARP_SERVER_IP\n"

var harpNatanchorRelaceString = "HARP_NAT_ANCHOR_REPLACE_STRING"
