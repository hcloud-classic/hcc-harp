package servicecontrol

import (
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpdext"
	"hcc/harp/lib/logger"
	"os/exec"
	"sync"
)

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
	output, err := cmd.CombinedOutput()
	if err != nil {
		dhcpdLock.Unlock()
		return errors.New(string(output))
	}

	dhcpdLock.Unlock()
	return nil
}
