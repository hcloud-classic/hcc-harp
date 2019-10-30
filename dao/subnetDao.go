package dao

import (
	"errors"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/uuidgen"
	"hcc/harp/model"
	"time"
)

// ReadSubnet - cgs
func ReadSubnet(args map[string]interface{}) (interface{}, error) {
	var subnet model.Subnet

	uuid := args["uuid"].(string)
	var networkIP string
	var netmask string
	var gateway string
	var nextServer string
	var nameServer string
	var domainName string
	var serverUUID string
	var leaderNodeUUID string
	var _os string
	var subnetName string
	var createdAt time.Time

	sql := "select network_ip, netmask, gateway, next_server, name_server, domain_name, server_uuid, leader_node_uuid, os, subnet_name, created_at from subnet where uuid = ?"
	err := mysql.Db.QueryRow(sql, uuid).Scan(
		&networkIP,
		&netmask,
		&gateway,
		&nextServer,
		&nameServer,
		&domainName,
		&serverUUID,
		&leaderNodeUUID,
		&_os,
		&subnetName,
		&createdAt)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	subnet.UUID = uuid
	subnet.NetworkIP = networkIP
	subnet.Netmask = netmask
	subnet.Gateway = gateway
	subnet.NextServer = nextServer
	subnet.NameServer = nameServer
	subnet.DomainName = domainName
	subnet.ServerUUID = serverUUID
	subnet.LeaderNodeUUID = leaderNodeUUID
	subnet.OS = _os
	subnet.SubnetName = subnetName
	subnet.CreatedAt = createdAt

	return subnet, nil
}

// ReadSubnetList - cgs
func ReadSubnetList(args map[string]interface{}) (interface{}, error) {
	var err error
	var subnets []model.Subnet
	var uuid string
	var createdAt time.Time

	networkIP, networkIPOk := args["network_ip"].(string)
	netmask, netmaskOk := args["netmask"].(string)
	gateway, gatewayOk := args["gateway"].(string)
	nextServer, nextServerOk := args["next_server"].(string)
	nameServer, nameServerOk := args["name_server"].(string)
	domainName, domainNameOk := args["domain_name"].(string)
	serverUUID, serverUUIDOk := args["sever_uuid"].(string)
	leaderNodeUUID, leaderNodeUUIDOk := args["leader_node_uuid"].(string)
	os, osOk := args["os"].(string)
	subnetName, subnetNameOk := args["subnet_name"].(string)

	row, rowOk := args["row"].(int)
	page, pageOk := args["page"].(int)
	if !rowOk || !pageOk {
		return nil, err
	}

	sql := "select * from subnet where 1 = 1 and "

	if networkIPOk {
		sql += " network_ip = '" + networkIP + "'"
		if netmaskOk || gatewayOk || nextServerOk || nameServerOk || domainNameOk || serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
			sql += " and"
		}
	}
	if netmaskOk {
		sql += " netmask = '" + netmask + "'"
		if gatewayOk || nextServerOk || nameServerOk || domainNameOk || serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
			sql += " and"
		}
	}
	if gatewayOk {
		sql += " gateway = '" + gateway + "'"
		if nextServerOk || nameServerOk || domainNameOk || serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
			sql += " and"
		}
	}
	if nextServerOk {
		sql += " next_server = '" + nextServer + "'"
		if nameServerOk || domainNameOk || serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
			sql += " and"
		}
	}
	if nameServerOk {
		sql += " name_server = '" + nameServer + "'"
		if domainNameOk || serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
			sql += " and"
		}
	}
	if domainNameOk {
		sql += " domain_name = '" + domainName + "'"
		if serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
			sql += " and"
		}
	}
	if serverUUIDOk {
		sql += " server_uuid = '" + serverUUID + "'"
		if leaderNodeUUIDOk || osOk || subnetNameOk {
			sql += " and"
		}
	}
	if leaderNodeUUIDOk {
		sql += " leader_node_uuid = '" + leaderNodeUUID + "'"
		if osOk || subnetNameOk {
			sql += " and"
		}
	}
	if osOk {
		sql += " os = '" + os + "'"
		if subnetNameOk {
			sql += " and"
		}
	}
	if subnetNameOk {
		sql += " subnet_name = '" + subnetName + "'"
	}

	sql += " order by created_at desc limit ? offset ?"

	stmt, err := mysql.Db.Query(sql, row, row*(page-1))
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &networkIP, &netmask, &gateway, &nextServer, &nameServer, &domainName, &serverUUID, &leaderNodeUUID, &os, &subnetName, &createdAt)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		subnet := model.Subnet{UUID: uuid, NetworkIP: networkIP, Netmask: netmask, Gateway: gateway, NextServer: nextServer, NameServer: nameServer, DomainName: domainName, ServerUUID: serverUUID, LeaderNodeUUID: leaderNodeUUID, OS: os, SubnetName: subnetName, CreatedAt: createdAt}
		subnets = append(subnets, subnet)
	}
	return subnets, nil
}

// ReadSubnetAll - cgs
func ReadSubnetAll(args map[string]interface{}) (interface{}, error) {

	var err error
	var subnets []model.Subnet
	var uuid string
	var networkIP string
	var netmask string
	var gateway string
	var nextServer string
	var nameServer string
	var domainName string
	var serverUUID string
	var leaderNodeUUID string
	var os string
	var subnetName string
	var createdAt time.Time
	row, rowOk := args["row"].(int)
	page, pageOk := args["page"].(int)
	if !rowOk || !pageOk {
		return nil, err
	}

	sql := "select * from subnet order by created_at desc limit ? offset ?"

	stmt, err := mysql.Db.Query(sql, row, row*(page-1))
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &networkIP, &netmask, &gateway, &nextServer, &nameServer, &domainName, &serverUUID, &leaderNodeUUID, &os, &subnetName, &createdAt)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		subnet := model.Subnet{UUID: uuid, NetworkIP: networkIP, Netmask: netmask, Gateway: gateway, NextServer: nextServer, NameServer: nameServer, DomainName: domainName, ServerUUID: serverUUID, LeaderNodeUUID: leaderNodeUUID, OS: os, SubnetName: subnetName, CreatedAt: createdAt}
		subnets = append(subnets, subnet)
	}

	return subnets, nil
}

// ReadSubnetNum - cgs
func ReadSubnetNum() (model.SubnetNum, error) {
	var subnetNum model.SubnetNum
	var subnetNr int
	var err error

	sql := "select count(*) from subnet"
	err = mysql.Db.QueryRow(sql).Scan(&subnetNr)
	if err != nil {
		logger.Logger.Println(err)
		return subnetNum, err
	}
	subnetNum.Number = subnetNr

	return subnetNum, nil
}

// CreateSubnet - cgs
func CreateSubnet(args map[string]interface{}) (interface{}, error) {
	uuid, err := uuidgen.UUIDgen()
	if err != nil {
		logger.Logger.Println("Failed to generate uuid!")
		return nil, err
	}

	subnet := model.Subnet{
		UUID:           uuid,
		NetworkIP:      args["network_ip"].(string),
		Netmask:        args["netmask"].(string),
		Gateway:        args["gateway"].(string),
		NextServer:     args["next_server"].(string),
		NameServer:     args["name_server"].(string),
		DomainName:     args["domain_name"].(string),
		ServerUUID:     args["server_uuid"].(string),
		LeaderNodeUUID: args["leader_node_uuid"].(string),
		OS:             args["os"].(string),
		SubnetName:     args["subnet_name"].(string),
	}

	sql := "insert into subnet(uuid, network_ip, netmask, gateway, next_server, name_server, domain_name, server_uuid, leader_node_uuid, os, subnet_name, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err := stmt.Exec(subnet.UUID, subnet.NetworkIP, subnet.Netmask, subnet.Gateway, subnet.NextServer, subnet.NameServer, subnet.DomainName, subnet.ServerUUID, subnet.LeaderNodeUUID, subnet.OS, subnet.SubnetName)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	logger.Logger.Println(result.LastInsertId())

	return subnet, nil
}

// UpdateSubnet - cgs
func UpdateSubnet(args map[string]interface{}) (interface{}, error) {
	requestedUUID, requestedUUIDOk := args["uuid"].(string)
	networkIP, networkIPOk := args["network_ip"].(string)
	netmask, netmaskOk := args["netmask"].(string)
	gateway, gatewayOk := args["gateway"].(string)
	nextServer, nextServerOk := args["next_server"].(string)
	nameServer, nameServerOk := args["name_server"].(string)
	domainName, domainNameOk := args["domain_name"].(string)
	serverUUID, serverUUIDOk := args["server_uuid"].(string)
	leaderNodeUUID, leaderNodeUUIDOk := args["leader_node_uuid"].(string)
	os, osOk := args["os"].(string)
	subnetName, subnetNameOk := args["subnet_name"].(string)

	subnet := new(model.Subnet)
	subnet.UUID = requestedUUID
	subnet.NetworkIP = networkIP
	subnet.Netmask = netmask
	subnet.Gateway = gateway
	subnet.NextServer = nextServer
	subnet.NameServer = nameServer
	subnet.DomainName = domainName
	subnet.ServerUUID = serverUUID
	subnet.LeaderNodeUUID = leaderNodeUUID
	subnet.OS = os
	subnet.SubnetName = subnetName

	if requestedUUIDOk {
		if !networkIPOk && !netmaskOk && !gatewayOk && !nextServerOk && !nameServerOk && !domainNameOk && !serverUUIDOk && !leaderNodeUUIDOk && !osOk && !subnetNameOk {
			return nil, errors.New("need some arguments")
		}

		sql := "update subnet set"
		if networkIPOk {
			sql += " network_ip = '" + subnet.NetworkIP + "'"
			if netmaskOk || gatewayOk || nextServerOk || nameServerOk || domainNameOk || serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
				sql += ", "
			}
		}
		if netmaskOk {
			sql += " netmask = '" + subnet.Netmask + "'"
			if gatewayOk || nextServerOk || nameServerOk || domainNameOk || serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
				sql += ", "
			}
		}
		if gatewayOk {
			sql += " gateway = '" + subnet.Gateway + "'"
			if nextServerOk || nameServerOk || domainNameOk || serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
				sql += ", "
			}
		}
		if nextServerOk {
			sql += " next_server = '" + subnet.NextServer + "'"
			if nameServerOk || domainNameOk || serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
				sql += ", "
			}
		}
		if nameServerOk {
			sql += " name_server = '" + subnet.NameServer + "'"
			if domainNameOk || serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
				sql += ", "
			}
		}
		if domainNameOk {
			sql += " domain_name = '" + subnet.DomainName + "'"
			if serverUUIDOk || leaderNodeUUIDOk || osOk || subnetNameOk {
				sql += ", "
			}
		}
		if serverUUIDOk {
			sql += " server_uuid = '" + subnet.ServerUUID + "'"
			if leaderNodeUUIDOk || osOk || subnetNameOk {
				sql += ", "
			}
		}
		if leaderNodeUUIDOk {
			sql += " leader_node_uuid = '" + subnet.LeaderNodeUUID + "'"
			if osOk || subnetNameOk {
				sql += ", "
			}
		}
		if osOk {
			sql += " os = '" + subnet.OS + "'"
			if subnetNameOk {
				sql += ", "
			}
		}
		if subnetNameOk {
			sql += " subnet_name = '" + subnet.SubnetName
		}
		sql += " where uuid = ?"

		logger.Logger.Println("update_subnet sql : ", sql)

		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()

		result, err2 := stmt.Exec(subnet.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err
		}
		logger.Logger.Println(result.LastInsertId())
		return subnet, nil
	}

	return nil, errors.New("need uuid argument")
}

// DeleteSubnet - cgs
func DeleteSubnet(args map[string]interface{}) (interface{}, error) {
	var err error

	requestedUUID, ok := args["uuid"].(string)
	if ok {
		sql := "delete from subnet where uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()
		result, err2 := stmt.Exec(requestedUUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err
		}
		logger.Logger.Println(result.RowsAffected())

		return requestedUUID, nil
	}

	return requestedUUID, err
}
