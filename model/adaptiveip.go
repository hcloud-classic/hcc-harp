package model

// AdaptiveIP : Struct of AdaptiveIP
type AdaptiveIP struct {
	ExtIfaceIPAddress string `json:"ext_iface_ip_address"`
	Netmask           string `json:"netmask"`
	GatewayAddress    string `json:"gateway"`
	StartIPAddress    string `json:"start_ip_address"`
	EndIPAddress      string `json:"end_ip_address"`
}

// AdaptiveIPAvailableIPList : Struct of AdaptiveIPAvailableIPList
type AdaptiveIPAvailableIPList struct {
	AvailableIPList []string `json:"available_ip_list"`
}
