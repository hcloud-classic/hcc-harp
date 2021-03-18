package iptables

import (
	"bufio"
	"errors"
	"innogrid.com/hcloud-classic/pb"
	"hcc/harp/dao"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/iptablesext"
	"hcc/harp/lib/logger"
	"os"
	"os/exec"
	"strings"
)

func checkNFTables() error {
	var nfTablesMatched = 0
	var nfTablesOk = false

	logger.Logger.Println("Checking available tables for iptables...")

	file, err := os.Open("/proc/net/ip_tables_names")
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		for _, table := range iptablesext.NeededTablesForHarp {
			if strings.TrimSuffix(line, "\n") == table {
				nfTablesMatched++
				break
			}
		}
	}

	if len(iptablesext.NeededTablesForHarp) == nfTablesMatched {
		nfTablesOk = true
	}

	if !nfTablesOk {
		logger.Logger.Println("checkNFTables(): Some of tables are not available from iptables")
		logger.Logger.Println("checkNFTables(): Please check if your kernel modules are loaded properly")
		logger.Logger.Println("checkNFTables(): Type 'lsmod' and check if these modules are loaded: " + iptablesext.NeededKernelModulesForHarp)
		return errors.New("some of tables are not available from iptables")
	}

	return nil
}

func flushOrAddHarpIPTABLESChain(table string, chain string) error {
	// Check if the chain is exist then create the chain if not exist or flushing it if exist
	cmd := exec.Command("iptables", "-t", table, "-n", "-L", iptablesext.HarpChainNamePrefix+chain)
	err := cmd.Run()
	if err == nil {
		cmd = exec.Command("iptables", "-t", table, "-F", iptablesext.HarpChainNamePrefix+chain)
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("iptables", "-t", table, "-Z", iptablesext.HarpChainNamePrefix+chain)
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command("iptables", "-t", table, "-N", iptablesext.HarpChainNamePrefix+chain)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	// Check if the chain is included in the table then insert to first line of the table
	cmd = exec.Command("iptables", "-t", table, "-C", chain, "-j", iptablesext.HarpChainNamePrefix+chain)
	err = cmd.Run()
	if err == nil {
		cmd := exec.Command("iptables", "-t", table, "-D", chain, "-j", iptablesext.HarpChainNamePrefix+chain)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	cmd = exec.Command("iptables", "-t", table, "-I", chain, "1", "-j", iptablesext.HarpChainNamePrefix+chain)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func prepareHarpNATIPTABLESChains() error {
	logger.Logger.Println("Preparing harp's NAT iptables chains...")

	err := flushOrAddHarpIPTABLESChain("filter", "FORWARD")
	if err != nil {
		return err
	}

	for _, chain := range iptablesext.NatChains {
		err := flushOrAddHarpIPTABLESChain("nat", chain)
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

	err = checkNFTables()
	if err != nil {
		return err
	}

	err = prepareHarpNATIPTABLESChains()
	if err != nil {
		return err
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
			"-A", iptablesext.HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.ExternalIfaceName,
			"-s", adaptiveIPServer.PrivateIP,
			"-j", "SNAT",
			"--to-source", adaptiveIPServer.PublicIP)
		err := cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("iptables", "-t", "nat",
			"-A", iptablesext.HarpChainNamePrefix+"PREROUTING", "-i", config.AdaptiveIP.ExternalIfaceName,
			"-d", adaptiveIPServer.PublicIP,
			"-j", "DNAT",
			"--to-destination", adaptiveIPServer.PrivateIP)
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("iptables", "-t", "filter",
			"-A", iptablesext.HarpChainNamePrefix+"FORWARD",
			"-s", adaptiveIPServer.PublicIP,
			"-j", "ACCEPT")
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("iptables", "-t", "filter",
			"-A", iptablesext.HarpChainNamePrefix+"FORWARD",
			"-d", adaptiveIPServer.PrivateIP,
			"-j", "ACCEPT")
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
