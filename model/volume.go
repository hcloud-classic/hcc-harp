package model

import "time"

// DefaultPXEdir : Default PXE directory
var DefaultPXEdir = "/HCC"

// OSDiskSize : Disk size for OS use
var OSDiskSize = 20

// Volume : Structure of Volume
type Volume struct {
	UUID       string    `json:"uuid"`
	Size       int       `json:"size"`
	Filesystem string    `json:"filesystem"`
	ServerUUID string    `json:"server_uuid"`
	UseType    string    `json:"use_type"`
	UserUUID   string    `json:"user_uuid"`
	CreatedAt  time.Time `json:"created_at"`
}

// Volumes : Structure of Volumes
type Volumes struct {
	Volumes []Volume `json:"volume"`
}

// VolumeNum : Structure of VolumeNum
type VolumeNum struct {
	Number int `json:"number"`
}
