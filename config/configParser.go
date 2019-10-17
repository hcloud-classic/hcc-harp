package config

import (
	"github.com/Terry-Mao/goconf"
	"hcc/harp/logger"
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

func parseHTTP() {
	config.HTTPConfig = conf.Get("http")
	if config.HTTPConfig == nil {
		logger.Logger.Panicln("no http section")
	}

	HTTP = http{}
	HTTP.Port, err = config.HTTPConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseRabbitMQ() {
	config.RabbitMQConfig = conf.Get("rabbitmq")
	if config.RabbitMQConfig == nil {
		logger.Logger.Panicln("no rabbitmq section")
	}

	RabbitMQ = rabbitmq{}
	RabbitMQ.ID, err = config.RabbitMQConfig.String("rabbitmq_id")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	RabbitMQ.Password, err = config.RabbitMQConfig.String("rabbitmq_password")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	RabbitMQ.Address, err = config.RabbitMQConfig.String("rabbitmq_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	RabbitMQ.Port, err = config.RabbitMQConfig.Int("rabbitmq_port")
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

// Parser : Parse config file
func Parser() {
	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}

	parseMysql()
	parseHTTP()
	parseRabbitMQ()
	parseFlute()
	parseViolin()
	parseDHCPD()
}
