package config

type adaptiveIP struct {
	ExternalIfaceName         string `goconf:"adaptiveip_external_iface_name"`                      // ExternalIfaceName : External interface name
	PFBaseConfigFileLocation  string `goconf:"adaptiveip:adaptiveip_pf_base_config_file_location"`  // PFBaseConfigFileLocation : Base configuration file location to make harp module's pf.rules
	PFRulesFileLocation       string `goconf:"adaptiveip:adaptiveip_pf_rules_file_location"`        // PFRulesFileLocation : pf.rules file location to use in harp module
	PFBinatConfigFileLocation string `goconf:"adaptiveip:adaptiveip_pf_binat_config_file_location"` // PFBinatConfigFileLocation : PF configuration file location of binat
	PublicNetworkAddress      string `goconf:"adaptiveip:adaptiveip_public_network_address"`        // PublicNetworkAddress : Public network address for allocate adaptive IP address
	PublicNetworkNetmask      string `goconf:"adaptiveip:adaptiveip_public_network_netmask"`        // PublicNetworkNetmask : Netmask of public network address
	PublicStartIP             string `goconf:"adaptiveip:adaptiveip_public_start_ip"`               // PublicStartIP : Public start IP address for using adaptive IP
	PublicEndIP               string `goconf:"adaptiveip:adaptiveip_public_end_ip"`                 // PublicEndIP : Public end IP address for using adaptive IP
	ArpingRetryCount          int64  `goconf:"adaptiveip:adaptiveip_arping_retry_count"`            // ArpingRetryCount : Retry count for arping to check duplicated IP address
	ArpingRoutineMaxNum       int64  `goconf:"adaptiveip:adaptiveip_arping_routine_max_num"`        // ArpingRoutineMaxNum : Max number of arping go routine jobs
}

// AdaptiveIP : adaptiveIP config structure
var AdaptiveIP adaptiveIP
