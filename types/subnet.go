package types

import "time"

// Subnet : Struct of Subnet
type Subnet struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	NetworkIP string    `json:"network_ip"`
	Netmask   string    `json:"netmask"`
	Os        string    `json:"os"`
	CreatedAt time.Time `json:"created_at"`
}

// Subnets : Array struct of subnets
type Subnets struct {
	Subnets []Subnet `json:"subnet"`
}
