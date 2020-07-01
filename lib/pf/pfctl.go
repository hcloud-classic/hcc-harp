package pf

import (
	"hcc/harp/lib/logger"
	"os/exec"
)

// FlushPFRules : Remove all of configured pf rules
func FlushPFRules() error {
	logger.Logger.Println("Flushing pf rules...")

	cmd := exec.Command("pfctl", "-F", "all")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// LoadPFRules : Load pf rules from config files created by harp module
func LoadPFRules(pfRulesFileLocation string) error {
	logger.Logger.Println("Loading pf rules from " + pfRulesFileLocation + "...")

	cmd := exec.Command("pfctl", "-f", pfRulesFileLocation)
	err := cmd.Run()
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
