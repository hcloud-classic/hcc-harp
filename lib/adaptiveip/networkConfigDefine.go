package adaptiveip

var extIfaceAddrReplaceString = "ADAPTIVEIP_NETWORK_EXT_IFACE_ADDR"
var netmaskReplaceString = "ADAPTIVEIP_NETWORK_NETMASK"
var gatewayAddrReplaceString = "ADAPTIVEIP_NETWORK_GATEWAY_ADDR"
var internalStartIPReplaceString = "ADAPTIVEIP_NETWORK_INTERNAL_START_IP"
var internalEndIPReplaceString = "ADAPTIVEIP_NETWORK_INTERNAL_END_IP"
var externalStartIPReplaceString = "ADAPTIVEIP_NETWORK_EXTERNAL_START_IP"
var externalEndIPReplaceString = "ADAPTIVEIP_NETWORK_EXTERNAL_END_IP"

var networkConfigBase = "[adaptiveip_network]\n" +
	"adaptiveip_ext_iface_ip_addr " + extIfaceAddrReplaceString + "\n" +
	"adaptiveip_netmask " + netmaskReplaceString + "\n" +
	"adaptiveip_gateway_addr " + gatewayAddrReplaceString + "\n" +
	"adaptiveip_internal_start_ip " + internalStartIPReplaceString + "\n" +
	"adaptiveip_internal_end_ip " + internalEndIPReplaceString + "\n" +
	"adaptiveip_external_start_ip " + externalStartIPReplaceString + "\n" +
	"adaptiveip_external_end_ip " + externalEndIPReplaceString + "\n"
