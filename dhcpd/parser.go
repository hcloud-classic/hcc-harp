package dhcpd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	"hcc/harp/config"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type nodeEntries struct {
	PXEMACAddress string
	IP            string
	NodeName      string
}

func GetPXEFilename(os string) (string, error) {
	// TODO: Need implement of how to get pxe file location for each OS

	return "", errors.New("not supported os")
}

type nodePXEMACAddr struct {
	Data struct {
		Node struct {
			PxeMacAddr string `json:"pxe_mac_addr"`
		} `json:"node"`
	} `json:"data"`
}

func GetPXEMACAddress(nodeUUID string) (string, error) {
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

func ConfParser(networkIP string, netmask string, nodeUUIDs []string,
	leaderUUID string, os string, name string) error {
	var err error = nil

	var maskPartsStr = strings.Split(netmask, ".")
	var maskParts [4]int

	for i := range maskPartsStr {
		maskParts[i], err = strconv.Atoi(maskPartsStr[i])
		if err != nil {
			return err
		}
	}

	var mask = net.IPv4Mask(
		byte(maskParts[0]),
		byte(maskParts[1]),
		byte(maskParts[2]),
		byte(maskParts[3]))

	ipNet := net.IPNet{
		IP:   net.ParseIP(networkIP).To4(),
		Mask: mask,
	}

	count := int(cidr.AddressCount(&ipNet) - 2)

	if len(nodeUUIDs) > count {
		return errors.New("nodes count is bigger than available IP addresses count")
	}

	firstIP, lastIP := cidr.AddressRange(&ipNet)
	firstIP = cidr.Inc(firstIP)
	lastIP = cidr.Dec(lastIP)

	nextIP := cidr.Inc(firstIP)

	pxeFileName, err := GetPXEFilename(os)
	if err != nil {
		return err
	}

	strings.Replace(confBase, "HARP_DHCPD_PXE_FILENAME", pxeFileName, -1)
	strings.Replace(confBase, "HARP_DHCPD_DOMAIN_NAME", name, -1)
	strings.Replace(confBase, "HARP_DHCPD_MIN_LEASE_TIME", strconv.Itoa(int(config.DHCPD.MinLeaseTime)), -1)
	strings.Replace(confBase, "HARP_DHCPD_DEFAULT_LEASE_TIME", strconv.Itoa(int(config.DHCPD.DefaultLeaseTime)), -1)
	strings.Replace(confBase, "HARP_DHCPD_MAX_LEASE_TIME", strconv.Itoa(int(config.DHCPD.MaxLeaseTime)), -1)

	var nodeEntryConfPart = ""
	for i, uuid := range nodeUUIDs {
		if nextIP.Equal(lastIP) {
			return errors.New("ip range exceeded")
		}

		pxeMacAddr, err := GetPXEMACAddress(uuid)
		if err != nil {
			return err
		}

		var node = new(nodeEntries)
		node.NodeName = "node" + strconv.Itoa(i) + "." + name
		node.PXEMACAddress = pxeMacAddr
		if uuid == leaderUUID {
			node.IP = firstIP.String()
		} else {
			node.IP = nextIP.String()
			nextIP = cidr.Inc(nextIP)
		}

		var nodeConfPart = nodeEntry
		strings.Replace(nodeConfPart, "HARP_DHCPD_NODE_NAME", node.NodeName, -1)
		strings.Replace(nodeConfPart, "HARP_DHCPD_NODE_PXE_MAC", node.PXEMACAddress, -1)
		strings.Replace(nodeConfPart, "HARP_DHCPD_NODE_IP", node.IP, -1)

		nodeEntryConfPart += nodeConfPart
	}

	strings.Replace(confBase, "HARP_DHCPD_NODES_ENTRIES", nodeEntryConfPart, -1)

	// TODO: Write config string to file
	fmt.Println(confBase)

	return err
}
