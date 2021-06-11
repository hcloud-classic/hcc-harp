package iptables

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/configadapriveipnetwork"
	"hcc/harp/lib/iptablesext"
	"hcc/harp/lib/logger"
)

// ForwardTimpani : Forwarding Timpani connection
func ForwardTimpani() error {
	logger.Logger.Println("Forwarding Timpani connection...")

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()
	err := iptablesext.AdaptiveIPServerForwarding(true, false, adaptiveIP.ExtIfaceIPAddress, config.Timpani.TimpaniAddress)
	if err != nil {
		return err
	}
	err = iptablesext.PortForwarding(true, true, true, false, adaptiveIP.ExtIfaceIPAddress,
		config.Timpani.TimpaniAddress, int(config.Timpani.TimpaniExternalPort), int(config.Timpani.TimpaniInternalPort))
	if err != nil {
		return err
	}

	return err
}
