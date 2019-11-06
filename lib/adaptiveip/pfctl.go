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

	for _, filename := range configFiles {
		if filename == config.AdaptiveIP.PFServersConfigFileLocation {
			continue
		}

		if !strings.Contains(filename, binatanchorFilenamePrefix) ||
			!strings.Contains(filename, ".conf") {
			logger.Logger.Println("Wrong binat anchor filename: " + filename)
			logger.Logger.Println("Filename must be as '" + binatanchorFilenamePrefix + "XXX'")
			continue
		}

		binatanchorName := filename[0 : len(filename)-len(".conf")]
		err := LoadPFBinatAnchorRule(binatanchorName, config.AdaptiveIP.PFServersConfigFileLocation+"/"+filename)
		if err != nil {
			logger.Logger.Println(err)
		}
	}

	return nil
}

func LoadPFBinatAnchorRule(binatanchorName string, binatanchorConfigFileLocation string) error {
	logger.Logger.Println("Loading binat anchor of " + binatanchorName + "...")

	cmd := exec.Command("pfctl", "-a", binatanchorName, "-f", binatanchorConfigFileLocation)
	err := cmd.Run()
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
