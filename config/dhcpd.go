package config

type dhcpd struct {
	MinLeaseTime       int64  `goconf:"dhcpd:dhcpd_min_lease_time"`       // MinLeaseTime : Minimum lease time for dhcpd
	DefaultLeaseTime   int64  `goconf:"dhcpd:dhcpd_default_lease_time"`   // DefaultLeaseTime : Default lease time for dhcpd
	MaxLeaseTime       int64  `goconf:"dhcpd:dhcpd_max_lease_time"`       // MaxLeaseTime : Max lease time for dhcpd
	ConfigFileLocation string `goconf:"dhcpd:dhcpd_config_file_location"` // ConfigFileLocation : Config file location need to include in dhcpd
}

// DHCPD : dhcpd config structure
var DHCPD dhcpd
