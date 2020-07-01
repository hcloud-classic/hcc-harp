package serviceControl

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"os/exec"
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

// RestartDHCPDServer : Run 'service isc-dhcpd restart' command to restart local dhcpd server
func RestartDHCPDServer() error {
	logger.Logger.Println("Restarting dhcpd service...")

	cmd := exec.Command("service", config.DHCPD.LocalDHCPDServiceName, "restart")
	err := cmd.Run()
	if err != nil {
		return err
	}

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

	err = RestartDHCPDServer()
	if err != nil {
		logger.Logger.Println(err)
	}

	return nil
}