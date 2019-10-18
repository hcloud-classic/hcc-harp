package config

type dhcpd struct {
	LocalDHCPDServiceName   string `goconf:"dhcpd:dhcpd_local_dhcpd_service_name"`   // LocalDHCPDServiceName : Local server's dhcpd service name
	LocalConfigFileLocation string `goconf:"dhcpd:dhcpd_local_config_file_location"` // LocalConfigFileLocation : Local server's dhcpd file location
	ConfigFileLocation      string `goconf:"dhcpd:dhcpd_config_file_location"`       // ConfigFileLocation : Config file location need to include in dhcpd
	MinLeaseTime            int64  `goconf:"dhcpd:dhcpd_min_lease_time"`             // MinLeaseTime : Minimum lease time for dhcpd
	DefaultLeaseTime        int64  `goconf:"dhcpd:dhcpd_default_lease_time"`         // DefaultLeaseTime : Default lease time for dhcpd
	MaxLeaseTime            int64  `goconf:"dhcpd:dhcpd_max_lease_time"`             // MaxLeaseTime : Max lease time for dhcpd
}

// DHCPD : dhcpd config structure
var DHCPD dhcpd
