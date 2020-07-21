package configext

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"hcc/harp/model"
)

// GetAdaptiveIPNetwork : Get AdaptiveIP's configured information
func GetAdaptiveIPNetwork() model.AdaptiveIP {
	var adaptiveIP model.AdaptiveIP

	err := config.AdaptiveIPNetworkConfigParser()
	if err != nil {
		logger.Logger.Println("AdaptiveIP network networkConfig not found! Using default network information.")

		adaptiveIP.ExtIfaceIPAddress = config.AdaptiveIP.DefaultExtIfaceIPAddr
		adaptiveIP.Netmask = config.AdaptiveIP.DefaultNetmask
		adaptiveIP.GatewayAddress = config.AdaptiveIP.DefaultGatewayAddr
		adaptiveIP.StartIPAddress = config.AdaptiveIP.DefaultStartIPAddr
		adaptiveIP.EndIPAddress = config.AdaptiveIP.DefaultEndIPAddr
	} else {
		adaptiveIP.ExtIfaceIPAddress = config.AdaptiveIPNetwork.ExtIfaceIPAddr
		adaptiveIP.Netmask = config.AdaptiveIPNetwork.Netmask
		adaptiveIP.GatewayAddress = config.AdaptiveIPNetwork.GatewayAddr
		adaptiveIP.StartIPAddress = config.AdaptiveIPNetwork.StartIPAddr
		adaptiveIP.EndIPAddress = config.AdaptiveIPNetwork.EndIPAddr
	}

	return adaptiveIP
}
