package dhcpd

import (
	"encoding/json"
	"hcc/harp/dao"
	"hcc/harp/lib/fileutil"
	"net/http"
	"time"

	// "encoding/json"
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/model"
	"io/ioutil"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"

	// "net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	// "time"
)

type nodeEntries struct {
	PXEMACAddress string
	IP            string
	NodeName      string
}

func getPXEFilename(serverUUID string) (string, error) {
	if len(serverUUID) == 0 {
		return "", errors.New("failed to get server UUID")
	}
	return model.DefaultPXEdir + "/" + serverUUID + "/pxelinux.0", nil
}

//type nodePXEMACAddr struct {
//	Data struct {
//		Node struct {
//			PxeMacAddr string `json:"pxe_mac_addr"`
//		} `json:"node"`
//	} `json:"data"`
//}

/*
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
*/

// Flute

// NodeData : Data structure of list_node
type NodeData struct {
	Data struct {
		Node model.Node `json:"node"`
	} `json:"data"`
}

func getNodePXEMACAddress(nodeUUID string) (NodeData, error) {
	var nodePXEMACAddressData NodeData

	client := &http.Client{Timeout: time.Duration(config.Flute.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "http://"+config.Flute.ServerAddress+":"+strconv.Itoa(int(config.Flute.ServerPort))+
		"/graphql?query=query%20%7B%0A%20%20node(uuid%3A%20%22"+nodeUUID+"%22)%20%7B%0A%20%20%20%20pxe_mac_addr%0A%20%20%7D%0A%7D", nil)

	if err != nil {
		return nodePXEMACAddressData, err
	}
	resp, err := client.Do(req)

	if err != nil {
		return nodePXEMACAddressData, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			err = json.Unmarshal([]byte(str), &nodePXEMACAddressData)
			if err != nil {
				return nodePXEMACAddressData, err
			}

			return nodePXEMACAddressData, nil
		}

		return nodePXEMACAddressData, err
	}

	return nodePXEMACAddressData, errors.New("http response returned error code")
}

func getPXEMACAddress(nodeUUID string) (string, error) {
	nodePXEMACAddress, err := getNodePXEMACAddress(nodeUUID)
	if err != nil {
		return "", err
	}

	return nodePXEMACAddress.Data.Node.PXEMacAddr, nil
}

func writeConfigFile(input string, name string) error {
	err := logger.CreateDirIfNotExist(config.DHCPD.ConfigFileLocation)
	if err != nil {
		return err
	}

	dhcpdConfLocation := config.DHCPD.ConfigFileLocation + "/" + name + ".conf"
	err = fileutil.WriteFile(dhcpdConfLocation, input)
	if err != nil {
		return err
	}

	return nil
}

// CheckNodeUUIDs : Check UUIDs of nodes
func CheckNodeUUIDs(subnet net.IPNet, nodeUUIDs []string, leaderNodeUUID string) error {
	if len(nodeUUIDs) == 0 {
		return errors.New("provided nodeUUIDs[] is empty")
	}
	count := int(cidr.AddressCount(&subnet) - 2)
	if len(nodeUUIDs) > count {
		return errors.New("node count is bigger than available IP addresses count")
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

func doWriteConfig(subnet model.Subnet, firstIP net.IP, lastIP net.IP, pxeFileName string, nodeUUIDs []string,
				   useSamePXEFileForCompute bool) error {
	confContent := confBase
	confContent = strings.Replace(confContent, "HARP_DHCPD_SUBNET", subnet.NetworkIP, -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_NETMASK", subnet.Netmask, -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_START_IP", firstIP.String(), -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_LAST_IP", lastIP.String(), -1)

	confContent = strings.Replace(confContent, "HARP_DHCPD_NEXT_SERVER", subnet.NextServer, -1)
	if useSamePXEFileForCompute {
		confContent = strings.Replace(confContent, "HARP_DHCPD_PXE_FILENAME", pxeFileName, -1)
	} else {
		confContent = strings.Replace(confContent, "    filename \"HARP_DHCPD_PXE_FILENAME\";", "", -1)
	}

	confContent = strings.Replace(confContent, "HARP_DHCPD_DOMAIN_NAME_SERVER", subnet.NameServer, -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_DOMAIN_NAME", subnet.DomainName, -1)
	confContent = strings.Replace(confContent, "HARP_DHCPD_GATEWAY", subnet.Gateway, -1)

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
		node.NodeName = "node" + strconv.Itoa(i) + "." + subnet.SubnetName
		node.PXEMACAddress = pxeMacAddr
		if uuid == subnet.LeaderNodeUUID {
			node.IP = firstIP.String()
		} else {
			node.IP = nextIP.String()
			nextIP = cidr.Inc(nextIP)
		}

		var nodeConfPart = nodeEntry
		nodeConfPart = strings.Replace(nodeConfPart, "HARP_DHCPD_NODE_NAME", node.NodeName, -1)
		nodeConfPart = strings.Replace(nodeConfPart, "HARP_DHCPD_NODE_PXE_MAC", node.PXEMACAddress, -1)
		nodeConfPart = strings.Replace(nodeConfPart, "HARP_DHCPD_NODE_IP", node.IP, -1)
		if useSamePXEFileForCompute {
			nodeConfPart = strings.Replace(nodeConfPart, "        filename \"HARP_DHCPD_PXE_FILENAME\";", "", -1)
		} else {
			var otherPXEFileName = pxeFileName
			if uuid == subnet.LeaderNodeUUID {
				otherPXEFileName = strings.Replace(pxeFileName, subnet.ServerUUID, subnet.ServerUUID + "/Leader", -1)
			} else {
				otherPXEFileName = strings.Replace(pxeFileName, subnet.ServerUUID, subnet.ServerUUID + "/Compute", -1)
			}
			nodeConfPart = strings.Replace(nodeConfPart, "HARP_DHCPD_PXE_FILENAME", otherPXEFileName, -1)
		}

		nodeEntryConfPart += nodeConfPart
	}

	confContent = strings.Replace(confContent, "HARP_DHCPD_NODES_ENTRIES", nodeEntryConfPart, -1)

	err := writeConfigFile(confContent, subnet.ServerUUID)
	if err != nil {
		return err
	}

	return nil
}

// CreateConfig : Get needed parameters for make dhcpd config file then generate config file for each subnet
func CreateConfig(subnetUUID string, nodeUUIDs []string) error {
	if len(subnetUUID) == 0 {
		return errors.New("subnetUUID is needed for make dhcpd config file")
	}

	args := make(map[string]interface{})
	args["uuid"] = subnetUUID

	subnetInterface, err := dao.ReadSubnet(args)
	if err != nil {
		return err
	}

	var subnet = subnetInterface.(model.Subnet)

	if len(subnet.SubnetName) == 0 {
		return errors.New("name is needed for make dhcpd config file")
	}

	netIPnetworkIP := iputil.CheckValidIP(subnet.NetworkIP)
	if netIPnetworkIP == nil {
		return errors.New("wrong network IP address")
	}

	mask, err := iputil.CheckNetmask(subnet.Netmask)
	if err != nil {
		return err
	}

	ipNet := net.IPNet{
		IP:   netIPnetworkIP,
		Mask: mask,
	}

	err = iputil.CheckGateway(ipNet, subnet.Gateway)
	if err != nil {
		return err
	}

	netIPnextServer := iputil.CheckValidIP(subnet.NextServer)
	if netIPnextServer == nil {
		return errors.New("wrong next server IP")
	}

	netIPnameServer := iputil.CheckValidIP(subnet.NameServer)
	if netIPnameServer == nil {
		return errors.New("wrong name server IP")
	}

	err = CheckNodeUUIDs(ipNet, nodeUUIDs, subnet.LeaderNodeUUID)
	if err != nil {
		return err
	}

	pxeFileName, err := getPXEFilename(subnet.ServerUUID)
	if err != nil {
		return err
	}

	firstIP, _ := cidr.AddressRange(&ipNet)
	firstIP = cidr.Inc(firstIP)
	lastIP := firstIP

	for i := 0; i < len(nodeUUIDs)-1; i++ {
		lastIP = cidr.Inc(lastIP)
	}

	err = doWriteConfig(subnet, firstIP, lastIP, pxeFileName, nodeUUIDs, false)
	if err != nil {
		return err
	}

	//logger.Logger.Println("CreateConfig: Registering first IP to adaptive IP...")
	//err = adaptiveip.CreateAndLoadBinatAnchorConfig(firstIP.String())
	//if err != nil {
	//	return err
	//}

	return nil
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

	err = fileutil.WriteFile(config.DHCPD.ConfigFileLocation+"/harp_dhcpd.conf", allIncludeLines)
	if err != nil {
		return err
	}

	return nil
}

// RestartDHCPDServer : Run 'service isc-dhcpd restart' command to restart local dhcpd server
func RestartDHCPDServer() error {
	cmd := exec.Command("service", config.DHCPD.LocalDHCPDServiceName, "restart")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
