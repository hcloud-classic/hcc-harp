package types

import "time"

type Subnet struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	NetworkIP string    `json:"network_ip"`
	Netmask   string    `json:"netmask"`
	Os        string    `json:"os"`
	CreatedAt time.Time `json:"created_at"`
}

type Subnets struct {
	Subnets []Subnet `json:"subnet"`
}
