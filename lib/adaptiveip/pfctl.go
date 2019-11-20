package adaptiveip

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func flushPFRules() error {
	logger.Logger.Println("Flushing pf rules...")

	cmd := exec.Command("pfctl", "-F", "all")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func loadPFRules(pfRulesFileLocation string) error {
	logger.Logger.Println("Loading pf rules from " + pfRulesFileLocation + "...")

	cmd := exec.Command("pfctl", "-f", pfRulesFileLocation)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getBinatAnchorConfigFiles() ([]string, error) {
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

func loadExstingBinatAndNATRules() error {
	logger.Logger.Println("Loading existing binat and NAT rules...")

	configFiles, err := getBinatAnchorConfigFiles()
	if err != nil {
		return err
	}
	if len(configFiles) == 1 {
		return nil
	}

	ipMap := getAvailableIPsStatusMap()

	var binatanchorFileName string
	var binatanchorName string
	var ip string

	for i := 0; i < len(configFiles); i++ {
		if configFiles[i] == config.AdaptiveIP.PFBinatConfigFileLocation {
			continue
		}

		binatanchorFileName = configFiles[i][len(config.AdaptiveIP.PFBinatConfigFileLocation+"/"):]
		if !strings.Contains(binatanchorFileName, binatanchorFilenamePrefix) ||
			!strings.Contains(binatanchorFileName, ".conf") {
			logger.Logger.Println("Wrong binat anchor filename: " + binatanchorFileName)
			logger.Logger.Println("Filename must be as '" + binatanchorFilenamePrefix + "XXX'")
			continue
		}

		binatanchorName = binatanchorFileName[0 : len(binatanchorFileName)-len(".conf")]
		ip = binatanchorName[len(binatanchorFilenamePrefix):]
		if !ipMap[ip] {
			logger.Logger.Println("Skipping for not available IP address: " + ip)
			continue
		}
		err = LoadPFBinatAnchorRule(binatanchorName, configFiles[i])
		if err != nil {
			logger.Logger.Println(err)
		}
	}

	return nil
}

// LoadPFBinatAnchorRule : Load binat anchor rules configuration file and apply it to pf firewall
func LoadPFBinatAnchorRule(binatanchorName string, binatanchorConfigFileLocation string) error {
	logger.Logger.Println("Loading binat anchor of " + binatanchorName + "...")

	cmd := exec.Command("pfctl", "-a", binatanchorName, "-f", binatanchorConfigFileLocation)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// RemvoePFBinatAnchorRule : Remove binat anchor rules of provided name from pf firewall
func RemvoePFBinatAnchorRule(binatanchorName string) error {
	logger.Logger.Println("Removing binat anchor rules of " + binatanchorName + "...")

	cmd := exec.Command("pfctl", "-a", binatanchorName, "-F", "all")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// LoadHarpPFRules : Load pf rules for harp module
func LoadHarpPFRules() error {
	err := flushPFRules()
	if err != nil {
		return err
	}

	err = loadPFRules(config.AdaptiveIP.PFRulesFileLocation)
	if err != nil {
		return err
	}

	err = loadExstingBinatAndNATRules()
	if err != nil {
		return err
	}

	return nil
}
