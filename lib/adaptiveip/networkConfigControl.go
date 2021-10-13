package adaptiveip

import (
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configadaptiveip"
	"hcc/harp/lib/fileutil"
	"innogrid.com/hcloud-classic/pb"
	"strings"
	"sync"
)

func checkWriteAdaptiveIPNetworkConfigAllArgs(adaptiveIPSetting *pb.AdaptiveIPSetting) bool {
	extIPAddressOk := len(adaptiveIPSetting.ExtIfaceIPAddress) != 0
	netmaskOk := len(adaptiveIPSetting.Netmask) != 0
	gatewayOk := len(adaptiveIPSetting.GatewayAddress) != 0
	internalStartIPAddressOk := len(adaptiveIPSetting.InternalStartIPAddress) != 0
	internalEndIPAddressOk := len(adaptiveIPSetting.InternalEndIPAddress) != 0
	externalStartIPAddressOk := len(adaptiveIPSetting.ExternalStartIPAddress) != 0
	externalEndIPAddressOk := len(adaptiveIPSetting.ExternalEndIPAddress) != 0

	return extIPAddressOk && netmaskOk && gatewayOk &&
		internalStartIPAddressOk && internalEndIPAddressOk &&
		externalStartIPAddressOk && externalEndIPAddressOk
}

func writeAdaptiveIPNetworkConfig(in *pb.ReqCreateAdaptiveIPSetting) (*pb.AdaptiveIPSetting, error) {
	adaptiveIPSetting := in.GetAdaptiveipSetting()
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
	internalStartIP := adaptiveIPSetting.InternalStartIPAddress
	internalEndIP := adaptiveIPSetting.InternalEndIPAddress
	externalStartIP := adaptiveIPSetting.ExternalStartIPAddress
	externalEndIP := adaptiveIPSetting.ExternalEndIPAddress

	var adaptiveIP pb.AdaptiveIPSetting
	adaptiveIP.ExtIfaceIPAddress = extIPAddress
	adaptiveIP.Netmask = netmask
	adaptiveIP.GatewayAddress = gateway
	adaptiveIP.InternalStartIPAddress = internalStartIP
	adaptiveIP.InternalEndIPAddress = internalEndIP
	adaptiveIP.ExternalStartIPAddress = externalStartIP
	adaptiveIP.ExternalEndIPAddress = externalEndIP

	err := configadaptiveip.CheckAdaptiveIPConfig(&adaptiveIP)
	if err != nil {
		return nil, err
	}

	var networkConfigData string

	networkConfigData = networkConfigBase
	networkConfigData = strings.Replace(networkConfigData, extIfaceAddrReplaceString, extIPAddress, -1)
	networkConfigData = strings.Replace(networkConfigData, netmaskReplaceString, netmask, -1)
	networkConfigData = strings.Replace(networkConfigData, gatewayAddrReplaceString, gateway, -1)
	networkConfigData = strings.Replace(networkConfigData, internalStartIPReplaceString, internalStartIP, -1)
	networkConfigData = strings.Replace(networkConfigData, internalEndIPReplaceString, internalEndIP, -1)
	networkConfigData = strings.Replace(networkConfigData, externalStartIPReplaceString, externalStartIP, -1)
	networkConfigData = strings.Replace(networkConfigData, externalEndIPReplaceString, externalEndIP, -1)

	err = fileutil.WriteFile(config.AdaptiveIP.NetworkConfigFile, networkConfigData)
	if err != nil {
		return nil, err
	}

	return &adaptiveIP, nil
}

var networkReloadLock sync.Mutex

// WriteNetworkConfigAndReloadHarpNetwork : Write network config files then reload network related services.
func WriteNetworkConfigAndReloadHarpNetwork(in *pb.ReqCreateAdaptiveIPSetting) (*pb.AdaptiveIPSetting, error) {
	networkReloadLock.Lock()

	adaptiveIP, err := writeAdaptiveIPNetworkConfig(in)
	if err != nil {
		goto ERROR
	}

	err = LoadFirewall()
	if err != nil {
		goto ERROR
	}

	networkReloadLock.Unlock()
	return adaptiveIP, nil
ERROR:
	networkReloadLock.Unlock()
	return nil, err
}
