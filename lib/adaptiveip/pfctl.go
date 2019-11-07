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

	folder := config.AdaptiveIP.PFServersConfigFileLocation
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func loadExstingBinatAnchorServersRules() error {
	logger.Logger.Println("Loading existing binat anchor servers rules...")

	configFiles, err := getBinatAnchorConfigFiles()
	if err != nil {
		return err
	}

	for _, filePath := range configFiles {
		if filePath == config.AdaptiveIP.PFServersConfigFileLocation {
			continue
		}

		binatanchorFileName := filePath[len(config.AdaptiveIP.PFServersConfigFileLocation+"/"):]
		if !strings.Contains(binatanchorFileName, binatanchorFilenamePrefix) ||
			!strings.Contains(binatanchorFileName, ".conf") {
			logger.Logger.Println("Wrong binat anchor filename: " + binatanchorFileName)
			logger.Logger.Println("Filename must be as '" + binatanchorFilenamePrefix + "XXX'")
			continue
		}

		binatanchorName := binatanchorFileName[0 : len(binatanchorFileName)-len(".conf")]
		err := LoadPFBinatAnchorRule(binatanchorName, filePath)
		if err != nil {
			logger.Logger.Println(err)
		}
	}

	return nil
}

func LoadPFBinatAnchorRule(binatanchorName string, binatanchorConfigFileLocation string) error {
	logger.Logger.Println("Loading binat anchor of " + binatanchorName + "...")

	IP := binatanchorName[len(binatanchorFilenamePrefix):]
	err := CheckDuplicatedIPAddress(IP)
	if err != nil {
		return err
	}

	cmd := exec.Command("pfctl", "-a", binatanchorName, "-f", binatanchorConfigFileLocation)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func RemvoePFBinatAnchorRule(binatanchorName string) error {
	logger.Logger.Println("Removing binat anchor rules of " + binatanchorName + "...")

	cmd := exec.Command("pfctl", "-a", binatanchorName, "-F", "all")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func LoadHarpPFRules() error {
	err := flushPFRules()
	if err != nil {
		return err
	}

	err = loadPFRules(config.AdaptiveIP.PFRulesFileLocation)
	if err != nil {
		return err
	}

	err = loadExstingBinatAnchorServersRules()
	if err != nil {
		return err
	}

	return nil
}
