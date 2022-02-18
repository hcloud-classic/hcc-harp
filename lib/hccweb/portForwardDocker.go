package hccweb

import (
	"errors"
	"hcc/harp/daoext"
	"hcc/harp/lib/config"
	"hcc/harp/lib/iplinkext"
	"hcc/harp/lib/iptablesext"
	"os/exec"
	"strconv"
	"strings"
)

// PortForwardDocker : Add port forward rules for the hccweb docker container
func PortForwardDocker(isAdd bool, serverUUID string) error {
	subnet, errCode, errText := daoext.ReadSubnetByServer(serverUUID)
	if errCode != 0 {
		return errors.New(errText)
	}
	harpVNUM := iplinkext.GetIfaceVNUM(subnet.Gateway)

	cmd := exec.Command("sh", "-c",
		"docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' hccweb_"+strconv.Itoa(harpVNUM))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	dockerIP := strings.TrimSpace(string(output))

	adaptiveIPServer, errCode, errText := daoext.ReadAdaptiveIPServer(serverUUID)
	if errCode != 0 {
		return errors.New(errText)
	}

	err = iptablesext.PortForwarding(isAdd, false, true, false, adaptiveIPServer.PublicIP,
		dockerIP, int(config.Timpani.TimpaniExternalPort), int(config.Hccweb.Port))
	if err != nil {
		return err
	}

	return nil
}
