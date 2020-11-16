package adaptiveip

var extIfaceAddrReplaceString = "ADAPTIVEIP_NETWORK_EXT_IFACE_ADDR"
var netmaskReplaceString = "ADAPTIVEIP_NETWORK_NETMASK"
var gatewayAddrReplaceString = "ADAPTIVEIP_NETWORK_GATEWAY_ADDR"
var startIPReplaceString = "ADAPTIVEIP_NETWORK_START_IP"
var endIPReplaceString = "ADAPTIVEIP_NETWORK_END_IP"

var networkConfigBase = "[adaptiveip_network]\n" +
	"adaptiveip_ext_iface_ip_addr " + extIfaceAddrReplaceString + "\n" +
	"adaptiveip_netmask " + netmaskReplaceString + "\n" +
	"adaptiveip_gateway_addr " + gatewayAddrReplaceString + "\n" +
	"adaptiveip_start_ip " + startIPReplaceString + "\n" +
	"adaptiveip_end_ip " + endIPReplaceString + "\n"
