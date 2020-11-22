package ifconfig

import (
	"hcc/harp/lib/logger"
	"hcc/harp/lib/syscheck"
	"os/exec"
)

func loadIfconfigScript(filepath string) error {
	logger.Logger.Println("Loading ifconfig script file: " + filepath)

	var cmd *exec.Cmd

	if syscheck.OS == "freebsd" {
		cmd = exec.Command("csh", filepath)
	} else {
		cmd = exec.Command("bash", filepath)
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
