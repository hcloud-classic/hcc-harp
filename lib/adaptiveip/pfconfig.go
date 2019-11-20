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
			commented, err := regexp.MatchString("#[ ]+"+harpBinatanchorRelaceString, scanText)
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

	err = scanner.Err()
	if err != nil {
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
	ipRangeCount, _ := iputil.GetIPRangeCount(netStartIP, netEndIP)

	var binatanchorConfPart = ""
	for i := 0; i < ipRangeCount; i++ {
		binatanchorConfLine := strings.Replace(binatanchorStr, "HARP_SERVER_IP", netStartIP.String(), -1)
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

func writeBinatAnchorConfigFile(binatAnchorFileLocation string, binatAnchorData string) error {
	err := fileutil.WriteFile(binatAnchorFileLocation, binatAnchorData)
	if err != nil {
		return err
	}

	return nil
}

func checkBinatAnchorFileExist(privateIP string) error {
	configFiles, err := getBinatAnchorConfigFiles()
	if err != nil {
		return err
	}
	if len(configFiles) == 1 {
		return nil
	}

	for _, file := range configFiles {
		if file == config.AdaptiveIP.PFBinatConfigFileLocation {
			continue
		}

		binatanchorFileName := file[len(config.AdaptiveIP.PFBinatConfigFileLocation+"/"):]
		if binatanchorFileName == binatanchorFilenamePrefix+privateIP+".conf" {
			return errors.New(privateIP + " is already used in binat anchor rules")
		}
	}

	return nil
}

// CreateAndLoadBinatAnchorConfig : Create binat anchor config file to match private IP address
// to available public IP address. Then load it to pf firewall.
func CreateAndLoadBinatAnchorConfig(privateIP string) error {
	netStartIP := iputil.CheckValidIP(config.AdaptiveIP.PublicStartIP)
	netEndIP := iputil.CheckValidIP(config.AdaptiveIP.PublicEndIP)
	ipRangeCount, _ := iputil.GetIPRangeCount(netStartIP, netEndIP)

	var binatanchorConfData string
	for i := 0; i < ipRangeCount; i++ {
		err := checkBinatAnchorFileExist(netStartIP.String())
		if err != nil {
			logger.Logger.Println(err)
			netStartIP = cidr.Inc(netStartIP)
			continue
		}

		err = checkDuplicatedIPAddress(netStartIP.String())
		if err != nil {
			logger.Logger.Println(err)
			netStartIP = cidr.Inc(netStartIP)
			continue
		}

		binatanchorConfData = binatStr
		binatanchorConfData = strings.Replace(binatanchorConfData, "HARP_EXTERNAL_IFACE_NAME", config.AdaptiveIP.ExternalIfaceName, -1)
		binatanchorConfData = strings.Replace(binatanchorConfData, "HARP_PF_PRIVATE_IP", privateIP, -1)
		binatanchorConfData = strings.Replace(binatanchorConfData, "HARP_PF_PUBLIC_IP", netStartIP.String(), -1)

		binatanchorName := binatanchorFilenamePrefix + netStartIP.String()
		logger.Logger.Println("CreateAndLoadBinatAnchorConfig: Creating config file for " + binatanchorName + " (privateIP: " + privateIP + ")")
		binatanchorConfigFileLocation := config.AdaptiveIP.PFBinatConfigFileLocation + "/" + binatanchorName + ".conf"
		err = writeBinatAnchorConfigFile(binatanchorConfigFileLocation, binatanchorConfData)
		if err != nil {
			return err
		}

		logger.Logger.Println("CreateAndLoadBinatAnchorConfig: Load binat anchor config file for " + binatanchorName)
		err = LoadPFBinatAnchorRule(binatanchorName, binatanchorConfigFileLocation)
		if err != nil {
			return err
		}

		break
	}

	return nil
}
