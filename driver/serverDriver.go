package driver

import (
	"hcc/harp/data"
	"hcc/harp/http"
)

// AllServerUUID : Get all of server UUIDs.
func AllServerUUID() (interface{}, error) {
	var query = "query { all_server { uuid } }"

	var allServerData data.AllServerData

	result, err := http.DoHTTPRequest("violin", true, allServerData, query)
	if err != nil {
		return allServerData, err
	}

	return result, nil
}