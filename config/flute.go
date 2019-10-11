package config

type flute struct {
	ServerAddress    string `goconf:"flute:flute_server_address"`     // ID : MySQL server login id
	ServerPort       int64  `goconf:"flute:flute_server_port"`        // Password : MySQL server login password
	RequestTimeoutMs int64  `goconf:"flute:flute_request_timeout_ms"` // Address : MySQL server address
}

// Flute : flute config structure
var Flute flute
