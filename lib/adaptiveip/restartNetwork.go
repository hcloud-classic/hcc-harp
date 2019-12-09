package adaptiveip

import (
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

	cmd := exec.Command("service", "netif", "restart")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func restartNetwork() error {
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
