package dao

import (
	"errors"
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

func ReadSubnetList(args map[string]interface{}) (interface{}, error) {
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

	row, _ := args["row"].(int)
	page, _ := args["page"].(int)
	if checkReadSubnetListPageRow(args) {
		return nil, errors.New("need row and page arguments")
	}

	sql := "select * from subnet where 1=1"

	if networkIPOk {
		sql += " and network_ip = '" + networkIP + "'"
	}
	if netmaskOk {
		sql += " and netmask = '" + netmask + "'"
	}
	if gatewayOk {
		sql += " and gateway = '" + gateway + "'"
	}
	if nextServerOk {
		sql += " and next_server = '" + nextServer + "'"
	}
	if nameServerOk {
		sql += " and name_server = '" + nameServer + "'"
	}
	if domainNameOk {
		sql += " and domain_name = '" + domainName + "'"
	}
	if serverUUIDOk {
		sql += " and server_uuid = '" + serverUUID + "'"
	}
	if leaderNodeUUIDOk {
		sql += " and leader_node_uuid = '" + leaderNodeUUID + "'"
	}
	if osOk {
		sql += " and os = '" + os + "'"
	}
	if subnetNameOk {
		sql += " and subnet_name = '" + subnetName + "'"
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

func (args map[string]interface{}) (interface{}, error) {
	out, err := gouuid.NewV4()
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	uuid := out.String()

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



func checkUpdateSubnetArgs(reqSubnet *pb.Subnet) bool {
	networkIPOk := len(reqSubnet.GetNetworkIP()) != 0
	netmaskOk := len(reqSubnet.GetNetmask()) != 0
	gatewayOk := len(reqSubnet.GetGateway()) != 0
	nextServerOk := len(reqSubnet.GetNextServer()) != 0
	nameServerOk := len(reqSubnet.GetNameServer()) != 0
	domainNameOk := len(reqSubnet.GetDomainName()) != 0
	serverUUIDOk := len(reqSubnet.GetServerUUID()) != 0
	leaderNodeUUIDOk := len(reqSubnet.GetLeaderNodeUUID()) != 0
	osOk := len(reqSubnet.GetOS()) != 0
	subnetNameOk := len(reqSubnet.GetSubnetName()) != 0

	return !networkIPOk && !netmaskOk && !gatewayOk && !nextServerOk && !nameServerOk && !domainNameOk && !serverUUIDOk && !leaderNodeUUIDOk && !osOk && !subnetNameOk
}

func UpdateSubnet(in *pb.ReqUpdateSubnet) (*pb.Subnet, error) {
	if in.Subnet == nil {
		return nil, errors.New("subnet is nil")
	}
	reqSubnet := in.Subnet

	requestedUUID := reqSubnet.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return nil, errors.New("need a uuid argument")
	}

	if checkUpdateSubnetArgs(reqSubnet) {
		return nil, errors.New("need some arguments")
	}

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

	networkIP = in.GetSubnet().NetworkIP
	networkIPOk := len(networkIP) != 0
	netmask = in.GetSubnet().Netmask
	netmaskOk := len(netmask) != 0
	gateway = in.GetSubnet().Gateway
	gatewayOk := len(gateway) != 0
	nextServer = in.GetSubnet().NextServer
	nextServerOk := len(nextServer) != 0
	nameServer = in.GetSubnet().NameServer
	nameServerOk := len(nameServer) != 0
	domainName = in.GetSubnet().DomainName
	domainNameOk := len(domainName) != 0
	serverUUID = in.GetSubnet().ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	leaderNodeUUID = in.GetSubnet().LeaderNodeUUID
	leaderNodeUUIDOk := len(leaderNodeUUID) != 0
	os = in.GetSubnet().OS
	osOk := len(os) != 0
	subnetName = in.GetSubnet().SubnetName
	subnetNameOk := len(subnetName) != 0

	subnet := new(pb.Subnet)
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

	oldSubnet, err := ReadSubnet(subnet.GetUUID())
	if err != nil {
		return nil, err
	}

	if !networkIPOk {
		subnet.NetworkIP = oldSubnet.NetworkIP
	}
	if !netmaskOk {
		subnet.Netmask = oldSubnet.Netmask
	}
	if !gatewayOk {
		subnet.Gateway = oldSubnet.Gateway
	}

	err = checkSubnet(subnet.NetworkIP, subnet.Netmask, subnet.Gateway, true, oldSubnet)
	if err != nil {
		return nil, err
	}

	if serverUUIDOk {
		err = checkServerUUID(subnet.ServerUUID)
		if err != nil {
			return nil, err
		}
	}

	sql := "update subnet set"
	var updateSet = ""
	if networkIPOk {
		updateSet += " network_ip = '" + subnet.NetworkIP + "', "
	}
	if netmaskOk {
		updateSet += " netmask = '" + subnet.Netmask + "', "
	}
	if gatewayOk {
		updateSet += " gateway = '" + subnet.Gateway + "', "
	}
	if nextServerOk {
		updateSet += " next_server = '" + subnet.NextServer + "', "
	}
	if nameServerOk {
		updateSet += " name_server = '" + subnet.NameServer + "', "
	}
	if domainNameOk {
		updateSet += " domain_name = '" + subnet.DomainName + "', "
	}
	if serverUUIDOk {
		updateSet += " server_uuid = '" + subnet.ServerUUID + "', "
	}
	if leaderNodeUUIDOk {
		updateSet += " leader_node_uuid = '" + subnet.LeaderNodeUUID + "', "
	}
	if osOk {
		updateSet += " os = '" + subnet.OS + "', "
	}
	if subnetNameOk {
		updateSet += " subnet_name = '" + subnet.SubnetName + "', "
	}
	sql += updateSet[0:len(updateSet)-2] + " where uuid = ?"

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

func DeleteSubnet(in *pb.ReqDeleteSubnet) (string, error) {
	var err error

	requestedUUID := in.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return "", errors.New("need a uuid argument")
	}

	subnet, err := ReadSubnet(requestedUUID)
	if err != nil {
		return "", err
	}

	if len(subnet.ServerUUID) == 0 {
		msg := "subnet is used by the server (UUID:" + subnet.ServerUUID + ")"
		logger.Logger.Println(msg)
		return "", errors.New(msg)
	}

	sql := "delete from subnet where uuid = ?"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return "", err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err2 := stmt.Exec(requestedUUID)
	if err2 != nil {
		logger.Logger.Println(err2)
		return "", err
	}
	logger.Logger.Println(result.RowsAffected())

	return requestedUUID, nil
}
