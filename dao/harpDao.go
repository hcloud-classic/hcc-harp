package dao

import (
	"errors"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/uuidgen"
	"hcc/harp/model"
	"net"
)

func CreateSubnet(args map[string]interface{}) (interface{}, error) {
	networkIP, networkIPOk := args["network_ip"].(string)
	netmask, netmaskOk := args["netmask"].(string)
	gateway, gatewayOk := args["gateway"].(string)
	nextServer, nextServerOk := args["next_server"].(string)
	name, nameOk := args["name"].(string)

	if !networkIPOk {
		return nil, errors.New("need network_ip argument")
	}
	if !netmaskOk {
		return nil, errors.New("need netmask argument")
	}
	if !gatewayOk {
		return nil, errors.New("need gateway argument")
	}
	if !nextServerOk {
		return nil, errors.New("need next_server argument")
	}
	if !nameOk {
		return nil, errors.New("need name argument")
	}

	netIPnetworkIP := iputil.CheckValidIP(networkIP)
	if netIPnetworkIP == nil {
		return nil, errors.New("wrong network IP")
	}

	mask, err := iputil.CheckNetmask(netmask)
	if err != nil {
		return nil, err
	}

	ipNet := net.IPNet{
		IP:   netIPnetworkIP,
		Mask: mask,
	}

	err = iputil.CheckGateway(ipNet, gateway)
	if err != nil {
		return nil, err
	}

	netIPnextServer := net.ParseIP(nextServer)
	if netIPnextServer == nil {
		return nil, errors.New("wrong next server IP")
	}

	nameServer, nameServerOk := args["name_server"].(string)
	if !nameServerOk {
		nameServer = ""
	}
	if len(nameServer) != 0 {
		netIPnameServer := net.ParseIP(nameServer)
		if netIPnameServer == nil {
			return nil, errors.New("wrong name server IP")
		}
	}

	domainName, domainNameOk := args["domain_name"].(string)
	if !domainNameOk {
		domainName = ""
	}

	uuid, err := uuidgen.UUIDgen()
	if err != nil {
		logger.Logger.Println("Failed to generate uuid!")
		return nil, err
	}

	subnet := model.Subnet{
		UUID:       uuid,
		NetworkIP:  networkIP,
		Netmask:    netmask,
		Gateway:    gateway,
		NextServer: nextServer,
		Name:       name,
		NameServer: nameServer,
		DomainName: domainName,
	}

	sql := "insert into subnet(network_ip, netmask, gateway, next_server, name, name_server, domain_name, created_at) values (?, ?, ?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err2 := stmt.Exec(subnet.NetworkIP, subnet.Netmask, subnet.Gateway, subnet.NextServer, subnet.Name, subnet.NameServer, subnet.DomainName)
	if err2 != nil {
		logger.Logger.Println(err2)
		return nil, err2
	}
	logger.Logger.Println(result.LastInsertId())

	return subnet, nil
}