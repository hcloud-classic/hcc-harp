package iptables

import (
	"bufio"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/logger"
	"os"
	"os/exec"
)

// InitIPTABLES : Prepare for use iptables
func InitIPTABLES() error {
	adaptiveIP := configext.GetAdaptiveIPNetwork()

	err := configext.CheckAdaptiveIPConfig(adaptiveIP)
	if err != nil {
		return err
	}

	//err = ReplaceBaseConfigAnchorStrings()
	//if err != nil {
	//	return err
	//}

	return nil
}

func getNFTables() ([]string, error) {
	var nfTables []string
	file, err := os.Open("/proc/net/ip_tables_names")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		nfTables = append(nfTables, line)
	}

	return nfTables, nil
}

// FlushIPTABLESRules : Remove all of configured iptables rules
func FlushIPTABLESRules() error {
	logger.Logger.Println("Flushing iptables rules...")

	cmd := exec.Command("iptables", "-F")
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("iptables", "-X")
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("iptables", "-Z")
	err = cmd.Run()
	if err != nil {
		return err
	}

	nfTables, err := getNFTables()
	if err != nil {
		return err
	}

	for _, table := range nfTables {
		cmd = exec.Command("iptables", "-t", table, "-F")
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("iptables", "-t", table, "-X")
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("iptables", "-t", table, "-Z")
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
