package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/harp/harp.conf"

type harpConfig struct {
	RsakeyConfig            *goconf.Section
	MysqlConfig             *goconf.Section
	GrpcConfig              *goconf.Section
	HornConfig              *goconf.Section
	CelloConfig             *goconf.Section
	FluteConfig             *goconf.Section
	ViolinConfig            *goconf.Section
	PiccoloConfig           *goconf.Section
	DHCPDConfig             *goconf.Section
	ARPINGConfig            *goconf.Section
	AdaptiveIPConfig        *goconf.Section
	AdaptiveIPNetworkConfig *goconf.Section
	VnStat                  *goconf.Section
	Timpani                 *goconf.Section
	HccwebConfig            *goconf.Section
}
