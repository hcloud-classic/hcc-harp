package dao

import (
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/model"
	"time"
)

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

func checkReadSubnetListPageRow(args map[string]interface{}) bool {
	_, rowOk := args["row"].(int)
	_, pageOk := args["page"].(int)

	return !rowOk || !pageOk
}
