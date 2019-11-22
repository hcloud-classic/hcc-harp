package model

// AdaptiveIPServer - ish
type AdaptiveIPServer struct {
	AdaptiveIPUUID string `json:"adaptiveip_uuid"`
	ServerUUID     string `json:"server_uuid"`
	PublicIP       string `json:"public_ip"`
	PrivateIP      string `json:"private_ip"`
	PrivateGateway string `json:"private_gateway"`
	Status string `json:"status"`
}

// AdaptiveIPServers - ish
type AdaptiveIPServers struct {
	AdaptiveIP []Subnet `json:"adaptiveip"`
}

// AdaptiveIPServerNum - ish
type AdaptiveIPServerNum struct {
	Number int `json:"number"`
}
