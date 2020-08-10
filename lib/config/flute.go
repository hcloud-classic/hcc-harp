package config

type flute struct {
	ServerAddress    string `goconf:"flute:flute_server_address"`     // ServerAddress : IP address of server which installed flute module
	ServerPort       int64  `goconf:"flute:flute_server_port"`        // ServerPort : gRPC listening port number of flute module
	RequestTimeoutMs int64  `goconf:"flute:flute_request_timeout_ms"` // RequestTimeoutMs : Timeout for gRPC request to flute module
}

// Flute : flute config structure
var Flute flute
