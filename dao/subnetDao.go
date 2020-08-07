package dao

import (
	dbsql "database/sql"
	"errors"
	"github.com/golang/protobuf/ptypes"
	gouuid "github.com/nu7hatch/gouuid"
	pb "hcc/harp/action/grpc/rpcharp"
	"hcc/harp/data"
	"hcc/harp/driver"
	"hcc/harp/lib/config"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"time"
)

// ReadSubnet : Get infos of a subnet
func ReadSubnet(uuid string) (*pb.Subnet, error) {
	var subnet pb.Subnet

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

	subnet.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	return &subnet, nil
}

// ReadSubnetByServer : Get infos of a subnet by server UUID
func ReadSubnetByServer(serverUUID string) (*pb.Subnet, error) {
	var subnet pb.Subnet

	var uuid string
	var networkIP string
	var netmask string
	var gateway string
	var nextServer string
	var nameServer string
	var domainName string
	var leaderNodeUUID string
	var _os string
	var subnetName string
	var createdAt time.Time

	sql := "select uuid, network_ip, netmask, gateway, next_server, name_server, domain_name, leader_node_uuid, os, subnet_name, created_at from subnet where server_uuid = ?"
	err := mysql.Db.QueryRow(sql, serverUUID).Scan(
		&uuid,
		&networkIP,
		&netmask,
		&gateway,
		&nextServer,
		&nameServer,
		&domainName,
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

	subnet.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	return &subnet, nil
}

// ReadSubnetList : Get list of subnet with selected infos
func ReadSubnetList(in *pb.ReqGetSubnetList) (*pb.ResGetSubnetList, error) {
	var subnetList pb.ResGetSubnetList
	var subnets []pb.Subnet
	var psubnets []*pb.Subnet

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

	var isLimit bool
	row := in.GetRow()
	rowOk := row != 0
	page := in.GetPage()
	pageOk := page != 0
	if !rowOk && !pageOk {
		isLimit = false
	} else if rowOk && pageOk {
		isLimit = true
	} else {
		return nil, errors.New("please insert row and page arguments or leave arguments as empty state")
	}

	sql := "select * from subnet where 1=1"

	if in.Subnet != nil {
		reqSubnet := in.Subnet

		networkIP = reqSubnet.NetworkIP
		networkIPOk := len(networkIP) != 0
		netmask = reqSubnet.Netmask
		netmaskOk := len(netmask) != 0
		gateway = reqSubnet.Gateway
		gatewayOk := len(gateway) != 0
		nextServer = reqSubnet.NextServer
		nextServerOk := len(nextServer) != 0
		nameServer = reqSubnet.NameServer
		nameServerOk := len(nameServer) != 0
		domainName = reqSubnet.DomainName
		domainNameOk := len(domainName) != 0
		serverUUID = reqSubnet.ServerUUID
		serverUUIDOk := len(serverUUID) != 0
		leaderNodeUUID = reqSubnet.LeaderNodeUUID
		leaderNodeUUIDOk := len(leaderNodeUUID) != 0
		os = reqSubnet.OS
		osOk := len(os) != 0
		subnetName = reqSubnet.SubnetName
		subnetNameOk := len(subnetName) != 0

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
	}

	var stmt *dbsql.Rows
	var err error
	if isLimit {
		sql += " order by created_at desc limit ? offset ?"
		stmt, err = mysql.Db.Query(sql, row, row*(page-1))
	} else {
		sql += " order by created_at desc"
		stmt, err = mysql.Db.Query(sql)
	}

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

		_createdAt, err := ptypes.TimestampProto(createdAt)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		subnets = append(subnets, pb.Subnet{
			UUID:           uuid,
			NetworkIP:      networkIP,
			Netmask:        netmask,
			Gateway:        gateway,
			NextServer:     nextServer,
			NameServer:     nameServer,
			DomainName:     domainName,
			ServerUUID:     serverUUID,
			LeaderNodeUUID: leaderNodeUUID,
			OS:             os,
			SubnetName:     subnetName,
			CreatedAt:      _createdAt})
	}

	for i := range subnets {
		psubnets = append(psubnets, &subnets[i])
	}

	subnetList.Subnet = psubnets

	return &subnetList, nil
}

// ReadSubnetNum : Get the number of subnets
func ReadSubnetNum() (*pb.ResGetSubnetNum, error) {
	var resSubnetNum pb.ResGetSubnetNum
	var subnetNr int64

	sql := "select count(*) from subnet"
	err := mysql.Db.QueryRow(sql).Scan(&subnetNr)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	resSubnetNum.Num = subnetNr

	return &resSubnetNum, nil
}

func checkSubnet(networkIP string, netmask string, gateway string, skipMine bool, oldSubnet interface{}) error {
	isPrivate, err := iputil.CheckPrivateSubnet(networkIP, netmask)
	if !isPrivate {
		return errors.New("given network IP address is not in private network")
	}
	if err != nil {
		return err
	}

	isConflict, err := iputil.CheckSubnetConflict(networkIP, netmask, skipMine, oldSubnet)
	if isConflict {
		return errors.New("given subnet is conflicted with one of subnet that stored in the database")
	}
	if err != nil {
		return err
	}

	netNetwork, err := iputil.CheckNetwork(networkIP, netmask)
	if err != nil {
		return err
	}

	err = iputil.CheckIPisInSubnet(*netNetwork, gateway)
	if err != nil {
		return err
	}

	return nil
}

func checkServerUUID(serverUUID string) error {
	allServerData, err := driver.AllServerUUID()
	if err != nil {
		return err
	}

	allServer := allServerData.(data.AllServerData).Data.AllServer
	for _, server := range allServer {
		_serverUUID := server.UUID
		if _serverUUID == serverUUID {
			return nil
		}
	}

	return errors.New("given server UUID is not in the database")
}

func checkCreateSubnetArgs(reqSubnet *pb.Subnet) bool {
	networkIPOk := len(reqSubnet.GetNetworkIP()) != 0
	netmaskOk := len(reqSubnet.GetNetmask()) != 0
	gatewayOk := len(reqSubnet.GetGateway()) != 0
	nextServerOk := len(reqSubnet.GetNextServer()) != 0
	nameServerOk := len(reqSubnet.GetNameServer()) != 0
	domainNameOk := len(reqSubnet.GetDomainName()) != 0
	osOk := len(reqSubnet.GetOS()) != 0
	subnetNameOk := len(reqSubnet.GetSubnetName()) != 0

	return !(networkIPOk && netmaskOk && gatewayOk && nextServerOk && nameServerOk && domainNameOk && osOk && subnetNameOk)
}

// CreateSubnet : Create a subnet
func CreateSubnet(in *pb.ReqCreateSubnet) (*pb.Subnet, error) {
	reqSubnet := in.GetSubnet()
	if reqSubnet == nil {
		return nil, errors.New("subnet is nil")
	}

	out, err := gouuid.NewV4()
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	uuid := out.String()

	if checkCreateSubnetArgs(reqSubnet) {
		return nil, errors.New("some of arguments are missing")
	}

	subnet := pb.Subnet{
		UUID:           uuid,
		NetworkIP:      reqSubnet.GetNetworkIP(),
		Netmask:        reqSubnet.GetNetmask(),
		Gateway:        reqSubnet.GetGateway(),
		NextServer:     reqSubnet.GetNextServer(),
		NameServer:     reqSubnet.GetNameServer(),
		DomainName:     reqSubnet.GetDomainName(),
		ServerUUID:     "",
		LeaderNodeUUID: "",
		OS:             reqSubnet.GetOS(),
		SubnetName:     reqSubnet.GetSubnetName(),
	}

	err = checkSubnet(subnet.NetworkIP, subnet.Netmask, subnet.Gateway, false, nil)
	if err != nil {
		return nil, err
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

	return &subnet, nil
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

// UpdateSubnet : Update infos of a subnet
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

func deleteDHCPConfigFile(serverUUID string) error {
	dhcpdConfLocation := config.DHCPD.ConfigFileLocation + "/" + serverUUID + ".conf"
	err := fileutil.DeleteFile(dhcpdConfLocation)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSubnet : Delete a subnet by UUID
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
