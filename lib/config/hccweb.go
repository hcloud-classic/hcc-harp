package config

type hccweb struct {
	Port int64 `goconf:"hccweb:port"` // Port : Port number of hccweb inside the container
}

// Hccweb : Hccweb config structure
var Hccweb hccweb
