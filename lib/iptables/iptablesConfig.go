package iptables

import "hcc/harp/lib/configext"

// InitIPTABLES : Prepare for use iptables
func InitIPTABLES() error {
	adaptiveIP := configext.GetAdaptiveIPNetwork()

	err := configext.CheckAdaptiveIPConfig(adaptiveIP)
	if err != nil {
		return err
	}

	//err = ReplaceBaseConfigAnchorStrings()
	//if err != nil {
	//	return err
	//}

	return nil
}
