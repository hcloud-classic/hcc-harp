package pf

import (
	"bufio"
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"hcc/harp/lib/config"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/ifconfig"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

// ReplaceBaseConfigAnchorStrings : Create binat anchor config file
func ReplaceBaseConfigAnchorStrings() error {
	err := checkPFBaseConfig()
	if err != nil {
		return err
	}

	baseConfigData, err := ioutil.ReadFile(config.AdaptiveIP.PFBaseConfigFileLocation)
	if err != nil {
		return errors.New("failed reading data from base pf config file location")
	}

	adaptiveIP := config.GetAdaptiveIPNetwork()
	netStartIP := iputil.CheckValidIP(adaptiveIP.StartIPAddress)
	netEndIP := iputil.CheckValidIP(adaptiveIP.EndIPAddress)
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

// GetBinatAnchorConfigFiles : Get locations of binat anchor config files
func GetBinatAnchorConfigFiles() ([]string, error) {
	var files []string

	folder := config.AdaptiveIP.PFBinatConfigFileLocation
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// CheckBinatAnchorFileExist : Check if binat anchor config file is exists
func CheckBinatAnchorFileExist(publicIP string) error {
	configFiles, err := GetBinatAnchorConfigFiles()
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
		if binatanchorFileName == binatanchorFilenamePrefix+publicIP+".conf" {
			return errors.New(publicIP + " is already used in binat anchor rules")
		}
	}

	return nil
}

func getnatAnchorConfigFiles() ([]string, error) {
	var files []string

	folder := config.AdaptiveIP.PFnatConfigFileLocation
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func loadExstingBinatRules() error {
	logger.Logger.Println("Loading existing binat rules...")

	configFiles, err := GetBinatAnchorConfigFiles()
	if err != nil {
		return err
	}
	if len(configFiles) == 1 {
		return nil
	}

	var binatanchorFileName string
	var binatanchorName string

	for i := 0; i < len(configFiles); i++ {
		if configFiles[i] == config.AdaptiveIP.PFBinatConfigFileLocation {
			continue
		}

		binatanchorFileName = configFiles[i][len(config.AdaptiveIP.PFBinatConfigFileLocation+"/"):]
		if !strings.Contains(binatanchorFileName, binatanchorFilenamePrefix) ||
			!strings.Contains(binatanchorFileName, ".conf") {
			logger.Logger.Println("Wrong binat anchor filename: " + binatanchorFileName)
			logger.Logger.Println("Filename must be as '" + binatanchorFilenamePrefix + "XXX.conf'")
			continue
		}

		binatanchorName = binatanchorFileName[0 : len(binatanchorFileName)-len(".conf")]
		err = LoadPFAnchorRule(binatanchorName, configFiles[i])
		if err != nil {
			logger.Logger.Println(err)
		}
	}

	return nil
}

func loadExstingnatRules() error {
	logger.Logger.Println("Loading existing NAT rules...")

	configFiles, err := getnatAnchorConfigFiles()
	if err != nil {
		return err
	}
	if len(configFiles) == 1 {
		return nil
	}

	var natanchorFileName string
	var natanchorName string

	for i := 0; i < len(configFiles); i++ {
		if configFiles[i] == config.AdaptiveIP.PFnatConfigFileLocation {
			continue
		}

		natanchorFileName = configFiles[i][len(config.AdaptiveIP.PFnatConfigFileLocation+"/"):]
		if !strings.Contains(natanchorFileName, natanchorFilenamePrefix) ||
			!strings.Contains(natanchorFileName, ".conf") {
			logger.Logger.Println("Wrong nat anchor filename: " + natanchorFileName)
			logger.Logger.Println("Filename must be as '" + natanchorFilenamePrefix + "XXX.conf'")
			continue
		}

		natanchorName = natanchorFileName[0 : len(natanchorFileName)-len(".conf")]
		err = LoadPFAnchorRule(natanchorName, configFiles[i])
		if err != nil {
			logger.Logger.Println(err)
		}
	}

	return nil
}

// LoadExstingBinatAndNATRules : Load binat and nat rules configured by harp module
func LoadExstingBinatAndNATRules() error {
	err := loadExstingBinatRules()
	if err != nil {
		return err
	}
	err = loadExstingnatRules()
	if err != nil {
		return err
	}

	return nil
}

// PreparePFConfigFiles : Prepare pf.rules config file for use in adaptive IP
func PreparePFConfigFiles() error {
	adaptiveIP := config.GetAdaptiveIPNetwork()

	err := config.CheckAdaptiveIPConfig(adaptiveIP)
	if err != nil {
		return err
	}

	err = ReplaceBaseConfigAnchorStrings()
	if err != nil {
		return err
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

	err := fileutil.CreateDirIfNotExist(config.AdaptiveIP.PFBinatConfigFileLocation)
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(binatanchorConfigFileLocation, binatanchorConfData)
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

func createAndLoadNatAnchorConfig(privateIP string, publicIP string) error {
	var natanchorConfData string

	natanchorConfData = natStr
	natanchorConfData = strings.Replace(natanchorConfData, "HARP_INTERNAL_IFACE_NAME", config.AdaptiveIP.InternalIfaceName, -1)
	natanchorConfData = strings.Replace(natanchorConfData, "HARP_PF_PRIVATE_IP", privateIP, -1)
	natanchorConfData = strings.Replace(natanchorConfData, "HARP_PF_PUBLIC_IP", publicIP, -1)

	natanchorName := natanchorFilenamePrefix + publicIP
	logger.Logger.Println("createAndLoadNatAnchorConfig: Creating config file for " + natanchorName +
		" (publicIP: " + publicIP + ", privateIP: " + privateIP + ")")
	natanchorConfigFileLocation := config.AdaptiveIP.PFnatConfigFileLocation + "/" + natanchorName + ".conf"

	err := fileutil.CreateDirIfNotExist(config.AdaptiveIP.PFnatConfigFileLocation)
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(natanchorConfigFileLocation, natanchorConfData)
	if err != nil {
		return err
	}

	logger.Logger.Println("createAndLoadNatAnchorConfig: Load nat anchor config file for " + natanchorName)
	err = LoadPFAnchorRule(natanchorName, natanchorConfigFileLocation)
	if err != nil {
		return err
	}

	return nil
}

// CreateAndLoadAnchorConfig : Create anchor config files to match private IP address
// to available public IP address. Then load them to pf firewall.
func CreateAndLoadAnchorConfig(publicIP string, privateIP string) error {
	adaptiveip := config.GetAdaptiveIPNetwork()

	err := CheckBinatAnchorFileExist(publicIP)
	if err != nil {
		goto Error
	}

	err = checkDuplicatedIPAddress(publicIP)
	if err != nil {
		goto Error
	}

	err = createAndLoadBinatAnchorConfig(privateIP, publicIP)
	if err != nil {
		goto Error
	}

	err = createAndLoadNatAnchorConfig(privateIP, publicIP)
	if err != nil {
		goto Error
	}

	err = ifconfig.CreateAndLoadIfconfigScriptExternal(config.AdaptiveIP.ExternalIfaceName, publicIP,
		adaptiveip.Netmask)
	if err != nil {
		goto Error
	}

	return nil
Error:
	return err
}

func deleteAndUnloadBinatAnchorConfig(publicIP string) error {
	binatanchorName := binatanchorFilenamePrefix + publicIP
	logger.Logger.Println("deleteAndUnloadBinatAnchorConfig: Deleting config file of " + binatanchorName +
		" (publicIP: " + publicIP + ")")
	binatanchorConfigFileLocation := config.AdaptiveIP.PFBinatConfigFileLocation + "/" + binatanchorName + ".conf"

	err := fileutil.DeleteFile(binatanchorConfigFileLocation)
	if err != nil {
		logger.Logger.Println(err.Error())
	}

	logger.Logger.Println("deleteAndUnloadBinatAnchorConfig: Remove binat anchor rules of " + binatanchorName)
	err = removePFAnchorRule(binatanchorName)
	if err != nil {
		return err
	}

	return nil
}

func deleteAndUnloadNatAnchorConfig(publicIP string) error {
	natanchorName := natanchorFilenamePrefix + publicIP
	logger.Logger.Println("deleteAndUnloadNatAnchorConfig: Deleting config file of " + natanchorName +
		" (publicIP: " + publicIP + ")")
	natanchorConfigFileLocation := config.AdaptiveIP.PFnatConfigFileLocation + "/" + natanchorName + ".conf"

	err := fileutil.DeleteFile(natanchorConfigFileLocation)
	if err != nil {
		logger.Logger.Println(err.Error())
	}

	logger.Logger.Println("deleteAndUnloadNatAnchorConfig: Remove nat anchor rules of " + natanchorName)
	err = removePFAnchorRule(natanchorName)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAndUnloadAnchorConfig : Delete anchor config files to match public IP address.
// Then unload them from pf firewall.
func DeleteAndUnloadAnchorConfig(publicIP string) error {
	adaptiveip := config.GetAdaptiveIPNetwork()

	err := deleteAndUnloadBinatAnchorConfig(publicIP)
	if err != nil {
		goto Error
	}

	err = deleteAndUnloadNatAnchorConfig(publicIP)
	if err != nil {
		goto Error
	}

	err = ifconfig.DeleteAndUnloadIfconfigScriptExternal(config.AdaptiveIP.ExternalIfaceName, publicIP,
		adaptiveip.Netmask)
	if err != nil {
		goto Error
	}

	return nil
Error:
	return err
}
