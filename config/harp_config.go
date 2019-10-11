package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/harp/harp.conf"

type fluteConfig struct {
	MysqlConfig *goconf.Section
	HTTPConfig  *goconf.Section
	FluteConfig *goconf.Section
	DHCPDConfig *goconf.Section
}

/*-----------------------------------
         Config File Example

##### CONFIG START #####
[mysql]
id user
password pass
address 111.111.111.111
port 9999
database db_name

[http]
port 8888

[flute]
flute_server_address 222.222.222.222
flute_server_port 3333
flute_request_timeout_ms 5000

[dhcpd]
dhcpd_min_lease_time 1200
dhcpd_default_lease_time 1800
dhcpd_max_lease_time 3600
-----------------------------------*/
