package adaptiveip

import (
	"hcc/harp/lib/config"
	"os/exec"
)

func flushPFRules() error {
	cmd := exec.Command("pfctl", "-F", "all")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func loadPFRules(pfRulesFileLocation string) error {
	cmd := exec.Command("pfctl", "-f", pfRulesFileLocation)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func LoadPFBinatAnchorRule(binatanchorName string, binatanchorConfigFileLocation string) error {
	cmd := exec.Command("pfctl", "-a", binatanchorName, "-f", binatanchorConfigFileLocation)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func RemvoePFBinatAnchorRule(binatanchorName string) error {
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

	return nil
}
