package config

type adaptiveIPNetwork struct {
	ExtIfaceIPAddr string `goconf:"adaptiveip_network:adaptiveip_ext_iface_ip_addr"` // ExtIfaceIPAddr : IP address of external interface for use adaptive IP
	Netmask        string `goconf:"adaptiveip_network:adaptiveip_netmask"`           // Netmask : netmask for use adaptive IP
	GatewayAddr    string `goconf:"adaptiveip_network:adaptiveip_gateway_addr"`      // GatewayAddr : gateway address for use adaptive IP
	StartIPAddr    string `goconf:"adaptiveip_network:adaptiveip_start_ip"`          // StartIPAddr : Start IP address for use adaptive IP
	EndIPAddr      string `goconf:"adaptiveip_network:adaptiveip_end_ip"`            // EndIPAddr : End IP address for use adaptive IP
}

// AdaptiveIPNetwork : adaptiveIP network networkConfig structure
var AdaptiveIPNetwork adaptiveIPNetwork
