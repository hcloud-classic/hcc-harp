package ifconfig

import (
	"hcc/harp/lib/logger"
	"os/exec"
)

func loadIfconfigScript(filepath string) error {
	logger.Logger.Println("Loading ifconfig script file: " + filepath)

	cmd := exec.Command("csh", filepath)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
