package config

type adaptiveIP struct {
	PFLocalConfigFileLocation string `goconf:"adaptiveip:adaptiveip_pf_local_config_file_location"` // PFLocalConfigFileLocation : PF local configuration file location
	PFConfigFileLocation      string `goconf:"adaptiveip:adaptiveip_pf_config_file_location"`       // PFConfigFileLocation : PF configuration file location
}

// AdaptiveIP : adaptiveIP config structure
var AdaptiveIP adaptiveIP
