package adaptiveip

import (
	"bufio"
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
)

func checkPFBaseConfig() error {
	file, err := os.Open(config.AdaptiveIP.PFBaseConfigFileLocation)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	isHarpbinatanchorRelaceStringIncluded := false
	var scanText string

	for scanner.Scan() {
		scanText = scanner.Text()
		isHarpbinatanchorRelaceStringIncluded = strings.Contains(scanText, harpBinatanchorRelaceString)
		if isHarpbinatanchorRelaceStringIncluded {
			break
		}
	}

	if isHarpbinatanchorRelaceStringIncluded {
		lineCheckOk := scanText == harpBinatanchorRelaceString
		if !lineCheckOk {
			commented, err := regexp.MatchString("#[ ]+" +harpBinatanchorRelaceString, scanText)
			if err != nil {
				return err
			}
			if commented {
				return errors.New("please comment out HARP_BINAT_ANCHOR_REPLACE_STRING string in base config file")
			}
			return errors.New("please add HARP_BINAT_ANCHOR_REPLACE_STRING string line correctly to base config file")
		}
	} else {
		logger.Logger.Println("Please add this line to pf base config file!\n" + harpBinatanchorRelaceString)
		return errors.New("cannot find binat anchor replace string from pf base config file")
	}

	if err := scanner.Err()
	err != nil {
		return err
	}

	return nil
}

func writePFRulesFile(pfRulesData string) error {
	err := fileutil.WriteFile(config.AdaptiveIP.PFRulesFileLocation, pfRulesData)
	if err != nil {
		return err
	}

	return nil
}

func replaceBaseConfigBinatAnchorString() error {
	err := checkPFBaseConfig()
	if err != nil {
		return err
	}

	baseConfigData, err := ioutil.ReadFile(config.AdaptiveIP.PFBaseConfigFileLocation)
	if err != nil {
		return errors.New("failed reading data from base pf config file location")
	}

	netStartIP := iputil.CheckValidIP(config.AdaptiveIP.PublicStartIP)
	netEndIP := iputil.CheckValidIP(config.AdaptiveIP.PublicEndIP)
	_, ipRangeCount := iputil.GetIPRangeCount(netStartIP, netEndIP)

	var binatanchorConfPart = ""
	for i := 0; i < ipRangeCount; i++ {
		var binatanchorConfLine = binatanchorStr
		binatanchorConfLine = strings.Replace(binatanchorStr, "HARP_SERVER_IP", netStartIP.String(), -1)
		netStartIP = cidr.Inc(netStartIP)

		binatanchorConfPart += binatanchorConfLine
	}

	pfRulesData := strings.Replace(string(baseConfigData), harpBinatanchorRelaceString,
		binatanchorConfPart, -1)

	err = writePFRulesFile(pfRulesData)
	if err != nil {
		return err
	}

	return nil
}

// PreparePFConfigFiles : Prepare pf.rules config file for use in adaptive IP
func PreparePFConfigFiles() error {
	err := checkHarpConfigNetwork()
	if err != nil {
		return err
	}

	err = replaceBaseConfigBinatAnchorString()
	if err != nil {
		return err
	}

	return nil
}
