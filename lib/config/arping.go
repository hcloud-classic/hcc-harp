package config

type arping struct {
	IfaceName string `goconf:"arping:arping_iface_name"` // IfaceName : Interface name for use arping command
}

// ARPING : arping config structure
var ARPING arping
