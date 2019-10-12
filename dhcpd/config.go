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

func writeFile(input string, fileLocation string) error {
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

	var maskPartsStr = strings.Split(netmask, ".")
	if len(maskPartsStr) != 4 {
		return errors.New("netmask should be X.X.X.X form")
	}

	var maskParts [4]int
	for i := range maskPartsStr {
		maskParts[i], err = strconv.Atoi(maskPartsStr[i])
		if err != nil {
			return errors.New("netmask contained none integer value")
		}
	}

	var mask = net.IPv4Mask(
		byte(maskParts[0]),
		byte(maskParts[1]),
		byte(maskParts[2]),
		byte(maskParts[3]))

	maskSizeOne, maskSizeBit := mask.Size()
	if maskSizeOne == 0 && maskSizeBit == 0 {
		return errors.New("invalid netmask")
	}

	if maskSizeOne > 30 {
		return errors.New("netmask bit should be equal or smaller than 30")
	}

	ipNet := net.IPNet{
		IP:   netIPnetworkIP,
		Mask: mask,
	}

	netIPgateway := net.ParseIP(gateway)
	if netIPgateway == nil {
		return errors.New("wrong gateway IP")
	}
	isGatewayInSubnet := ipNet.Contains(netIPgateway)
	if isGatewayInSubnet == false {
		return errors.New("gateway IP is not in the subnet")
	}

	netIPnextServer := net.ParseIP(nextServer)
	if netIPnextServer == nil {
		return errors.New("wrong next server IP")
	}

	netIPnameServer := net.ParseIP(nameServer)
	if netIPnameServer == nil {
		return errors.New("wrong name server IP")
	}

	if len(nodeUUIDs) == 0 {
		return errors.New("provided nodeUUIDs[] is empty")
	}
	count := int(cidr.AddressCount(&ipNet) - 2)
	if len(nodeUUIDs) > maxNodes {
		return errors.New("nodes count is bigger than provided max nodes value")
	}
	if maxNodes > count {
		return errors.New("provided max nodes value is bigger than available IP addresses count")
	}

	firstIP, _ := cidr.AddressRange(&ipNet)
	firstIP = cidr.Inc(firstIP)
	lastIP := firstIP

	for i := 0; i < maxNodes-1; i++ {
		lastIP = cidr.Inc(lastIP)
	}

	pxeFileName, err := getPXEFilename(os)
	if err != nil {
		return err
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

	err = logger.CreateDirIfNotExist(config.DHCPD.ConfigFileLocation)
	if err != nil {
		return err
	}

	dhcpdConfLocation := config.DHCPD.ConfigFileLocation + "/" + name + ".conf"
	err = writeFile(confContent, dhcpdConfLocation)
	if err != nil {
		return err
	}

	return err
}
