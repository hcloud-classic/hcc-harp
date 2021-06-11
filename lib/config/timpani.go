package config

type timpani struct {
	TimpaniTargetIfaceName string `goconf:"timpani_target_iface_name"` // TimpaniTargetIfaceName : Target interface name for Timpani
	TimpaniExternalPort    int64  `goconf:"timpani_external_port"`     // TimpaniExternalPort : External port number for Timpani
	TimpaniInternalPort    int64  `goconf:"timpani_internal_port"`     // TimpaniInternalPort : Internal port number for Timpani
	TimpaniAddress         string `goconf:"timpani_address"`           // TimpaniAddress : Address of Timpani
}

// Timpani : timpani config structure
var Timpani timpani
