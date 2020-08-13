package config

import (
	"github.com/Terry-Mao/goconf"
	"hcc/harp/lib/logger"
)

var conf = goconf.New()
var config = harpConfig{}
var err error

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

	Mysql.Password, err = config.MysqlConfig.String("password")
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

	DHCPD.IfconfigScriptFileLocation, err = config.DHCPDConfig.String("dhcpd_ifconfig_file_location")
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

	AdaptiveIP.ExternalIfaceName, err = config.AdaptiveIPConfig.String("adaptiveip_external_iface_name")
	if err != nil {
		logger.Logger.Panic(err)
	}

	AdaptiveIP.InternalIfaceName, err = config.AdaptiveIPConfig.String("adaptiveip_internal_iface_name")
	if err != nil {
		logger.Logger.Panic(err)
	}

	AdaptiveIP.PFBaseConfigFileLocation, err = config.AdaptiveIPConfig.String("adaptiveip_pf_base_config_file_location")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.PFRulesFileLocation, err = config.AdaptiveIPConfig.String("adaptiveip_pf_rules_file_location")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.PFBinatConfigFileLocation, err = config.AdaptiveIPConfig.String("adaptiveip_pf_binat_config_file_location")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.PFnatConfigFileLocation, err = config.AdaptiveIPConfig.String("adaptiveip_pf_nat_config_file_location")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.IfconfigScriptFileLocation, err = config.AdaptiveIPConfig.String("adaptiveip_ifconfig_script_file_location")
	if err != nil {
		logger.Logger.Panicln(err)
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

	AdaptiveIP.DefaultStartIPAddr, err = config.AdaptiveIPConfig.String("adaptiveip_default_start_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.DefaultEndIPAddr, err = config.AdaptiveIPConfig.String("adaptiveip_default_end_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.ArpingRetryCount, err = config.AdaptiveIPConfig.Int("adaptiveip_arping_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIP.ArpingRoutineMaxNum, err = config.AdaptiveIPConfig.Int("adaptiveip_arping_routine_max_num")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

// Parser : Parse config file
func Parser() {
	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}

	parseMysql()
	parseGrpc()
	parseFlute()
	parseViolin()
	parseDHCPD()
	parseAdaptiveIP()
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

	AdaptiveIPNetwork.StartIPAddr, err = config.AdaptiveIPNetworkConfig.String("adaptiveip_start_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	AdaptiveIPNetwork.EndIPAddr, err = config.AdaptiveIPNetworkConfig.String("adaptiveip_end_ip")
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
