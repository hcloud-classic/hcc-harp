package adaptiveip

import (
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"io/ioutil"
	"strings"
)

func CheckLocalPFConfig() error {
	include := includeStr
	include = strings.Replace(include, "HARP_PF_CONF_LOCATION",
		config.AdaptiveIP.PFConfigFileLocation+"/harp_pf.conf", -1)

	data, err := ioutil.ReadFile(config.AdaptiveIP.PFLocalConfigFileLocation)
	if err != nil {
		return errors.New("failed reading data from local pf config file location")
	}

	isHarpPFIncluded := strings.Contains(string(data), include)
	if !isHarpPFIncluded {
		logger.Logger.Println("Please add this line to pf config file!\n" + include)
		return errors.New("cannot find harp pf config include line from local pf config file")
	}

	return nil
}
