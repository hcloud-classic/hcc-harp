package iptables

import (
	"bufio"
	"errors"
	"hcc/harp/dao"
	"hcc/harp/daoext"
	"hcc/harp/lib/config"
	"hcc/harp/lib/configadapriveipnetwork"
	"hcc/harp/lib/configadaptiveip"
	"hcc/harp/lib/iptablesext"
	"hcc/harp/lib/logger"
	"innogrid.com/hcloud-classic/pb"
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

func flushOrAddHarpIPTABLESChainAdaptiveIPInputDrop() error {
	// Check if the chain is exist then create the chain if not exist or flushing it if exist
	cmd := exec.Command("iptables", "-t", "filter", "-n", "-L", iptablesext.HarpAdaptiveIPInputDropChainName)
	err := cmd.Run()
	if err == nil {
		cmd = exec.Command("iptables", "-t", "filter", "-F", iptablesext.HarpAdaptiveIPInputDropChainName)
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("iptables", "-t", "filter", "-Z", iptablesext.HarpAdaptiveIPInputDropChainName)
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command("iptables", "-t", "filter", "-N", iptablesext.HarpAdaptiveIPInputDropChainName)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	// Check if the chain is included in the table then insert to first line of the table
	cmd = exec.Command("iptables", "-t", "filter", "-C", iptablesext.HarpChainNamePrefix+"INPUT",
		"-j", iptablesext.HarpAdaptiveIPInputDropChainName)
	err = cmd.Run()
	if err == nil {
		cmd := exec.Command("iptables", "-t", "filter", "-D", iptablesext.HarpChainNamePrefix+"INPUT",
			"-j", iptablesext.HarpAdaptiveIPInputDropChainName)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	cmd = exec.Command("iptables", "-t", "filter", "-A", iptablesext.HarpChainNamePrefix+"INPUT",
		"-j", iptablesext.HarpAdaptiveIPInputDropChainName)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func addHarpExternalMasqueradeRule() error {
	cmd := exec.Command("iptables", "-t", "nat",
		"-C", iptablesext.HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.ExternalIfaceName,
		"-j", "MASQUERADE")
	err := cmd.Run()
	isExist := err == nil

	if !isExist {
		cmd = exec.Command("iptables", "-t", "nat",
			"-A", iptablesext.HarpChainNamePrefix+"POSTROUTING", "-o", config.AdaptiveIP.ExternalIfaceName,
			"-j", "MASQUERADE")
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to add MASQUERADE rule for external interface")
		}
	}

	return nil
}

func addNATSecurityRule() error {
	logger.Logger.Println("Adding NAT security rule...")

	cmd := exec.Command("iptables", "-t", "nat",
		"-C", "PREROUTING",
		"-i", config.AdaptiveIP.ExternalIfaceName,
		"-j", "RETURN")
	err := cmd.Run()
	isExist := err == nil

	if isExist {
		cmd = exec.Command("iptables", "-t", "nat",
			"-D", "PREROUTING",
			"-i", config.AdaptiveIP.ExternalIfaceName,
			"-j", "RETURN")
		err = cmd.Run()
		if err != nil {
			return errors.New("failed to delete NAT security rule")
		}
	}

	cmd = exec.Command("iptables", "-t", "nat",
		"-I", "PREROUTING", "2",
		"-i", config.AdaptiveIP.ExternalIfaceName,
		"-j", "RETURN")
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to add NAT security rule")
	}

	return nil
}

func prepareHarpIPTABLESChains() error {
	logger.Logger.Println("Preparing harp's iptables chains...")

	err := flushOrAddHarpIPTABLESChain("filter", "INPUT")
	if err != nil {
		return err
	}

	err = flushOrAddHarpIPTABLESChainAdaptiveIPInputDrop()
	if err != nil {
		return err
	}

	err = flushOrAddHarpIPTABLESChain("filter", "FORWARD")
	if err != nil {
		return err
	}

	for _, chain := range iptablesext.NatChains {
		err := flushOrAddHarpIPTABLESChain("nat", chain)
		if err != nil {
			return err
		}
	}

	err = addNATSecurityRule()
	if err != nil {
		return err
	}

	err = addHarpExternalMasqueradeRule()
	if err != nil {
		return err
	}

	return nil
}

// InitIPTABLES : Prepare for use iptables
func InitIPTABLES() error {
	logger.Logger.Println("Initializing iptables...")

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()

	err := configadaptiveip.CheckAdaptiveIPConfig(adaptiveIP)
	if err != nil {
		return err
	}

	err = checkNFTables()
	if err != nil {
		return err
	}

	err = prepareHarpIPTABLESChains()
	if err != nil {
		return err
	}

	err = masterInputControl()
	if err != nil {
		return err
	}

	return nil
}

// LoadAdaptiveIPNetDevAndIPTABLESRules : Load interfaces and iptables rules for AdaptiveIP
func LoadAdaptiveIPNetDevAndIPTABLESRules() error {
	logger.Logger.Println("Loading interfaces and iptables rules for AdaptiveIPs...")

	var adaptiveIPServer pb.AdaptiveIPServer
	in := &pb.ReqGetAdaptiveIPServerList{
		AdaptiveipServer: &adaptiveIPServer,
		Row:              0,
		Page:             0,
	}

	adaptiveIPServerList, errCode, errString := daoext.ReadAdaptiveIPServerList(in)
	if errCode != 0 {
		return errors.New(errString)
	}

	for _, adaptiveIPServer := range adaptiveIPServerList.AdaptiveipServer {
		err := iptablesext.ControlNetDevAndIPTABLES(true, adaptiveIPServer.PublicIP,
			adaptiveIPServer.PrivateIP)
		if err != nil {
			logger.Logger.Println(err.Error())
		}

		portForwardingList, _, _ := dao.ReadPortForwardingList(&pb.ReqGetPortForwardingList{
			PortForwarding: &pb.PortForwarding{
				ServerUUID: adaptiveIPServer.ServerUUID,
			},
		})
		if portForwardingList != nil {
			for _, portForwarding := range portForwardingList.PortForwarding {
				err = iptablesext.PortForwarding(true, false, portForwarding.ForwardTCP, portForwarding.ForwardUDP,
					adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP,
					int(portForwarding.ExternalPort), int(portForwarding.InternalPort))
				if err != nil {
					logger.Logger.Println(err.Error())
				}
			}
		}
	}

	portForwardingList, _, _ := dao.ReadPortForwardingList(&pb.ReqGetPortForwardingList{
		PortForwarding: &pb.PortForwarding{
			ServerUUID: "master",
		},
	})
	if portForwardingList != nil {
		adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()

		for _, masterInput := range portForwardingList.PortForwarding {
			err := iptablesext.PortForwarding(true, false, masterInput.ForwardTCP, masterInput.ForwardUDP,
				adaptiveIP.ExtIfaceIPAddress, "",
				int(masterInput.ExternalPort), 0)
			if err != nil {
				logger.Logger.Println(err.Error())
			}
		}
	}

	return nil
}
