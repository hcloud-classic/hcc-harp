package adaptiveip

import (
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/pf"
	"hcc/harp/pb"
	"strings"
)

func checkWriteAdaptiveIPNetworkConfigAllArgs(args map[string]interface{}) bool {
	_, extIPAddressOk := args["ext_iface_ip_address"].(string)
	_, netmaskOk := args["netmask"].(string)
	_, gatewayOk := args["gateway"].(string)
	_, startIPAdressOk := args["start_ip_address"].(string)
	_, endIPAdressOk := args["end_ip_address"].(string)

	return extIPAddressOk && netmaskOk && gatewayOk && startIPAdressOk && endIPAdressOk
}

func writeAdaptiveIPNetworkConfig(args map[string]interface{}) (interface{}, error) {
	if !checkWriteAdaptiveIPNetworkConfigAllArgs(args) {
		return nil, errors.New("needed arguments: ext_iface_ip_address, netmask, gateway, start_ip_address," +
			"end_ip_address")
	}

	extIPAddress, _ := args["ext_iface_ip_address"].(string)
	netmask, _ := args["netmask"].(string)
	gateway, _ := args["gateway"].(string)
	startIP, _ := args["start_ip_address"].(string)
	endIP, _ := args["end_ip_address"].(string)

	var adaptiveIP pb.AdaptiveIPSetting
	adaptiveIP.ExtIfaceIpAddress = extIPAddress
	adaptiveIP.Netmask = netmask
	adaptiveIP.Gateway = gateway
	adaptiveIP.StartIpAddress = startIP
	adaptiveIP.EndIpAddress = endIP

	err := configext.CheckAdaptiveIPConfig(&adaptiveIP)
	if err != nil {
		return nil, err
	}

	var networkConfigData string

	networkConfigData = networkConfigBase
	networkConfigData = strings.Replace(networkConfigData, extIfaceAddrReplaceString, extIPAddress, -1)
	networkConfigData = strings.Replace(networkConfigData, netmaskReplaceString, netmask, -1)
	networkConfigData = strings.Replace(networkConfigData, gatewayAddrReplaceString, gateway, -1)
	networkConfigData = strings.Replace(networkConfigData, startIPReplaceString, startIP, -1)
	networkConfigData = strings.Replace(networkConfigData, endIPReplaceString, endIP, -1)

	err = fileutil.WriteFile(config.AdaptiveIP.NetworkConfigFile, networkConfigData)
	if err != nil {
		return nil, err
	}

	return adaptiveIP, nil
}

// WriteNetworkConfigAndReloadHarpNetwork : Write network config files then reload network related services.
func WriteNetworkConfigAndReloadHarpNetwork(args map[string]interface{}) (interface{}, error) {
	adaptiveIP, err := writeAdaptiveIPNetworkConfig(args)
	if err != nil {
		return nil, err
	}

	err = pf.PreparePFConfigFiles()
	if err != nil {
		return nil, err
	}

	err = LoadHarpPFRules()
	if err != nil {
		return nil, err
	}

	return adaptiveIP, nil
}
