package dao

import (
	dbsql "database/sql"
	"github.com/golang/protobuf/ptypes"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/iptablesext"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/pf"
	"hcc/harp/lib/syscheck"
	"net"
	"strconv"
	"strings"
	"time"
)

// ReadAdaptiveIPServer : Get a information of AdaptiveIP server setting
func ReadAdaptiveIPServer(serverUUID string) (*pb.AdaptiveIPServer, uint64, string) {
	var adaptiveIPServer pb.AdaptiveIPServer

	var groupID int64
	var publicIP string
	var privateIP string
	var privateGateway string
	var createdAt time.Time

	sql := "select group_id, public_ip, private_ip, private_gateway, created_at from adaptiveip_server where server_uuid = ?"
	row := mysql.Db.QueryRow(sql, serverUUID)
	err := mysql.QueryRowScan(row,
		&groupID,
		&publicIP,
		&privateIP,
		&privateGateway,
		&createdAt)
	if err != nil {
		errStr := "ReadAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}

	adaptiveIPServer.ServerUUID = serverUUID
	adaptiveIPServer.GroupID = groupID
	adaptiveIPServer.PublicIP = publicIP
	adaptiveIPServer.PrivateIP = privateIP
	adaptiveIPServer.PrivateGateway = privateGateway

	_createdAt, err := ptypes.TimestampProto(createdAt)
	if err != nil {
		errStr := "ReadAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpInternalTimeStampConversionError, errStr
	}
	adaptiveIPServer.CreatedAt = _createdAt

	return &adaptiveIPServer, 0, ""
}

// ReadAdaptiveIPServerList : Get the list of AdaptiveIP server settings
func ReadAdaptiveIPServerList(in *pb.ReqGetAdaptiveIPServerList) (*pb.ResGetAdaptiveIPServerList, uint64, string) {
	var adaptiveIPList pb.ResGetAdaptiveIPServerList
	var adaptiveIPServers []pb.AdaptiveIPServer
	var padaptiveIPServers []*pb.AdaptiveIPServer

	var serverUUID string
	var groupID int64
	var publicIP string
	var privateIP string
	var privateGateway string
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
		return nil, hcc_errors.HarpGrpcArgumentError, "ReadAdaptiveIPServerList(): please insert row and page arguments or leave arguments as empty state"
	}

	sql := "select * from adaptiveip_server where 1=1"

	if in.AdaptiveipServer != nil {
		reqAdaptiveIPServer := in.AdaptiveipServer

		publicIP = reqAdaptiveIPServer.PublicIP
		publicIPOk := len(publicIP) != 0
		privateIP = reqAdaptiveIPServer.PrivateIP
		privateIPOk := len(privateIP) != 0
		privateGateway = reqAdaptiveIPServer.PrivateGateway
		privateGatewayOk := len(privateGateway) != 0

		if publicIPOk {
			sql += " and public_ip = '" + publicIP + "'"
		}
		if privateIPOk {
			sql += " and private_ip = '" + privateIP + "'"
		}
		if privateGatewayOk {
			sql += " and private_gateway = '" + privateGateway + "'"
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
		errStr := "ReadAdaptiveIPServerList(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&serverUUID, &groupID, &publicIP, &privateIP, &privateGateway, &createdAt)
		if err != nil {
			errStr := "ReadAdaptiveIPServerList(): " + err.Error()
			logger.Logger.Println(errStr)
			if strings.Contains(err.Error(), "no rows in result set") {
				return nil, hcc_errors.HarpSQLNoResult, errStr
			}
			return nil, hcc_errors.HarpSQLOperationFail, errStr
		}

		_createdAt, err := ptypes.TimestampProto(createdAt)
		if err != nil {
			logger.Logger.Println(err)
			errStr := "ReadAdaptiveIPServerList(): " + err.Error()
			logger.Logger.Println(errStr)
			return nil, hcc_errors.HarpInternalTimeStampConversionError, errStr
		}

		adaptiveIPServers = append(adaptiveIPServers, pb.AdaptiveIPServer{
			GroupID: groupID,
			ServerUUID:     serverUUID,
			PublicIP:       publicIP,
			PrivateIP:      privateIP,
			PrivateGateway: privateGateway,
			CreatedAt:      _createdAt,
		})
	}

	adaptiveIP := configext.GetAdaptiveIPNetwork()
	netNetwork, err := iputil.CheckNetwork(adaptiveIP.ExtIfaceIPAddress, adaptiveIP.Netmask)
	if err != nil {
		logger.Logger.Println(err)
		errStr := "ReadAdaptiveIPServerList(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpInternalIPAddressError, errStr
	}

	for i := range adaptiveIPServers {
		netIP := iputil.CheckValidIP(adaptiveIPServers[i].PublicIP)
		if !netNetwork.Contains(netIP) {
			adaptiveIPServers[i].Status = "Invalid"
			continue
		}
		adaptiveIPServers[i].Status = "Using"

		padaptiveIPServers = append(padaptiveIPServers, &adaptiveIPServers[i])
	}

	adaptiveIPList.AdaptiveipServer = padaptiveIPServers

	return &adaptiveIPList, 0, ""
}

// ReadAdaptiveIPServerNum : Get the number of AdaptiveIPServer
func ReadAdaptiveIPServerNum() (*pb.ResGetAdaptiveIPServerNum, uint64, string) {
	var adaptiveIPServerNum pb.ResGetAdaptiveIPServerNum
	var adaptiveIPServerNr int64

	sql := "select count(*) from adaptiveip_server"
	row := mysql.Db.QueryRow(sql)
	err := mysql.QueryRowScan(row, &adaptiveIPServerNr)
	if err != nil {
		errStr := "ReadAdaptiveIPServerNum(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	adaptiveIPServerNum.Num = adaptiveIPServerNr

	return &adaptiveIPServerNum, 0, ""
}

// CreateAdaptiveIPServer : Create AdaptiveIP of server
func CreateAdaptiveIPServer(in *pb.ReqCreateAdaptiveIPServer) (*pb.AdaptiveIPServer, uint64, string) {
	serverUUID := in.ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	groupID := in.GroupID
	groupIDOk := groupID != 0
	publicIP := in.PublicIP
	publicIPOk := len(publicIP) != 0

	if !serverUUIDOk || !groupIDOk || !publicIPOk {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreateAdaptiveIPServer(): need ServerUUID and GroupID, PublicIP arguments"
	}

	oldAdaptiveIPServer, _, _ := ReadAdaptiveIPServer(serverUUID)
	if oldAdaptiveIPServer != nil {
		return nil, hcc_errors.HarpInternalAdaptiveIPAllocatedError, "CreateAdaptiveIPServer(): provided ServerUUID is already allocated to one of adaptiveIP"
	}

	subnet, errCode, _ := ReadSubnetByServer(serverUUID)
	if errCode != 0 {
		return nil, hcc_errors.HarpInternalSubnetNotAllocatedError, "CreateAdaptiveIPServer(): provided ServerUUID is not allocated to one of private subnet"
	}

	adaptiveIP := configext.GetAdaptiveIPNetwork()
	netNetwork, _ := iputil.CheckNetwork(adaptiveIP.ExtIfaceIPAddress, adaptiveIP.Netmask)
	mask, _ := iputil.CheckNetmask(adaptiveIP.Netmask)
	netIP := net.IPNet{
		IP:   netNetwork.IP,
		Mask: mask,
	}

	err := iputil.CheckIPisInSubnet(netIP, publicIP)
	if err != nil {
		return nil, hcc_errors.HarpInternalIPAddressError, "CreateAdaptiveIPServer(): " + err.Error()
	}

	var startIPSum = 0
	var endIPSsum = 0
	var publicIPSum = 0

	startIPSplit := strings.Split(adaptiveIP.StartIPAddress, ".")
	endIPSplit := strings.Split(adaptiveIP.EndIPAddress, ".")
	publicIPSplit := strings.Split(publicIP, ".")

	for _, startIPSplited := range startIPSplit {
		num, _ := strconv.Atoi(startIPSplited)
		startIPSum += num
	}
	for _, endIPSplited := range endIPSplit {
		num, _ := strconv.Atoi(endIPSplited)
		endIPSsum += num
	}
	for _, publicIPSplited := range publicIPSplit {
		num, _ := strconv.Atoi(publicIPSplited)
		publicIPSum += num
	}

	if publicIPSum < startIPSum || publicIPSum > endIPSsum {
		return nil, hcc_errors.HarpInternalIPAddressError,
			"CreateAdaptiveIPServer(): Provided public IP address is out of range. Check AdaptiveIP settings."
	}

	adaptiveIPServer := pb.AdaptiveIPServer{
		ServerUUID: serverUUID,
		GroupID: groupID,
		PublicIP:   publicIP,
	}

	firstIP, _, err := iputil.GetFirstAndLastIPs(subnet.NetworkIP, subnet.Netmask)
	if err != nil {
		return nil, hcc_errors.HarpInternalIPAddressError, "CreateAdaptiveIPServer(): " + err.Error()
	}

	adaptiveIPServer.PrivateIP = firstIP.String()
	adaptiveIPServer.PrivateGateway = subnet.Gateway

	if syscheck.OS == "freebsd" {
		err = pf.CreateAndLoadAnchorConfig(adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP)
	} else {
		err = iptablesext.CreateIPTABLESRulesAndExtIface(adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP)
	}

	if err != nil {
		return nil, hcc_errors.HarpInternalOperationFail, "CreateAdaptiveIPServer(): " + err.Error()
	}

	sql := "insert into adaptiveip_server(server_uuid, group_id, public_ip, private_ip, private_gateway, created_at) values (?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "CreateAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(adaptiveIPServer.ServerUUID, adaptiveIPServer.GroupID, adaptiveIPServer.PublicIP,
		adaptiveIPServer.PrivateIP, adaptiveIPServer.PrivateGateway)
	if err != nil {
		errStr := "CreateAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}

	return &adaptiveIPServer, 0, ""
}

// DeleteAdaptiveIPServer : Delete AdaptiveIP of the server
func DeleteAdaptiveIPServer(in *pb.ReqDeleteAdaptiveIPServer) (string, uint64, string) {
	var err error

	serverUUID := in.ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	if !serverUUIDOk {
		return "", hcc_errors.HarpGrpcArgumentError, "DeleteAdaptiveIPServer(): need a server_uuid argument"
	}

	adaptiveIPServer, _, _ := ReadAdaptiveIPServer(serverUUID)
	if adaptiveIPServer == nil {
		return "", hcc_errors.HarpGrpcArgumentError, "DeleteAdaptiveIPServer(): adaptiveIPServer is nil"
	}

	if syscheck.OS == "freebsd" {
		err = pf.DeleteAndUnloadAnchorConfig(adaptiveIPServer.PublicIP)
	} else {
		err = iptablesext.DeleteIPTABLESRulesAndExtIface(adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP)
	}
	if err != nil {
		errStr := "DeleteAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hcc_errors.HarpInternalOperationFail, errStr
	}

	sql := "delete from adaptiveip_server where server_uuid = ?"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "DeleteAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err2 := stmt.Exec(serverUUID)
	if err2 != nil {
		errStr := "DeleteAdaptiveIPServer(): " + err2.Error()
		logger.Logger.Println(errStr)
		return "", hcc_errors.HarpSQLOperationFail, errStr
	}

	return serverUUID, 0, ""
}
