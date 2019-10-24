package graphql

import (
	"encoding/json"
	"errors"
	"hcc/harp/lib/config"
	"hcc/harp/model"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Flute

// NodeData : Data structure of list_node
type NodeData struct {
	Data struct {
		Node model.Node `json:"node"`
	} `json:"data"`
}

// GetNodePXEMACAddress : Get mac address of node
func GetNodePXEMACAddress(nodeUUID string) (NodeData, error) {
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
