package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/harp/harp.conf"

type harpConfig struct {
	MysqlConfig             *goconf.Section
	GrpcConfig              *goconf.Section
	CelloConfig             *goconf.Section
	FluteConfig             *goconf.Section
	ViolinConfig            *goconf.Section
	PiccoloConfig           *goconf.Section
	DHCPDConfig             *goconf.Section
	ARPINGConfig            *goconf.Section
	AdaptiveIPConfig        *goconf.Section
	AdaptiveIPNetworkConfig *goconf.Section
}
