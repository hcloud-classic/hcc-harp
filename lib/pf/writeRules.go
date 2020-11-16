package pf

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/logger"
)

func writePFRulesFile(pfRulesData string) error {
	err := logger.CreateDirIfNotExist(config.AdaptiveIP.PFRulesFileLocation)
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(config.AdaptiveIP.PFRulesFileLocation, pfRulesData)
	if err != nil {
		return err
	}

	return nil
}