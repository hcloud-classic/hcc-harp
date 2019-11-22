package adaptiveip

import (
	"bufio"
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/model"
	"io/ioutil"
	"os"
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
	isHarpnatanchorRelaceStringIncluded := false
	var scanText string

	// binat
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
			return errors.New("please add HARP_BINAT_ANCHOR_REPLACE_STRING string line correctly to base config file")
		}
	} else {
		logger.Logger.Println("Please add this line to pf base config file!\n" + harpBinatanchorRelaceString)
		return errors.New("cannot find binat anchor replace string from pf base config file")
	}

	// NAT
	for scanner.Scan() {
		scanText = scanner.Text()
		isHarpnatanchorRelaceStringIncluded = strings.Contains(scanText, harpNatanchorRelaceString)
		if isHarpnatanchorRelaceStringIncluded {
			break
		}
	}

	if isHarpnatanchorRelaceStringIncluded {
		lineCheckOk := scanText == harpNatanchorRelaceString
		if !lineCheckOk {
			return errors.New("please add HARP_NAT_ANCHOR_REPLACE_STRING string line correctly to base config file")
		}
	} else {
		logger.Logger.Println("Please add this line to pf base config file!\n" + harpNatanchorRelaceString)
		return errors.New("cannot find nat anchor replace string from pf base config file")
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

func replaceBaseConfigAnchorStrings() error {
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

	// binat
	var binatanchorConfPart = ""
	netStartIPtemp := netStartIP
	for i := 0; i < ipRangeCount; i++ {
		binatanchorConfLine := strings.Replace(binatanchorStr, "HARP_SERVER_IP", netStartIPtemp.String(), -1)
		netStartIPtemp = cidr.Inc(netStartIPtemp)

		binatanchorConfPart += binatanchorConfLine
	}

	pfRulesData := strings.Replace(string(baseConfigData), harpBinatanchorRelaceString,
		binatanchorConfPart, -1)

	// NAT
	var natanchorConfPart = ""
	netStartIPtemp = netStartIP
	for i := 0; i < ipRangeCount; i++ {
		natanchorConfLine := strings.Replace(natanchorStr, "HARP_SERVER_IP", netStartIPtemp.String(), -1)
		netStartIPtemp = cidr.Inc(netStartIPtemp)

		natanchorConfPart += natanchorConfLine
	}

	pfRulesData = strings.Replace(pfRulesData, harpNatanchorRelaceString,
		natanchorConfPart, -1)

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

	err = replaceBaseConfigAnchorStrings()
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

func createAndLoadBinatAnchorConfig(privateIP string, publicIP string) error {
	var binatanchorConfData string

	binatanchorConfData = binatStr
	binatanchorConfData = strings.Replace(binatanchorConfData, "HARP_EXTERNAL_IFACE_NAME", config.AdaptiveIP.ExternalIfaceName, -1)
	binatanchorConfData = strings.Replace(binatanchorConfData, "HARP_PF_PRIVATE_IP", privateIP, -1)
	binatanchorConfData = strings.Replace(binatanchorConfData, "HARP_PF_PUBLIC_IP", publicIP, -1)

	binatanchorName := binatanchorFilenamePrefix + publicIP
	logger.Logger.Println("createAndLoadBinatAnchorConfig: Creating config file for " + binatanchorName +
		" (publicIP: " + publicIP + ", privateIP: " + privateIP + ")")
	binatanchorConfigFileLocation := config.AdaptiveIP.PFBinatConfigFileLocation + "/" + binatanchorName + ".conf"
	err := fileutil.WriteFile(binatanchorConfigFileLocation, binatanchorConfData)
	if err != nil {
		return err
	}

	logger.Logger.Println("createAndLoadBinatAnchorConfig: Load binat anchor config file for " + binatanchorName)
	err = LoadPFAnchorRule(binatanchorName, binatanchorConfigFileLocation)
	if err != nil {
		return err
	}

	return nil
}

func createAndLoadnatAnchorConfig(privateIP string, publicIP string) error {
	var natanchorConfData string

	natanchorConfData = natStr
	natanchorConfData = strings.Replace(natanchorConfData, "HARP_INTERNAL_IFACE_NAME", config.AdaptiveIP.InternalIfaceName, -1)
	natanchorConfData = strings.Replace(natanchorConfData, "HARP_PF_PRIVATE_IP", privateIP, -1)
	natanchorConfData = strings.Replace(natanchorConfData, "HARP_PF_PUBLIC_IP", publicIP, -1)

	natanchorName := natanchorFilenamePrefix + publicIP
	logger.Logger.Println("createAndLoadnatAnchorConfig: Creating config file for " + natanchorName +
		" (publicIP: " + publicIP + ", privateIP: " + privateIP + ")")
	natanchorConfigFileLocation := config.AdaptiveIP.PFnatConfigFileLocation + "/" + natanchorName + ".conf"
	err := fileutil.WriteFile(natanchorConfigFileLocation, natanchorConfData)
	if err != nil {
		return err
	}

	logger.Logger.Println("createAndLoadnatAnchorConfig: Load nat anchor config file for " + natanchorName)
	err = LoadPFAnchorRule(natanchorName, natanchorConfigFileLocation)
	if err != nil {
		return err
	}

	return nil
}

// CreateAndLoadAnchorConfig : Create anchor config files to match private IP address
// to available public IP address. Then load them to pf firewall.
func CreateAndLoadAnchorConfig(publicIP string, privateIP string, subnet model.Subnet) error {
	ipMap := getAvailableIPsStatusMap()

	err := checkBinatAnchorFileExist(publicIP)
	if err != nil {
		goto Error
	}

	if !ipMap[publicIP] {
		err = errors.New("CreateAndLoadAnchorConfig: " + publicIP + " is a duplicated IP address")
		goto Error
	}

	err = createAndLoadBinatAnchorConfig(privateIP, publicIP)
	if err != nil {
		goto Error
	}

	err = createAndLoadnatAnchorConfig(privateIP, publicIP)
	if err != nil {
		goto Error
	}

	err = createAndLoadIfconfigScript(config.AdaptiveIP.InternalIfaceName, config.AdaptiveIP.ExternalIfaceName,
		subnet.Gateway, publicIP, subnet.Netmask, config.AdaptiveIP.PublicNetworkNetmask)
	if err != nil {
		goto Error
	}

	return nil
Error:
	return err
}
