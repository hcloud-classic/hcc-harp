package servicecontrol

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpdext"
	"hcc/harp/lib/logger"
	"os/exec"
	"sync"
)

func restartNetif() error {
	logger.Logger.Println("Restarting netif service...")

	cmd := exec.Command("service", "netif", "restart")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func restartRouting() error {
	logger.Logger.Println("Restarting routing service...")

	cmd := exec.Command("service", "routing", "restart")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

var dhcpdLock sync.Mutex

// RestartDHCPDServer : Run 'service isc-dhcpd restart' command to restart local dhcpd server
func RestartDHCPDServer() error {
	dhcpdLock.Lock()

	configFiles, err := dhcpdext.GetSubnetConfFiles()
	if err != nil {
		dhcpdLock.Unlock()
		return err
	}
	if len(configFiles) == 0 {
		logger.Logger.Println("No need to restart dhcpd service.")

		logger.Logger.Println("Stopping dhcpd service...")
		cmd := exec.Command("service", config.DHCPD.LocalDHCPDServiceName, "stop")
		_ = cmd.Run()

		dhcpdLock.Unlock()
		return nil
	}

	logger.Logger.Println("Restarting dhcpd service...")

	cmd := exec.Command("service", config.DHCPD.LocalDHCPDServiceName, "restart")
	err = cmd.Run()
	if err != nil {
		dhcpdLock.Unlock()
		return err
	}

	dhcpdLock.Unlock()
	return nil
}

// RestartNetwork : Restart network related services
func RestartNetwork() error {
	logger.Logger.Println("Restarting network services...")

	err := restartNetif()
	if err != nil {
		logger.Logger.Println(err)
	}

	err = restartRouting()
	if err != nil {
		logger.Logger.Println(err)
	}

	return nil
}
