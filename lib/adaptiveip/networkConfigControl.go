package adaptiveip

import (
	"errors"
	pb "hcc/harp/action/grpc/rpcharp"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/pf"
	"strings"
)

func checkWriteAdaptiveIPNetworkConfigAllArgs(adaptiveIPSetting *pb.AdaptiveIPSetting) bool {
	extIPAddressOk := len(adaptiveIPSetting.ExtIfaceIPAddress) != 0
	netmaskOk := len(adaptiveIPSetting.Netmask) != 0
	gatewayOk := len(adaptiveIPSetting.GatewayAddress) != 0
	startIPAddressOk := len(adaptiveIPSetting.StartIPAddress) != 0
	endIPAddressOk := len(adaptiveIPSetting.EndIPAddress) != 0

	return extIPAddressOk && netmaskOk && gatewayOk && startIPAddressOk && endIPAddressOk
}

func writeAdaptiveIPNetworkConfig(in *pb.ReqCreateAdaptiveIPSetting) (*pb.AdaptiveIPSetting, error) {
	adaptiveIPSetting := in.GetAdaptiveIPSetting()
	if adaptiveIPSetting == nil {
		return nil, errors.New("AdaptiveIPSetting is nil")
	}

	if !checkWriteAdaptiveIPNetworkConfigAllArgs(adaptiveIPSetting) {
		return nil, errors.New("needed arguments: ext_iface_ip_address, netmask, gateway, start_ip_address," +
			"end_ip_address")
	}

	extIPAddress := adaptiveIPSetting.ExtIfaceIPAddress
	netmask := adaptiveIPSetting.Netmask
	gateway := adaptiveIPSetting.GatewayAddress
	startIP := adaptiveIPSetting.StartIPAddress
	endIP := adaptiveIPSetting.EndIPAddress

	var adaptiveIP pb.AdaptiveIPSetting
	adaptiveIP.ExtIfaceIPAddress = extIPAddress
	adaptiveIP.Netmask = netmask
	adaptiveIP.GatewayAddress = gateway
	adaptiveIP.StartIPAddress = startIP
	adaptiveIP.EndIPAddress = endIP

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

	return &adaptiveIP, nil
}

// WriteNetworkConfigAndReloadHarpNetwork : Write network config files then reload network related services.
func WriteNetworkConfigAndReloadHarpNetwork(in *pb.ReqCreateAdaptiveIPSetting) (*pb.AdaptiveIPSetting, error) {
	adaptiveIP, err := writeAdaptiveIPNetworkConfig(in)
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
