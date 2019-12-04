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

	configFiles, err := getBinatAnchorConfigFiles()
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

func loadExstingBinatAndNATRules() error {
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

// LoadPFAnchorRule : Load anchor rules configuration file and apply it to pf firewall
func LoadPFAnchorRule(anchorName string, anchorConfigFileLocation string) error {
	err := removePFAnchorRule(anchorName)
	if err != nil {
		return err
	}

	logger.Logger.Println("Loading anchor of " + anchorName + "...")

	cmd := exec.Command("pfctl", "-a", anchorName, "-f", anchorConfigFileLocation)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// removePFAnchorRule : Remove anchor rules of provided name from pf firewall
func removePFAnchorRule(anchorName string) error {
	logger.Logger.Println("Removing anchor rules of " + anchorName + "...")

	cmd := exec.Command("pfctl", "-a", anchorName, "-F", "all")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// LoadHarpPFRules : Load pf rules for harp module
func LoadHarpPFRules() error {
	err := loadExistingIfconfigScriptsInternal()
	if err != nil {
		return err
	}

	err = flushPFRules()
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

	err = loadExistingIfconfigScriptsExternal()
	if err != nil {
		return err
	}

	return nil
}
