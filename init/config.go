package init

import (
	"hcc/harp/lib/adaptiveip"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/syscheck"
)

func configInit() error {
	config.Parser()
	_, err := syscheck.CheckIfaceExist(config.AdaptiveIP.ExternalIfaceName)
	if err != nil {
		return err
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.InternalIfaceName)
	if err != nil {
		return err
	}

	err = dhcpd.CheckLocalDHCPDConfig()
	if err != nil {
		return err
	}

	err = adaptiveip.PreparePFConfigFiles()
	if err != nil {
		return err
	}

	err = adaptiveip.LoadHarpPFRules()
	if err != nil {
		return err
	}

	return nil
}
