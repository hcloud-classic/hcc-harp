package configext

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"hcc/harp/pb"
)

// GetAdaptiveIPNetwork : Get AdaptiveIP's configured information
func GetAdaptiveIPNetwork() *pb.AdaptiveIPSetting {
	var adaptiveIP pb.AdaptiveIPSetting

	err := config.AdaptiveIPNetworkConfigParser()
	if err != nil {
		logger.Logger.Println("AdaptiveIP network networkConfig not found! Using default network information.")

		adaptiveIP.ExtIfaceIpAddress = config.AdaptiveIP.DefaultExtIfaceIPAddr
		adaptiveIP.Netmask = config.AdaptiveIP.DefaultNetmask
		adaptiveIP.Gateway = config.AdaptiveIP.DefaultGatewayAddr
		adaptiveIP.StartIpAddress = config.AdaptiveIP.DefaultStartIPAddr
		adaptiveIP.EndIpAddress = config.AdaptiveIP.DefaultEndIPAddr
	} else {
		adaptiveIP.ExtIfaceIpAddress = config.AdaptiveIPNetwork.ExtIfaceIPAddr
		adaptiveIP.Netmask = config.AdaptiveIPNetwork.Netmask
		adaptiveIP.Gateway = config.AdaptiveIPNetwork.GatewayAddr
		adaptiveIP.StartIpAddress = config.AdaptiveIPNetwork.StartIPAddr
		adaptiveIP.EndIpAddress = config.AdaptiveIPNetwork.EndIPAddr
	}

	return &adaptiveIP
}
