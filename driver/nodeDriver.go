package driver

import (
	"hcc/harp/data"
	"hcc/harp/http"
)

// ListNode : Get the list of nodes by server UUID.
func ListNode(serverUUID string) (interface{}, error) {
	arguments := "server_uuid:\"" + serverUUID + "\", active:1"

	var listNodeData data.ListNodeData
	query := "query { list_node(" + arguments + ") { uuid } }"

	result, err :=  http.DoHTTPRequest("flute", true, listNodeData, query)
	if err != nil {
		return listNodeData, err
	}

	return result, nil
}