package config

type vnStat struct {
	Debug                     string `goconf:"vnstat:vnstat_debug"`                        // Debug : Enable debug logs for VnStat
	DatabaseUpdateIntervalSec int64  `goconf:"vnstat:vnstat_database_update_interval_sec"` // DatabaseUpdateIntervalSec : Update interval of vnStat database (Seconds)
}

// VnStat : vnStat config structure
var VnStat vnStat
