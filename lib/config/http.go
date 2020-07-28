package config

type http struct {
	Port                 int64 `goconf:"http:port"`                    // Port : Port number for listening graphql request via http server
	RequestRetryCount    int64 `goconf:"http:request_retry_count"`     // RequestRetryCount : Retry count for http request
	RequestRetryDelaySec int64 `goconf:"http:request_retry_delay_sec"` // RequestRetryDelaySec : Retry delay for http request (seconds)
}

// HTTP : http config structure
var HTTP http
