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

func checkDuplicatedIPAddress(IP string) error {
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
