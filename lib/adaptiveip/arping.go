package adaptiveip

import (
	"bytes"
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"os/exec"
	"strconv"
	"strings"
)

// CheckDuplicatedIPAddress : Check duplicated IP address by running arping command
func CheckDuplicatedIPAddress(IP string) error {
	logger.Logger.Println("Checking duplicated IP address for " + IP + " by running arping command...")

	if localPublicIP == IP {
		return errors.New(IP + " is your local public ip address")
	}

	cmd := exec.Command("arping", "-i", config.AdaptiveIP.ExternalIfaceName, "-c",
		strconv.Itoa(int(config.AdaptiveIP.ArpingRetryCount)), IP)

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	cmdOutputStr := string(cmdOutput.Bytes())

	if strings.Contains(cmdOutputStr, "Timeout") ||
		strings.Contains(cmdOutputStr, "timeout") {
		return nil
	}

	if err != nil {
		logger.Logger.Println("arping: " + err.Error())
	}

	if strings.Contains(cmdOutputStr, "from") {
		return errors.New("Found duplicated IP address for " + IP)
	}

	return nil
}
