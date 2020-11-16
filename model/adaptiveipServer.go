package model

import "time"

// AdaptiveIPServer : Struct of AdaptiveIPServer
type AdaptiveIPServer struct {
	ServerUUID     string    `json:"server_uuid"`
	PublicIP       string    `json:"public_ip"`
	PrivateIP      string    `json:"private_ip"`
	PrivateGateway string    `json:"private_gateway"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// AdaptiveIPServers : Struct of AdaptiveIPServers
type AdaptiveIPServers struct {
	AdaptiveIP []Subnet `json:"adaptiveip"`
}

// AdaptiveIPServerNum : Struct of AdaptiveIPServerNum
type AdaptiveIPServerNum struct {
	Number int `json:"number"`
}