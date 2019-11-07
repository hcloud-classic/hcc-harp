package adaptiveip

import (
	"bytes"
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/syscheck"
	"os/exec"
	"strconv"
	"strings"
)

func CheckDuplicatedIPAddress(IP string) error {
	logger.Logger.Println("Checking duplicated IP address for " + IP + " by running arping command...")

	if syscheck.LocalPublicIP == IP {
		return errors.New(IP + " is your local public ip address")
	}

	cmd := exec.Command("arping", "-i", config.AdaptiveIP.ExternalIfaceName, "-c",
		strconv.Itoa(int(config.AdaptiveIP.ArpingRetryCount)), IP)

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		logger.Logger.Println("arping: " + err.Error())
	}

	cmdOutputStr := string(cmdOutput.Bytes())
	if strings.Contains(cmdOutputStr, "from") {
		return errors.New("Found duplicated IP address for " + IP)
	}

	return nil
}
