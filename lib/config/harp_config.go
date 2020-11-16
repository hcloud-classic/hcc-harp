package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/harp/harp.conf"

type harpConfig struct {
	MysqlConfig             *goconf.Section
	HTTPConfig              *goconf.Section
	FluteConfig             *goconf.Section
	ViolinConfig            *goconf.Section
	DHCPDConfig             *goconf.Section
	ARPINGConfig            *goconf.Section
	AdaptiveIPConfig        *goconf.Section
	AdaptiveIPNetworkConfig *goconf.Section
}
