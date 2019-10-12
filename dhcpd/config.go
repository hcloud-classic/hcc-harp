package dhcpd

import (
	"encoding/json"
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"hcc/harp/config"
	"hcc/harp/logger"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type nodeEntries struct {
	PXEMACAddress string
	IP            string
	NodeName      string
}

func getPXEFilename(os string) (string, error) {
	// TODO: Need implement of how to get pxe file location for each OS
	if os == "CentOS 6" {
		return "/boot/pxeboot/centos6", nil
	}

	return "", errors.New("not supported os")
}

type nodePXEMACAddr struct {
	Data struct {
		Node struct {
			PxeMacAddr string `json:"pxe_mac_addr"`
		} `json:"node"`
	} `json:"data"`
}

func getPXEMACAddress(nodeUUID string) (string, error) {
	client := &http.Client{Timeout: time.Duration(config.Flute.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "http://"+config.Flute.ServerAddress+":"+strconv.Itoa(int(config.Flute.ServerPort))+"/graphql?query={node(uuid:%22"+
		nodeUUID+"%22){pxe_mac_addr}}", nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var MACAddr nodePXEMACAddr
			err = json.Unmarshal([]byte(str), &MACAddr)
			if err != nil {
				return "", err
			}

			return MACAddr.Data.Node.PxeMacAddr, nil
		}

		return "", err
	}

	return "", errors.New("http response returned error code")
}

func writeFile(fileLocation string, input string) error {
	file, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.WriteString(input)
	if err != nil {
		return err
	}

	return nil
}

func writeConfigFile(input string, name string) error {
	err := logger.CreateDirIfNotExist(config.DHCPD.ConfigFileLocation)
	if err != nil {
		return err
	}

	dhcpdConfLocation := config.DHCPD.ConfigFileLocation + "/" + name + ".conf"
	err = writeFile(dhcpdConfLocation, input)
	if err != nil {
		return err
	}

	return nil
}

func checkNetmask(netmask string) (net.IPMask, error) {
	var err error

	var maskPartsStr = strings.Split(netmask, ".")
	if len(maskPartsStr) != 4 {
		return nil, errors.New("netmask should be X.X.X.X form")
	}

	var maskParts [4]int
	for i := range maskPartsStr {
		maskParts[i], err = strconv.Atoi(maskPartsStr[i])
		if err != nil {
			return nil, errors.New("netmask contained none integer value")
		}
	}

	var mask = net.IPv4Mask(
		byte(maskParts[0]),
		byte(maskParts[1]),
		byte(maskParts[2]),
		byte(maskParts[3]))

	maskSizeOne, maskSizeBit := mask.Size()
	if maskSizeOne == 0 && maskSizeBit == 0 {
		return nil, errors.New("invalid netmask")
	}

	if maskSizeOne > 30 {
		return nil, errors.New("netmask bit should be equal or smaller than 30")
	}

	return mask, err
}

func checkGateway(subnet net.IPNet, gateway string) error {
	netIPgateway := net.ParseIP(gateway)
	if netIPgateway == nil {
		return errors.New("wrong gateway IP")
	}
	isGatewayInSubnet := subnet.Contains(netIPgateway)
	if isGatewayInSubnet == false {
		return errors.New("gateway IP is not in the subnet")
	}

	return nil
}

func checkNodeUUIDs(subnet net.IPNet, maxNodes int, nodeUUIDs []string, leaderNodeUUID string) error {
	if len(nodeUUIDs) == 0 {
		return errors.New("provided nodeUUIDs[] is empty")
	}
	count := int(cidr.AddressCount(&subnet) - 2)
	if len(nodeUUIDs) > maxNodes {
		return errors.New("nodes count is bigger than provided max nodes value")
	}
	if maxNodes > count {
		return errors.New("provided max nodes value is bigger than available IP addresses count")
	}

	var leaderNodeUUIDfound = false
	for _, uuid := range nodeUUIDs {
		if uuid == leaderNodeUUID {
			leaderNodeUUIDfound = true
			break
		}
	}
	if leaderNodeUUIDfound == false {
		return errors.New("leaderNodeUUID is not found from provided nodeUUIDs[]")
	}

	return nil
}

// CreateConfig : Get needed parameters for make dhcpd config file then generate config file for each subnet
func CreateConfig(networkIP string, netmask string, gateway string,
	nextServer string, nameServer string,
	domainName string, maxNodes int, nodeUUIDs []string,
	leaderNodeUUID string, os string, name string) error {
	var err error

	if len(name) == 0 {
		return errors.New("name is needed for make dhcpd config file")
	}

	netIPnetworkIP := net.ParseIP(networkIP).To4()
	if netIPnetworkIP == nil {
		return errors.New("wrong network IP")
	}

	mask, err := checkNetmask(netmask)
	if err != nil {
		return err
	}

	ipNet := net.IPNet{
		IP:   netIPnetworkIP,
		Mask: mask,
	}

	err = checkGateway(ipNet, gateway)
	if err != nil {
		return err
	}

	netIPnextServer := net.ParseIP(nextServer)
	if netIPnextServer == nil {
		return errors.New("wrong next server IP")
	}

	netIPnameServer := net.ParseIP(nameServer)
	if netIPnameServer == nil {
		return errors.New("wrong name server IP")
	}

	err = checkNodeUUIDs(ipNet, maxNodes, nodeUUIDs, leaderNodeUUID)
	if err != nil {
		return err
	}

	pxeFileName, err := getPXEFilename(os)
	if err != nil {
		return err
	}

	firstIP, _ := cidr.AddressRange(&ipNet)
	firstIP = cidr.Inc(firstIP)
	lastIP := firstIP

	for i := 0; i < maxNodes-1; i++ {
		lastIP = cidr.Inc(lastIP)
	}

	confContent := confBase
	confContent = strings.Replace(confContent, "HARP_DHCPD_SUBNET", networkIP, -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_NETMASK", netmask, -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_START_IP", firstIP.String(), -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_LAST_IP", lastIP.String(), -1)

	confContent = strings.Replace(confContent, "HARP_DHCPD_NEXT_SERVER", nextServer, -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_PXE_FILENAME", pxeFileName, -1)

	confContent = strings.Replace(confContent, "HARP_DHCPD_DOMAIN_NAME_SERVER", nameServer, -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_DOMAIN_NAME", domainName, -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_GATEWAY", gateway, -1)

	confContent = strings.Replace(confContent, "HARP_DHCPD_MIN_LEASE_TIME", strconv.Itoa(int(config.DHCPD.MinLeaseTime)), -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_DEFAULT_LEASE_TIME", strconv.Itoa(int(config.DHCPD.DefaultLeaseTime)), -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_MAX_LEASE_TIME", strconv.Itoa(int(config.DHCPD.MaxLeaseTime)), -1)

	nextIP := cidr.Inc(firstIP)

	var nodeEntryConfPart = ""
	for i, uuid := range nodeUUIDs {
		pxeMacAddr, err := getPXEMACAddress(uuid)
		if err != nil {
			return err
		}
		pxeMacAddr = strings.Replace(pxeMacAddr, "-", ":", -1)

		var node = new(nodeEntries)
		node.NodeName = "node" + strconv.Itoa(i) + "." + name
		node.PXEMACAddress = pxeMacAddr
		if uuid == leaderNodeUUID {
			node.IP = firstIP.String()
		} else {
			node.IP = nextIP.String()
			nextIP = cidr.Inc(nextIP)
		}

		var nodeConfPart = nodeEntry
		nodeConfPart = strings.Replace(nodeConfPart, "HARP_DHCPD_NODE_NAME", node.NodeName, -1)
		nodeConfPart = strings.Replace(nodeConfPart, "HARP_DHCPD_NODE_PXE_MAC", node.PXEMACAddress, -1)
		nodeConfPart = strings.Replace(nodeConfPart, "HARP_DHCPD_NODE_IP", node.IP, -1)

		nodeEntryConfPart += nodeConfPart
	}

	confContent = strings.Replace(confContent, "HARP_DHCPD_NODES_ENTRIES", nodeEntryConfPart, -1)

	err = writeConfigFile(confContent, name)
	if err != nil {
		return err
	}

	return err
}

// CheckLocalDHCPDConfig : Check if harp dhcpd config file is included in local dhcpd server config file
func CheckLocalDHCPDConfig() error {
	include := includeStr
	include = strings.Replace(include, "HARP_DHCPD_CONF_LOCATION",
		config.DHCPD.ConfigFileLocation+"/harp_dhcpd.conf", -1)

	data, err := ioutil.ReadFile(config.DHCPD.LocalConfigFileLocation)
	if err != nil {
		return errors.New("failed reading data from local dhcpd config file location")
	}

	isHarpDHCPDIncluded := strings.Contains(string(data), include)
	if !isHarpDHCPDIncluded {
		logger.Logger.Println("Please add this line to dhcpd config file!\n" + include)
		return errors.New("cannot find harp dhcp config include line from local dhcpd config file")
	}

	return nil
}

func getSubnetConfFiles() ([]string, error) {
	var files []string

	folder := config.DHCPD.ConfigFileLocation
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// UpdateHarpDHCPDConfig : Update harp dhcpd main config file. Write subnet config files include lines to 'harp_dhcpd.conf'
func UpdateHarpDHCPDConfig() error {
	configFiles, err := getSubnetConfFiles()
	if err != nil {
		return err
	}

	var allIncludeLines = ""
	for _, filename := range configFiles {
		if strings.Contains(filename, "harp_dhcpd.conf") ||
			filename == config.DHCPD.ConfigFileLocation {
			continue
		}

		include := includeStr
		include = strings.Replace(include, "HARP_DHCPD_CONF_LOCATION", filename, -1)
		allIncludeLines += include + "\n"
	}

	err = logger.CreateDirIfNotExist(config.DHCPD.ConfigFileLocation)
	if err != nil {
		return err
	}

	err = writeFile(config.DHCPD.ConfigFileLocation+"/harp_dhcpd.conf", allIncludeLines)
	if err != nil {
		return err
	}

	return nil
}

func RestartDHCPDServer() error {
	cmd := exec.Command("service", "isc-dhcpd", "restart")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
