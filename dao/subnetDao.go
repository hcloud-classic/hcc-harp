package dao

import (
	dbsql "database/sql"
	"errors"
	"github.com/golang/protobuf/ptypes"
	gouuid "github.com/nu7hatch/gouuid"
	"hcc/harp/action/grpc/client"
	pb "hcc/harp/action/grpc/pb/rpcharp"
	hccerr "hcc/harp/lib/errors"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"strings"
	"time"
)

// ReadSubnet : Get infos of a subnet
func ReadSubnet(uuid string) (*pb.Subnet, uint64, string) {
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
		errStr := "ReadSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hccerr.HarpSQLNoResult, errStr
		}
		return nil, hccerr.HarpSQLOperationFail, errStr
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
		errStr := "ReadSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpInternalTimeStampConversionError, errStr
	}

	return &subnet, 0, ""
}

// ReadSubnetByServer : Get infos of a subnet by server UUID
func ReadSubnetByServer(serverUUID string) (*pb.Subnet, uint64, string) {
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
		errStr := "ReadSubnetByServer(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hccerr.HarpSQLNoResult, errStr
		}
		return nil, hccerr.HarpSQLOperationFail, errStr
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
		errStr := "ReadSubnetByServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpInternalTimeStampConversionError, errStr
	}

	return &subnet, 0, ""
}

// ReadSubnetList : Get list of subnets with selected infos
func ReadSubnetList(in *pb.ReqGetSubnetList) (*pb.ResGetSubnetList, uint64, string) {
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
		return nil, hccerr.HarpGrpcArgumentError, "ReadSubnetList(): please insert row and page arguments or leave arguments as empty state"
	}

	sql := "select * from subnet where 1=1"

	if in.Subnet != nil {
		reqSubnet := in.Subnet

		uuid = reqSubnet.UUID
		uuidOk := len(uuid) != 0
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

		if uuidOk {
			sql += " and uuid = '" + uuid + "'"
		}
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
		errStr := "ReadSubnetList(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &networkIP, &netmask, &gateway, &nextServer, &nameServer, &domainName, &serverUUID, &leaderNodeUUID, &os, &subnetName, &createdAt)
		if err != nil {
			errStr := "ReadSubnetList(): " + err.Error()
			logger.Logger.Println(errStr)
			if strings.Contains(err.Error(), "no rows in result set") {
				return nil, hccerr.HarpSQLNoResult, errStr
			}
			return nil, hccerr.HarpSQLOperationFail, errStr
		}

		_createdAt, err := ptypes.TimestampProto(createdAt)
		if err != nil {
			errStr := "ReadSubnetList(): " + err.Error()
			logger.Logger.Println(errStr)
			return nil, hccerr.HarpInternalTimeStampConversionError, errStr
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

	return &subnetList, 0, ""
}

// ReadSubnetNum : Get the number of subnets
func ReadSubnetNum() (*pb.ResGetSubnetNum, uint64, string) {
	var resSubnetNum pb.ResGetSubnetNum
	var subnetNr int64

	sql := "select count(*) from subnet"
	err := mysql.Db.QueryRow(sql).Scan(&subnetNr)
	if err != nil {
		errStr := "ReadSubnetNum(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}
	resSubnetNum.Num = subnetNr

	return &resSubnetNum, 0, ""
}

func checkSubnet(networkIP string, netmask string, gateway string, skipMine bool, oldSubnet *pb.Subnet) error {
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

	isPrivate, err := iputil.CheckPrivateSubnet(netNetwork.IP.String(), netmask)
	if !isPrivate {
		return errors.New("given network IP address is not in private network")
	}
	if err != nil {
		return err
	}

	err = iputil.CheckIPisInSubnet(*netNetwork, gateway)
	if err != nil {
		return err
	}

	return nil
}

func checkServerUUID(serverUUID string) *hccerr.HccErrorStack {
	serverUUIDs, errStack := client.RC.AllServerUUID() // passing HccErrorStack to err
	if errStack != nil {
		return errStack
	}

	for i := range serverUUIDs {
		if serverUUIDs[i] == serverUUID {
			return nil
		}
	}

	hccErrStack := hccerr.ReturnHccError(hccerr.HarpSQLNoResult, "given server UUID is not in the database")
	return &hccErrStack
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
func CreateSubnet(in *pb.ReqCreateSubnet) (*pb.Subnet, uint64, string) {
	reqSubnet := in.GetSubnet()
	if reqSubnet == nil {
		return nil, hccerr.HarpGrpcArgumentError, "CreateSubnet(): subnet is nil"
	}

	out, err := gouuid.NewV4()
	if err != nil {
		errStr := "CreateSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpInternalUUIDGenerationError, errStr
	}
	uuid := out.String()

	if checkCreateSubnetArgs(reqSubnet) {
		return nil, hccerr.HarpGrpcArgumentError, "CreateSubnet(): some of arguments are missing"
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
		return nil, hccerr.HarpInternalIPAddressError, "CreateSubnet(): " + err.Error()
	}

	sql := "insert into subnet(uuid, network_ip, netmask, gateway, next_server, name_server, domain_name, server_uuid, leader_node_uuid, os, subnet_name, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		errStr := "CreateSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(subnet.UUID, subnet.NetworkIP, subnet.Netmask, subnet.Gateway, subnet.NextServer, subnet.NameServer, subnet.DomainName, subnet.ServerUUID, subnet.LeaderNodeUUID, subnet.OS, subnet.SubnetName)
	if err != nil {
		errStr := "CreateSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}

	return &subnet, 0, ""
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

// UpdateSubnet : Update infos of the subnet
func UpdateSubnet(in *pb.ReqUpdateSubnet) (*pb.Subnet, uint64, string) {
	if in.Subnet == nil {
		return nil, hccerr.HarpGrpcArgumentError, "UpdateSubnet(): subnet is nil"
	}
	reqSubnet := in.Subnet

	requestedUUID := reqSubnet.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return nil, hccerr.HarpGrpcArgumentError, "UpdateSubnet(): need a uuid argument"
	}

	if checkUpdateSubnetArgs(reqSubnet) {
		return nil, hccerr.HarpGrpcArgumentError, "UpdateSubnet(): need some arguments"
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

	oldSubnet, errCode, errStr := ReadSubnet(subnet.GetUUID())
	if errCode != 0 {
		return nil, errCode, "UpdateSubnet(): " + errStr
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

	err := checkSubnet(subnet.NetworkIP, subnet.Netmask, subnet.Gateway, true, oldSubnet)
	if err != nil {
		return nil, hccerr.HarpInternalIPAddressError, "UpdateSubnet(): " + err.Error()
	}

	if serverUUIDOk && subnet.ServerUUID != "-" {
		errStack := checkServerUUID(subnet.ServerUUID)
		if errStack != nil {
			return nil, (*errStack)[errStack.Len()].ErrCode, "UpdateSubnet(): " + (*errStack)[errStack.Len()].ErrText
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
		if subnet.ServerUUID == "-" {
			subnet.ServerUUID = ""
		}
		updateSet += " server_uuid = '" + subnet.ServerUUID + "', "
	}
	if leaderNodeUUIDOk {
		if subnet.LeaderNodeUUID == "-" {
			subnet.LeaderNodeUUID = ""
		}
		updateSet += " leader_node_uuid = '" + subnet.LeaderNodeUUID + "', "
	}
	if osOk {
		updateSet += " os = '" + subnet.OS + "', "
	}
	if subnetNameOk {
		updateSet += " subnet_name = '" + subnet.SubnetName + "', "
	}
	sql += updateSet[0:len(updateSet)-2] + " where uuid = ?"

	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		errStr := "UpdateSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err2 := stmt.Exec(subnet.UUID)
	if err2 != nil {
		errStr := "UpdateSubnet(): " + err2.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}

	subnet, errCode, errStr = ReadSubnet(subnet.UUID)
	if errCode != 0 {
		logger.Logger.Println("UpdateSubnet(): " + errStr)
	}

	return subnet, 0, ""
}

// DeleteSubnet : Delete a subnet by UUID
func DeleteSubnet(in *pb.ReqDeleteSubnet) (*pb.Subnet, uint64, string) {
	var err error

	requestedUUID := in.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return nil, hccerr.HarpGrpcArgumentError, "DeleteSubnet(): need a uuid argument"
	}

	subnet, errCode, errStr := ReadSubnet(requestedUUID)
	if errCode != 0 {
		return nil, errCode, "DeleteSubnet(): " + errStr
	}

	if len(subnet.ServerUUID) != 0 {
		errStr := "DeleteSubnet(): subnet is used by the server (UUID:" + subnet.ServerUUID + ")"
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpInternalSubnetInUseError, errStr
	}

	sql := "delete from subnet where uuid = ?"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		errStr := "DeleteSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err2 := stmt.Exec(requestedUUID)
	if err2 != nil {
		errStr := "DeleteSubnet(): " + err2.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}

	return subnet, 0, ""
}
