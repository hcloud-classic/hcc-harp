package config

type cello struct {
	ServerAddress string `goconf:"cello:cello_server_address"` // ServerAddress : IP address of server which installed cello module
}

// Cello : cello config structure
var Cello cello
