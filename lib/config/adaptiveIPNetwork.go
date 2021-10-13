package config

type adaptiveIPNetwork struct {
	ExtIfaceIPAddr      string `goconf:"adaptiveip_network:adaptiveip_ext_iface_ip_addr"` // ExtIfaceIPAddr : IP address of external interface for use adaptive IP
	Netmask             string `goconf:"adaptiveip_network:adaptiveip_netmask"`           // Netmask : netmask for use adaptive IP
	GatewayAddr         string `goconf:"adaptiveip_network:adaptiveip_gateway_addr"`      // GatewayAddr : gateway address for use adaptive IP
	InternalStartIPAddr string `goconf:"adaptiveip_network:adaptiveip_internal_start_ip"` // InternalStartIPAddr : Internal start IP address for use with AdaptiveIP
	InternalEndIPAddr   string `goconf:"adaptiveip_network:adaptiveip_internal_end_ip"`   // InternalEndIPAddr : Internal end IP address for use with AdaptiveIP
	ExternalStartIPAddr string `goconf:"adaptiveip_network:adaptiveip_external_start_ip"` // ExternalStartIPAddr : External start IP address for use with AdaptiveIP
	ExternalEndIPAddr   string `goconf:"adaptiveip_network:adaptiveip_external_end_ip"`   // ExternalEndIPAddr : External end IP address for use with AdaptiveIP
}

// AdaptiveIPNetwork : adaptiveIP network networkConfig structure
var AdaptiveIPNetwork adaptiveIPNetwork
