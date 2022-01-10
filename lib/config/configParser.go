package config

import (
	"github.com/Terry-Mao/goconf"
	"hcc/harp/lib/logger"
	"net"
)

var conf = goconf.New()
var config = harpConfig{}
var err error

func parseRsakey() {
	config.RsakeyConfig = conf.Get("rsakey")
	if config.RsakeyConfig == nil {
		logger.Logger.Panicln("no rsakey section")
	}

	Rsakey.PrivateKeyFile, err = config.RsakeyConfig.String("private_key_file")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseMysql() {
	config.MysqlConfig = conf.Get("mysql")
	if config.MysqlConfig == nil {
		logger.Logger.Panicln("no mysql section")
	}

	Mysql = mysql{}
	Mysql.ID, err = config.MysqlConfig.String("id")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Address, err = config.MysqlConfig.String("address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Port, err = config.MysqlConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Database, err = config.MysqlConfig.String("database")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.ConnectionRetryCount, err = config.MysqlConfig.Int("connection_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.ConnectionRetryIntervalMs, err = config.MysqlConfig.Int("connection_retry_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseGrpc() {
	config.GrpcConfig = conf.Get("grpc")
	if config.GrpcConfig == nil {
		logger.Logger.Panicln("no grpc section")
	}

	Grpc.Port, err = config.GrpcConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseHorn() {
	config.HornConfig = conf.Get("horn")
	if config.HornConfig == nil {
		logger.Logger.Panicln("no horn section")
	}

	Horn = horn{}
	Horn.ServerAddress, err = config.HornConfig.String("horn_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Horn.ServerPort, err = config.HornConfig.Int("horn_server_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Horn.ConnectionTimeOutMs, err = config.HornConfig.Int("horn_connection_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Horn.ConnectionRetryCount, err = config.HornConfig.Int("horn_connection_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Horn.RequestTimeoutMs, err = config.HornConfig.Int("horn_request_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseCello() {
	config.CelloConfig = conf.Get("cello")
	if config.CelloConfig == nil {
		logger.Logger.Panicln("no cello section")
	}

	Cello = cello{}
	Cello.ServerAddress, err = config.CelloConfig.String("cello_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	netIP := net.ParseIP(Cello.ServerAddress).To4()
	if netIP == nil {
		logger.Logger.Panicln("Cello server address is configured incorrectly")
	}
}

func parseFlute() {
	config.FluteConfig = conf.Get("flute")
	if config.FluteConfig == nil {
		logger.Logger.Panicln("no flute section")
	}

	Flute = flute{}
	Flute.ServerAddress, err = config.FluteConfig.String("flute_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Flute.ServerPort, err = config.FluteConfig.Int("flute_server_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Flute.ConnectionTimeOutMs, err = config.FluteConfig.Int("flute_connection_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Flute.ConnectionRetryCount, err = config.FluteConfig.Int("flute_connection_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Flute.RequestTimeoutMs, err = config.FluteConfig.Int("flute_request_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseViolin() {
	config.ViolinConfig = conf.Get("violin")
	if config.ViolinConfig == nil {
		logger.Logger.Panicln("no violin section")
	}

	Violin = violin{}
	Violin.ServerAddress, err = config.ViolinConfig.String("violin_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Violin.ServerPort, err = config.ViolinConfig.Int("violin_server_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Violin.ConnectionTimeOutMs, err = config.ViolinConfig.Int("violin_connection_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Violin.ConnectionRetryCount, err = config.ViolinConfig.Int("violin_connection_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Violin.RequestTimeoutMs, err = config.ViolinConfig.Int("violin_request_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parsePiccolo() {
	config.PiccoloConfig = conf.Get("piccolo")
	if config.PiccoloConfig == nil {
		logger.Logger.Panicln("no piccolo section")
	}

	Piccolo = piccolo{}
	Piccolo.ServerAddress, err = config.PiccoloConfig.String("piccolo_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Piccolo.ServerPort, err = config.PiccoloConfig.Int("piccolo_server_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Piccolo.ConnectionTimeOutMs, err = config.PiccoloConfig.Int("piccolo_connection_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Piccolo.ConnectionRetryCount, err = config.PiccoloConfig.Int("piccolo_connection_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Piccolo.RequestTimeoutMs, err = config.PiccoloConfig.Int("piccolo_request_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseDHCPD() {
	config.DHCPDConfig = conf.Get("dhcpd")
	if config.DHCPDConfig == nil {
		logger.Logger.Panicln("no dhcpd section")
	}

	DHCPD = dhcpd{}

	DHCPD.LocalDHCPDServiceName, err = config.DHCPDConfig.String("dhcpd_local_dhcpd_service_name")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	DHCPD.LocalConfigFileLocation, err = config.DHCPDConfig.String("dhcpd_local_config_file_location")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	DHCPD.ConfigFileLocation, err = config.DHCPDConfig.String("dhcpd_config_file_location")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	DHCPD.MinLeaseTime, err = config.DHCPDConfig.Int("dhcpd_min_lease_time")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	DHCPD.DefaultLeaseTime, err = config.DHCPDConfig.Int("dhcpd_default_lease_time")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	DHCPD.MaxLeaseTime, err = config.DHCPDConfig.Int("dhcpd_max_lease_time")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseAdaptiveIP() {
	config.AdaptiveIPConfig = conf.Get("adaptiveip")
	if config.AdaptiveIPConfig == nil {
		logger.Logger.Panicln("no adaptiveip section")
	}

	AdaptiveIP = adaptiveIP{}

	AdaptiveIP.CustomScriptsLocation, err = config.AdaptiveIPConfig.String("adaptiveip_custom_scripts_location")
	if err != nil {
		logger.Logger.Panic(err)
	}

	AdaptiveIP.ExternalIfaceName, err = config.AdaptiveIPConfig.String("adaptiveip_external_iface_name")
	if err != nil {
		logger.Logger.Panic(err)
	}

	AdaptiveIP.InternalIfaceName, err = config.AdaptiveIPConfig.String("adaptiveip_internal_iface_name")
	if err != nil {
		logger.Logger.Panic(err)
	}

	AdaptiveIP.NetworkConfigFile, err = config.AdaptiveIPConfig.String("adaptiveip_network_config_file")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.DefaultExtIfaceIPAddr, err = config.AdaptiveIPConfig.String("adaptiveip_default_ext_iface_ip_addr")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.DefaultNetmask, err = config.AdaptiveIPConfig.String("adaptiveip_default_netmask")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.DefaultGatewayAddr, err = config.AdaptiveIPConfig.String("adaptiveip_default_gateway_addr")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.DefaultInternalStartIPAddr, err = config.AdaptiveIPConfig.String("adaptiveip_default_internal_start_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.DefaultInternalEndIPAddr, err = config.AdaptiveIPConfig.String("adaptiveip_default_internal_end_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.DefaultExternalStartIPAddr, err = config.AdaptiveIPConfig.String("adaptiveip_default_external_start_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.DefaultExternalEndIPAddr, err = config.AdaptiveIPConfig.String("adaptiveip_default_external_end_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.ArpingRoutineMaxNum, err = config.AdaptiveIPConfig.Int("adaptiveip_arping_routine_max_num")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseVnStat() {
	config.VnStat = conf.Get("vnstat")
	if config.VnStat == nil {
		logger.Logger.Panicln("no vnstat section")
	}

	VnStat = vnStat{}

	VnStat.Debug, err = config.VnStat.String("vnstat_debug")
	if err != nil {
		logger.Logger.Panic(err)
	}

	VnStat.DatabaseUpdateIntervalSec, err = config.VnStat.Int("vnstat_database_update_interval_sec")
	if err != nil {
		logger.Logger.Panic(err)
	}
}

func parseTimpani() {
	config.Timpani = conf.Get("timpani")
	if config.Timpani == nil {
		logger.Logger.Panicln("no timpani section")
	}

	Timpani = timpani{}

	Timpani.TimpaniTargetIfaceName, err = config.Timpani.String("timpani_target_iface_name")
	if err != nil {
		logger.Logger.Panic(err)
	}

	Timpani.TimpaniExternalPort, err = config.Timpani.Int("timpani_external_port")
	if err != nil {
		logger.Logger.Panic(err)
	}
	if Timpani.TimpaniExternalPort < 0 || Timpani.TimpaniExternalPort > 65535 {
		logger.Logger.Panic("Port number is out of range (timpani_external_port)")
	}

	Timpani.TimpaniInternalPort, err = config.Timpani.Int("timpani_internal_port")
	if err != nil {
		logger.Logger.Panic(err)
	}
	if Timpani.TimpaniExternalPort < 0 || Timpani.TimpaniExternalPort > 65535 {
		logger.Logger.Panic("Port number is out of range (timpani_internal_port)")
	}

	Timpani.TimpaniAddress, err = config.Timpani.String("timpani_address")
	if err != nil {
		logger.Logger.Panic(err)
	}
}

// Init : Parse config file and initialize config structure
func Init() {
	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}

	parseRsakey()
	parseMysql()
	parseGrpc()
	parseHorn()
	parseCello()
	parseFlute()
	parseViolin()
	parsePiccolo()
	parseDHCPD()
	parseAdaptiveIP()
	parseVnStat()
	parseTimpani()
}

func parseAdaptiveIPNetwork(adaptiveipNetworkConf *goconf.Config) {
	config.AdaptiveIPNetworkConfig = adaptiveipNetworkConf.Get("adaptiveip_network")
	if config.AdaptiveIPNetworkConfig == nil {
		logger.Logger.Panicln("no adaptiveip_network section")
	}

	AdaptiveIPNetwork = adaptiveIPNetwork{}

	AdaptiveIPNetwork.ExtIfaceIPAddr, err = config.AdaptiveIPNetworkConfig.String("adaptiveip_ext_iface_ip_addr")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIPNetwork.Netmask, err = config.AdaptiveIPNetworkConfig.String("adaptiveip_netmask")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIPNetwork.GatewayAddr, err = config.AdaptiveIPNetworkConfig.String("adaptiveip_gateway_addr")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIPNetwork.InternalStartIPAddr, err = config.AdaptiveIPNetworkConfig.String("adaptiveip_internal_start_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIPNetwork.InternalEndIPAddr, err = config.AdaptiveIPNetworkConfig.String("adaptiveip_internal_end_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIPNetwork.ExternalStartIPAddr, err = config.AdaptiveIPNetworkConfig.String("adaptiveip_external_start_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIPNetwork.ExternalEndIPAddr, err = config.AdaptiveIPNetworkConfig.String("adaptiveip_external_end_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

// AdaptiveIPNetworkConfigParser : Parse Adaptive IP network config
func AdaptiveIPNetworkConfigParser() error {
	adaptiveipNetworkConf := goconf.New()

	err := adaptiveipNetworkConf.Parse(AdaptiveIP.NetworkConfigFile)
	if err != nil {
		logger.Logger.Println(err)
		return err
	}

	parseAdaptiveIPNetwork(adaptiveipNetworkConf)
	return nil
}
