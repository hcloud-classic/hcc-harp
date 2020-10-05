package config

type adaptiveIP struct {
	ExternalIfaceName          string `goconf:"adaptiveip_external_iface_name"`                      // ExternalIfaceName : External interface name
	InternalIfaceName          string `goconf:"adaptiveip_internal_iface_name"`                      // InternalIfaceName : Internal interface name
	PFBaseConfigFileLocation   string `goconf:"adaptiveip:adaptiveip_pf_base_config_file_location"`  // PFBaseConfigFileLocation : Base configuration file location to make harp module's pf.rules
	PFRulesFileLocation        string `goconf:"adaptiveip:adaptiveip_pf_rules_file_location"`        // PFRulesFileLocation : pf.rules file location to use in harp module
	PFBinatConfigFileLocation  string `goconf:"adaptiveip:adaptiveip_pf_binat_config_file_location"` // PFBinatConfigFileLocation : PF configuration file location of binat
	PFnatConfigFileLocation    string `goconf:"adaptiveip:adaptiveip_pf_nat_config_file_location"`   // PFnatConfigFileLocation : PF configuration file location of nat
	IPTABLESInitConfigFileLocation string `goconf:"adaptiveip:adaptiveip_iptables_init_config_file_location"`   // IPTABLESInitConfigFileLocation : iptables init configuration file location
	IfconfigScriptFileLocation string `goconf:"adaptiveip:adaptiveip_ifconfig_script_file_location"` // IfconfigScriptFileLocation : Script file location of ifconfig
	NetworkConfigFile          string `goconf:"adaptiveip:adaptiveip_network_config_file"`           // NetworkConfigFile : Adaptive IP network networkConfig file location
	DefaultExtIfaceIPAddr      string `goconf:"adaptiveip:adaptiveip_default_ext_iface_ip_addr"`     // DefaultExtIfaceIPAddr : Default IP address of external interface for use adaptive IP
	DefaultNetmask             string `goconf:"adaptiveip:adaptiveip_default_netmask"`               // DefaultNetmask : Default netmask for use adaptive IP
	DefaultGatewayAddr         string `goconf:"adaptiveip:adaptiveip_default_gateway_addr"`          // DefaultGatewayAddr : Default gateway address for use adaptive IP
	DefaultStartIPAddr         string `goconf:"adaptiveip:adaptiveip_default_start_ip"`              // DefaultStartIPAddr : Default start IP address for use adaptive IP
	DefaultEndIPAddr           string `goconf:"adaptiveip:adaptiveip_default_end_ip"`                // DefaultEndIPAddr : Default end IP address for use adaptive IP
	ArpingRoutineMaxNum        int64  `goconf:"adaptiveip:adaptiveip_arping_routine_max_num"`        // ArpingRoutineMaxNum : Max number of arping go routine jobs
}

// AdaptiveIP : adaptiveIP config structure
var AdaptiveIP adaptiveIP
