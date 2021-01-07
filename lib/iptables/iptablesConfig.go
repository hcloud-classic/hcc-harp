package iptables

import (
	"bufio"
	"errors"
	"github.com/hcloud-classic/pb"
	"hcc/harp/dao"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/logger"
	"os"
	"os/exec"
)

func getNFTables() ([]string, error) {
	logger.Logger.Println("Checking available tables for iptables...")

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

// InitIPTABLES : Prepare for use iptables
func InitIPTABLES() error {
	logger.Logger.Println("Initializing iptables...")

	adaptiveIP := configext.GetAdaptiveIPNetwork()

	err := configext.CheckAdaptiveIPConfig(adaptiveIP)
	if err != nil {
		return err
	}

	err = FlushIPTABLESRules()
	if err != nil {
		return err
	}

	logger.Logger.Println("Restoring iptables rules from " + config.AdaptiveIP.IPTABLESInitConfigFileLocation + "...")
	cmd := exec.Command("iptables-restore", config.AdaptiveIP.IPTABLESInitConfigFileLocation)
	err = cmd.Run()
	if err != nil {
		logger.Logger.Println("Failed to restore iptables rules. Skipping...")
	}

	return nil
}

// LoadAdaptiveIPIPTABLESRules : Load iptables rules for AdaptiveIP
func LoadAdaptiveIPIPTABLESRules() error {
	logger.Logger.Println("Loading iptables rules for AdaptiveIP...")

	var adaptiveIPServer pb.AdaptiveIPServer
	in := &pb.ReqGetAdaptiveIPServerList{
		AdaptiveipServer: &adaptiveIPServer,
		Row:              0,
		Page:             0,
	}

	adaptiveIPServerList, errCode, errString := dao.ReadAdaptiveIPServerList(in)
	if errCode != 0 {
		return errors.New(errString)
	}

	for _, adaptiveIPServer := range adaptiveIPServerList.AdaptiveipServer {
		cmd := exec.Command("iptables", "-t", "nat",
			"-A", "POSTROUTING", "-o", config.AdaptiveIP.ExternalIfaceName,
			"-s", adaptiveIPServer.PrivateIP,
			"-j", "SNAT",
			"--to-source", adaptiveIPServer.PublicIP)
		err := cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("iptables", "-t", "nat",
			"-A", "PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
			"-d", adaptiveIPServer.PublicIP,
			"-j", "DNAT",
			"--to-destination", adaptiveIPServer.PrivateIP)
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("iptables",
			"-A", "FORWARD",
			"-s", adaptiveIPServer.PublicIP,
			"-j", "ACCEPT")
		err = cmd.Run()

		cmd = exec.Command("iptables",
			"-A", "FORWARD",
			"-d", adaptiveIPServer.PrivateIP,
			"-j", "ACCEPT")
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
