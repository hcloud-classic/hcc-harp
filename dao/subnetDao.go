package dao

import (
	dbsql "database/sql"
	"errors"
	gouuid "github.com/nu7hatch/gouuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"hcc/harp/action/grpc/client"
	"hcc/harp/daoext"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/iputilext"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"strings"
	"time"
)

// ReadSubnet : Get infos of a subnet
func ReadSubnet(uuid string) (*pb.Subnet, uint64, string) {
	var subnet pb.Subnet

	var groupID int64
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

	sql := "select group_id, network_ip, netmask, gateway, next_server, name_server, domain_name, server_uuid, leader_node_uuid, os, subnet_name, created_at from subnet where uuid = ?"
	row := mysql.Db.QueryRow(sql, uuid)
	err := mysql.QueryRowScan(row,
		&groupID,
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
			return nil, hcc_errors.HarpSQLNoResult, errStr
		}
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}

	subnet.UUID = uuid
	subnet.GroupID = groupID
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
	subnet.CreatedAt = timestamppb.New(createdAt)

	return &subnet, 0, ""
}

// ReadSubnetList : Get list of subnets with selected infos
func ReadSubnetList(in *pb.ReqGetSubnetList) (*pb.ResGetSubnetList, uint64, string) {
	var subnetList pb.ResGetSubnetList
	var subnets []pb.Subnet
	var psubnets []*pb.Subnet

	var uuid string
	var groupID int64
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
		return nil, hcc_errors.HarpGrpcArgumentError, "ReadSubnetList(): please insert row and page arguments or leave arguments as empty state"
	}

	sql := "select * from subnet where 1=1"

	if in.Subnet != nil {
		reqSubnet := in.Subnet

		uuid = reqSubnet.UUID
		uuidOk := len(uuid) != 0
		groupID = reqSubnet.GroupID
		groupIDOk := groupID != 0
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
			sql += " and uuid like '%" + uuid + "%'"
		}
		if groupIDOk {
			sql += " and group_id =" + strconv.Itoa(int(groupID))
		}
		if networkIPOk {
			sql += " and network_ip like '%" + networkIP + "%'"
		}
		if netmaskOk {
			sql += " and netmask like '%" + netmask + "%'"
		}
		if gatewayOk {
			sql += " and gateway like '%" + gateway + "%'"
		}
		if nextServerOk {
			sql += " and next_server like '%" + nextServer + "%'"
		}
		if nameServerOk {
			sql += " and name_server like '%" + nameServer + "%'"
		}
		if domainNameOk {
			sql += " and domain_name like '%" + domainName + "%'"
		}
		if serverUUIDOk {
			sql += " and server_uuid like '%" + serverUUID + "%'"
		}
		if leaderNodeUUIDOk {
			sql += " and leader_node_uuid like '%" + leaderNodeUUID + "%'"
		}
		if osOk {
			sql += " and os like '%" + os + "%'"
		}
		if subnetNameOk {
			sql += " and subnet_name like '%" + subnetName + "%'"
		}
	}

	var stmt *dbsql.Rows
	var err error
	if isLimit {
		sql += " order by created_at desc limit ? offset ?"
		stmt, err = mysql.Query(sql, row, row*(page-1))
	} else {
		sql += " order by created_at desc"
		stmt, err = mysql.Query(sql)
	}

	if err != nil {
		errStr := "ReadSubnetList(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &groupID, &networkIP, &netmask, &gateway, &nextServer, &nameServer, &domainName, &serverUUID, &leaderNodeUUID, &os, &subnetName, &createdAt)
		if err != nil {
			errStr := "ReadSubnetList(): " + err.Error()
			logger.Logger.Println(errStr)
			if strings.Contains(err.Error(), "no rows in result set") {
				return nil, hcc_errors.HarpSQLNoResult, errStr
			}
			return nil, hcc_errors.HarpSQLOperationFail, errStr
		}

		subnets = append(subnets, pb.Subnet{
			UUID:           uuid,
			GroupID:        groupID,
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
			CreatedAt:      timestamppb.New(createdAt)})
	}

	for i := range subnets {
		psubnets = append(psubnets, &subnets[i])
	}

	subnetList.Subnet = psubnets

	return &subnetList, 0, ""
}

// ReadAvailableSubnetList : Get list of available subnets
func ReadAvailableSubnetList(in *pb.ReqGetAvailableSubnetList) (*pb.ResGetSubnetList, uint64, string) {
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

	var groupID = in.GetGroupID()
	if groupID == 0 {
		return nil, hcc_errors.HarpGrpcArgumentError, "ReadAvailableSubnetList(): please insert a group_id argument"
	}

	sql := "select * from subnet where server_uuid = '' and group_id = " +
		strconv.Itoa(int(groupID)) + " order by created_at desc"
	stmt, err := mysql.Query(sql)
	if err != nil {
		errStr := "ReadAvailableSubnetList(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &groupID, &networkIP, &netmask, &gateway, &nextServer, &nameServer, &domainName, &serverUUID, &leaderNodeUUID, &os, &subnetName, &createdAt)
		if err != nil {
			errStr := "ReadAvailableSubnetList(): " + err.Error()
			logger.Logger.Println(errStr)
			if strings.Contains(err.Error(), "no rows in result set") {
				return nil, hcc_errors.HarpSQLNoResult, errStr
			}
			return nil, hcc_errors.HarpSQLOperationFail, errStr
		}

		subnets = append(subnets, pb.Subnet{
			UUID:           uuid,
			GroupID:        groupID,
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
			CreatedAt:      timestamppb.New(createdAt)})
	}

	for i := range subnets {
		psubnets = append(psubnets, &subnets[i])
	}

	subnetList.Subnet = psubnets

	return &subnetList, 0, ""
}

// ReadSubnetNum : Get the number of subnets
func ReadSubnetNum(in *pb.ReqGetSubnetNum) (*pb.ResGetSubnetNum, uint64, string) {
	var resSubnetNum pb.ResGetSubnetNum
	var subnetNr int64
	var groupID = in.GetGroupID()

	sql := "select count(*) from subnet"
	if groupID != 0 {
		sql = "select count(*) from subnet where group_id = " + strconv.Itoa(int(groupID))
	}
	row := mysql.Db.QueryRow(sql)
	err := mysql.QueryRowScan(row, &subnetNr)
	if err != nil {
		errStr := "ReadSubnetNum(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	resSubnetNum.Num = subnetNr

	return &resSubnetNum, 0, ""
}

func checkGroupIDExist(groupID int64) error {
	resGetGroupList, hccErrStack := client.RC.GetGroupList()
	if hccErrStack != nil {
		return hccErrStack.Pop().ToError()
	}

	for _, pGroup := range resGetGroupList.Group {
		if pGroup.Id == groupID {
			return nil
		}
	}

	return errors.New("given group ID is not in the database")
}

func checkSubnet(networkIP string, netmask string, gateway string, skipMine bool, oldSubnet *pb.Subnet,
	resValidCheckSubnet *pb.ResValidCheckSubnet, isUpdate bool) error {
	isConflict, err := iputilext.CheckSubnetConflict(networkIP, netmask, skipMine, oldSubnet, resValidCheckSubnet)
	if isConflict {
		return errors.New("given subnet is conflicted with one of subnet that stored in the database")
	}
	if err != nil {
		return err
	}

	netNetwork, _ := iputil.CheckNetwork(networkIP, netmask)

	isPrivate, _ := iputilext.CheckPrivateSubnet(netNetwork.IP.String(), netmask)
	if !isPrivate {
		if resValidCheckSubnet != nil {
			resValidCheckSubnet.ErrorCode = daoext.SubnetValidErrorNotPrivate
		}
		return errors.New("given network IP address is not in private network")
	}

	firstIP, _, _ := iputil.GetFirstAndLastIPs(networkIP, netmask)
	if firstIP.To4()[3] != 1 {
		if resValidCheckSubnet != nil {
			resValidCheckSubnet.ErrorCode = daoext.SubnetValidErrorStartIPNot1
		}
		return errors.New("start IP address must be x.x.x.1")
	}

	if !isUpdate {
		err := iputil.CheckIPisInSubnet(*netNetwork, gateway)
		if err != nil {
			if resValidCheckSubnet != nil {
				if strings.Contains(err.Error(), "wrong") ||
					strings.Contains(err.Error(), "network") ||
					strings.Contains(err.Error(), "broadcast") {
					resValidCheckSubnet.ErrorCode = daoext.SubnetValidErrorInvalidGatewayAddress
				} else {
					resValidCheckSubnet.ErrorCode = daoext.SubnetValidErrorGatewayNotInSubnet
				}
			}
			return err
		}

		err = iputil.CheckSubnetIsUsedByIface(*netNetwork)
		if err != nil {
			if resValidCheckSubnet != nil && strings.Contains(err.Error(), "conflicted") {
				resValidCheckSubnet.ErrorCode = daoext.SubnetValidErrorSubnetIsUsedByIface
			}
			return err
		}
	}

	return nil
}

func checkServerUUID(serverUUID string) *hcc_errors.HccErrorStack {
	serverUUIDs, errStack := client.RC.AllServerUUID() // passing HccErrorStack to err
	if errStack != nil {
		return errStack
	}

	for i := range serverUUIDs {
		if serverUUIDs[i] == serverUUID {
			return nil
		}
	}

	hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpSQLNoResult, "given server UUID is not in the database"))
	return hccErrStack
}

func checkCreateSubnetArgs(reqSubnet *pb.Subnet) bool {
	groupIDOk := reqSubnet.GroupID != 0
	networkIPOk := len(reqSubnet.GetNetworkIP()) != 0
	netmaskOk := len(reqSubnet.GetNetmask()) != 0
	gatewayOk := len(reqSubnet.GetGateway()) != 0
	nextServerOk := len(reqSubnet.GetNextServer()) != 0
	nameServerOk := len(reqSubnet.GetNameServer()) != 0
	domainNameOk := len(reqSubnet.GetDomainName()) != 0
	osOk := len(reqSubnet.GetOS()) != 0
	subnetNameOk := len(reqSubnet.GetSubnetName()) != 0

	return !(groupIDOk && networkIPOk && netmaskOk && gatewayOk && nextServerOk && nameServerOk && domainNameOk && osOk && subnetNameOk)
}

// CreateSubnet : Create a subnet
func CreateSubnet(in *pb.ReqCreateSubnet) (*pb.Subnet, uint64, string) {
	reqSubnet := in.GetSubnet()
	if reqSubnet == nil {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreateSubnet(): subnet is nil"
	}

	out, err := gouuid.NewV4()
	if err != nil {
		errStr := "CreateSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpInternalUUIDGenerationError, errStr
	}
	uuid := out.String()

	if checkCreateSubnetArgs(reqSubnet) {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreateSubnet(): some of arguments are missing"
	}

	subnet := pb.Subnet{
		UUID:           uuid,
		GroupID:        reqSubnet.GetGroupID(),
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

	err = checkGroupIDExist(subnet.GroupID)
	if err != nil {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreateSubnet(): " + err.Error()
	}

	resGetQuota, errStack := client.RC.GetQuota(subnet.GroupID)
	if errStack != nil {
		return nil, hcc_errors.HarpGrpcRequestError, "CreateSubnet(): " + errStack.Pop().Text()
	}
	subnetNum, errCode, errText := ReadSubnetNum(&pb.ReqGetSubnetNum{GroupID: subnet.GroupID})
	if errCode != 0 {
		return nil, errCode, "CreateSubnet(): " + errText
	}
	if subnetNum.Num+1 > int64(resGetQuota.Quota.LimitSubnetCnt) {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreateSubnet(): Subnet count quota exceeded"
	}

	err = checkSubnet(subnet.NetworkIP, subnet.Netmask, subnet.Gateway, false, nil, nil, false)
	if err != nil {
		return nil, hcc_errors.HarpInternalIPAddressError, "CreateSubnet(): " + err.Error()
	}

	sql := "insert into subnet(uuid, group_id, network_ip, netmask, gateway, next_server, name_server, domain_name, server_uuid, leader_node_uuid, os, subnet_name, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "CreateSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(subnet.UUID, subnet.GroupID, subnet.NetworkIP, subnet.Netmask, subnet.Gateway, subnet.NextServer, subnet.NameServer, subnet.DomainName, subnet.ServerUUID, subnet.LeaderNodeUUID, subnet.OS, subnet.SubnetName)
	if err != nil {
		errStr := "CreateSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}

	return &subnet, 0, ""
}

func checkValidCheckSubnetArgs(reqSubnet *pb.Subnet) bool {
	networkIPOk := len(reqSubnet.GetNetworkIP()) != 0
	netmaskOk := len(reqSubnet.GetNetmask()) != 0
	gatewayOk := len(reqSubnet.GetGateway()) != 0

	return !(networkIPOk && netmaskOk && gatewayOk)
}

// ValidCheckSubnet : Check if we can create the subnet with provided network address and subnet mask, gateway
func ValidCheckSubnet(in *pb.ReqValidCheckSubnet) *pb.ResValidCheckSubnet {
	var oldSubnet *pb.Subnet

	reqSubnet := in.GetSubnet()
	if reqSubnet == nil {
		return &pb.ResValidCheckSubnet{
			ErrorCode: daoext.SubnetValidErrorArgumentError,
		}
	}

	if checkValidCheckSubnetArgs(reqSubnet) {
		return &pb.ResValidCheckSubnet{
			ErrorCode: daoext.SubnetValidErrorArgumentError,
		}
	}

	subnet := pb.Subnet{
		NetworkIP: reqSubnet.GetNetworkIP(),
		Netmask:   reqSubnet.GetNetmask(),
		Gateway:   reqSubnet.GetGateway(),
	}

	if in.GetIsUpdate() {
		oldSubnet, _, _ = ReadSubnet(reqSubnet.GetUUID())
	}

	var resValidCheckSubnet pb.ResValidCheckSubnet
	err := checkSubnet(subnet.NetworkIP, subnet.Netmask, subnet.Gateway, in.GetIsUpdate(), oldSubnet,
		&resValidCheckSubnet, in.GetIsUpdate())
	if err != nil {
		return &pb.ResValidCheckSubnet{
			ErrorCode: resValidCheckSubnet.ErrorCode,
		}
	}

	return &pb.ResValidCheckSubnet{
		ErrorCode: daoext.SubnetValid,
	}
}

func checkUpdateSubnetArgs(reqSubnet *pb.Subnet) bool {
	networkIPOk := len(reqSubnet.GetNetworkIP()) != 0
	netmaskOk := len(reqSubnet.GetNetmask()) != 0
	gatewayOk := len(reqSubnet.GetGateway()) != 0
	nextServerOk := len(reqSubnet.GetNextServer()) != 0
	nameServerOk := len(reqSubnet.GetNameServer()) != 0
	domainNameOk := len(reqSubnet.GetDomainName()) != 0
	osOk := len(reqSubnet.GetOS()) != 0
	subnetNameOk := len(reqSubnet.GetSubnetName()) != 0

	return !networkIPOk && !netmaskOk && !gatewayOk && !nextServerOk && !nameServerOk && !domainNameOk && !osOk && !subnetNameOk
}

// UpdateSubnet : Update infos of the subnet
func UpdateSubnet(in *pb.ReqUpdateSubnet) (*pb.Subnet, uint64, string) {
	if in.Subnet == nil {
		return nil, hcc_errors.HarpGrpcArgumentError, "UpdateSubnet(): subnet is nil"
	}
	reqSubnet := in.Subnet

	requestedUUID := reqSubnet.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return nil, hcc_errors.HarpGrpcArgumentError, "UpdateSubnet(): need a uuid argument"
	}

	if checkUpdateSubnetArgs(reqSubnet) {
		return nil, hcc_errors.HarpGrpcArgumentError, "UpdateSubnet(): need some arguments"
	}

	var networkIP string
	var netmask string
	var gateway string
	var nextServer string
	var nameServer string
	var domainName string
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
	subnet.OS = os
	subnet.SubnetName = subnetName

	oldSubnet, errCode, errStr := ReadSubnet(subnet.GetUUID())
	if errCode != 0 {
		return nil, errCode, "UpdateSubnet(): " + errStr
	}

	if oldSubnet.ServerUUID != "" {
		return nil, hcc_errors.HarpInternalSubnetInUseError, "UpdateSubnet(): Subnet is in use by the server (ServerUUID=" + oldSubnet.ServerUUID + ")"
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

	err := checkSubnet(subnet.NetworkIP, subnet.Netmask, subnet.Gateway, true, oldSubnet, nil, true)
	if err != nil {
		return nil, hcc_errors.HarpInternalIPAddressError, "UpdateSubnet(): " + err.Error()
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
	if osOk {
		updateSet += " os = '" + subnet.OS + "', "
	}
	if subnetNameOk {
		updateSet += " subnet_name = '" + subnet.SubnetName + "', "
	}
	sql += updateSet[0:len(updateSet)-2] + " where uuid = ?"

	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "UpdateSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(subnet.UUID)
	if err != nil {
		errStr := "UpdateSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
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
		return nil, hcc_errors.HarpGrpcArgumentError, "DeleteSubnet(): need a uuid argument"
	}

	subnet, errCode, errStr := ReadSubnet(requestedUUID)
	if errCode != 0 {
		return nil, errCode, "DeleteSubnet(): " + errStr
	}

	if len(subnet.ServerUUID) != 0 {
		errStr := "DeleteSubnet(): subnet is used by the server (UUID:" + subnet.ServerUUID + ")"
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpInternalSubnetInUseError, errStr
	}

	sql := "delete from subnet where uuid = ?"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "DeleteSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(requestedUUID)
	if err != nil {
		errStr := "DeleteSubnet(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}

	return subnet, 0, ""
}
