package data

import "hcc/harp/model"

// AllServerData : Data structure of all_server
type AllServerData struct {
	Data struct {
		AllServer []model.Server `json:"all_server"`
	} `json:"data"`
}
