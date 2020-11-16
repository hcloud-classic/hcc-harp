package config

import (
	"errors"
	"hcc/harp/lib/iputil"
	"hcc/harp/model"
	"net"
)

type adaptiveIP struct {
	ExternalIfaceName          string `goconf:"adaptiveip_external_iface_name"`                      // ExternalIfaceName : External interface name
	InternalIfaceName          string `goconf:"adaptiveip_internal_iface_name"`                      // InternalIfaceName : Internal interface name
	PFBaseConfigFileLocation   string `goconf:"adaptiveip:adaptiveip_pf_base_config_file_location"`  // PFBaseConfigFileLocation : Base configuration file location to make harp module's pf.rules
	PFRulesFileLocation        string `goconf:"adaptiveip:adaptiveip_pf_rules_file_location"`        // PFRulesFileLocation : pf.rules file location to use in harp module
	PFBinatConfigFileLocation  string `goconf:"adaptiveip:adaptiveip_pf_binat_config_file_location"` // PFBinatConfigFileLocation : PF configuration file location of binat
	PFnatConfigFileLocation    string `goconf:"adaptiveip:adaptiveip_pf_nat_config_file_location"`   // PFnatConfigFileLocation : PF configuration file location of nat
	IfconfigScriptFileLocation string `goconf:"adaptiveip:adaptiveip_ifconfig_script_file_location"` // IfconfigScriptFileLocation : Script file location of ifconfig
	NetworkConfigFile          string `goconf:"adaptiveip:adaptiveip_network_config_file"`           // NetworkConfigFile : Adaptive IP network networkConfig file location
	DefaultExtIfaceIPAddr      string `goconf:"adaptiveip:adaptiveip_default_ext_iface_ip_addr"`     // DefaultExtIfaceIPAddr : Default IP address of external interface for use adaptive IP
	DefaultNetmask             string `goconf:"adaptiveip:adaptiveip_default_netmask"`               // DefaultNetmask : Default netmask for use adaptive IP
	DefaultGatewayAddr         string `goconf:"adaptiveip:adaptiveip_default_gateway_addr"`          // DefaultGatewayAddr : Default gateway address for use adaptive IP
	DefaultStartIPAddr         string `goconf:"adaptiveip:adaptiveip_default_start_ip"`              // DefaultStartIPAddr : Default start IP address for use adaptive IP
	DefaultEndIPAddr           string `goconf:"adaptiveip:adaptiveip_default_end_ip"`                // DefaultEndIPAddr : Default end IP address for use adaptive IP
	ArpingRetryCount           int64  `goconf:"adaptiveip:adaptiveip_arping_retry_count"`            // ArpingRetryCount : Retry count for arping to check duplicated IP address
	ArpingRoutineMaxNum        int64  `goconf:"adaptiveip:adaptiveip_arping_routine_max_num"`        // ArpingRoutineMaxNum : Max number of arping go routine jobs
}

// AdaptiveIP : adaptiveIP config structure
var AdaptiveIP adaptiveIP

// CheckAdaptiveIPConfig : Check configuration of Adaptive IP
func CheckAdaptiveIPConfig(adaptiveIP model.AdaptiveIP) error {
	netNetwork, err := iputil.CheckNetwork(adaptiveIP.ExtIfaceIPAddress,
		adaptiveIP.Netmask)
	if err != nil {
		return err
	}

	err = iputil.CheckGateway(*netNetwork, adaptiveIP.GatewayAddress)
	if err != nil {
		return err
	}

	netStartIP := iputil.CheckValidIP(adaptiveIP.StartIPAddress)
	if netStartIP == nil {
		return errors.New("wrong public start IP address")
	}

	netEndIP := iputil.CheckValidIP(adaptiveIP.EndIPAddress)
	if netEndIP == nil {
		return errors.New("wrong public end IP address")
	}

	isStartIPContainedInNetwork := netNetwork.Contains(netStartIP)
	if isStartIPContainedInNetwork == false {
		return errors.New("start IP address is not in the public network address")
	}

	isEndIPContainedInNetwork := netNetwork.Contains(netEndIP)
	if isEndIPContainedInNetwork == false {
		return errors.New("end IP address is not in the public network address")
	}

	totalAvailableIPs, err := iputil.GetTotalAvailableIPs(netNetwork.IP.String(), net.IP(netNetwork.Mask).String())
	if err != nil {
		return err
	}

	ipRangeCount, err := iputil.GetIPRangeCount(netStartIP, netEndIP)
	if err != nil {
		return err
	}

	if ipRangeCount > totalAvailableIPs {
		return errors.New("IP range count is bigger than total available IPs")
	}

	return nil
}
