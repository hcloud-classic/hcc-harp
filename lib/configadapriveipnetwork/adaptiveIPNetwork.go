package configadapriveipnetwork

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"innogrid.com/hcloud-classic/pb"
)

// GetAdaptiveIPNetwork : Get AdaptiveIP's configured information
func GetAdaptiveIPNetwork() *pb.AdaptiveIPSetting {
	var adaptiveIP pb.AdaptiveIPSetting

	err := config.AdaptiveIPNetworkConfigParser()
	if err != nil {
		logger.Logger.Println("AdaptiveIP network networkConfig not found! Using default network information.")

		adaptiveIP.ExtIfaceIPAddress = config.AdaptiveIP.DefaultExtIfaceIPAddr
		adaptiveIP.Netmask = config.AdaptiveIP.DefaultNetmask
		adaptiveIP.GatewayAddress = config.AdaptiveIP.DefaultGatewayAddr
		adaptiveIP.InternalStartIPAddress = config.AdaptiveIP.DefaultInternalStartIPAddr
		adaptiveIP.InternalEndIPAddress = config.AdaptiveIP.DefaultInternalEndIPAddr
		adaptiveIP.ExternalStartIPAddress = config.AdaptiveIP.DefaultExternalStartIPAddr
		adaptiveIP.ExternalEndIPAddress = config.AdaptiveIP.DefaultExternalEndIPAddr
	} else {
		adaptiveIP.ExtIfaceIPAddress = config.AdaptiveIPNetwork.ExtIfaceIPAddr
		adaptiveIP.Netmask = config.AdaptiveIPNetwork.Netmask
		adaptiveIP.GatewayAddress = config.AdaptiveIPNetwork.GatewayAddr
		adaptiveIP.InternalStartIPAddress = config.AdaptiveIPNetwork.InternalStartIPAddr
		adaptiveIP.InternalEndIPAddress = config.AdaptiveIPNetwork.InternalEndIPAddr
		adaptiveIP.ExternalStartIPAddress = config.AdaptiveIPNetwork.ExternalStartIPAddr
		adaptiveIP.ExternalEndIPAddress = config.AdaptiveIPNetwork.ExternalEndIPAddr
	}

	return &adaptiveIP
}
