package config

type adaptiveIP struct {
	ExternalIfaceName          string `goconf:"adaptiveip_external_iface_name"`                  // ExternalIfaceName : External interface name
	InternalIfaceName          string `goconf:"adaptiveip_internal_iface_name"`                  // InternalIfaceName : Internal interface name
	NetworkConfigFile          string `goconf:"adaptiveip:adaptiveip_network_config_file"`       // NetworkConfigFile : Adaptive IP network networkConfig file location
	DefaultExtIfaceIPAddr      string `goconf:"adaptiveip:adaptiveip_default_ext_iface_ip_addr"` // DefaultExtIfaceIPAddr : Default IP address of external interface for use AdaptiveIP
	DefaultNetmask             string `goconf:"adaptiveip:adaptiveip_default_netmask"`           // DefaultNetmask : Default netmask for use adaptive IP
	DefaultGatewayAddr         string `goconf:"adaptiveip:adaptiveip_default_gateway_addr"`      // DefaultGatewayAddr : Default gateway address for use adaptive IP
	DefaultInternalStartIPAddr string `goconf:"adaptiveip:adaptiveip_default_internal_start_ip"` // DefaultInternalStartIPAddr : Default internal start IP address for use AdaptiveIP
	DefaultInternalEndIPAddr   string `goconf:"adaptiveip:adaptiveip_default_internal_end_ip"`   // DefaultInternalEndIPAddr : Default internal end IP address for use AdaptiveIP
	DefaultExternalStartIPAddr string `goconf:"adaptiveip:adaptiveip_default_external_start_ip"` // DefaultExternalStartIPAddr : Default external start IP address for use AdaptiveIP
	DefaultExternalEndIPAddr   string `goconf:"adaptiveip:adaptiveip_default_external_end_ip"`   // DefaultExternalEndIPAddr : Default external end IP address for use AdaptiveIP
	ArpingRoutineMaxNum        int64  `goconf:"adaptiveip:adaptiveip_arping_routine_max_num"`    // ArpingRoutineMaxNum : Max number of arping go routine jobs
}

// AdaptiveIP : adaptiveIP config structure
var AdaptiveIP adaptiveIP
