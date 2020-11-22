package dao

import (
	dbsql "database/sql"
	"github.com/golang/protobuf/ptypes"
	pb "hcc/harp/action/grpc/pb/rpcharp"
	"hcc/harp/lib/configext"
	hccerr "hcc/harp/lib/errors"
	"hcc/harp/lib/iptablesext"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/pf"
	"hcc/harp/lib/syscheck"
	"net"
	"strings"
	"time"
)

// ReadAdaptiveIPServer : Get a information of AdaptiveIP server setting
func ReadAdaptiveIPServer(serverUUID string) (*pb.AdaptiveIPServer, uint64, string) {
	var adaptiveIPServer pb.AdaptiveIPServer

	var publicIP string
	var privateIP string
	var privateGateway string
	var createdAt time.Time

	sql := "select server_uuid, public_ip, private_ip, private_gateway, created_at from adaptiveip_server where server_uuid = ?"
	err := mysql.Db.QueryRow(sql, serverUUID).Scan(
		&serverUUID,
		&publicIP,
		&privateIP,
		&privateGateway,
		&createdAt)
	if err != nil {
		errStr := "ReadAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}

	adaptiveIPServer.PublicIP = publicIP
	adaptiveIPServer.PrivateIP = privateIP
	adaptiveIPServer.PrivateGateway = privateGateway

	_createdAt, err := ptypes.TimestampProto(createdAt)
	if err != nil {
		errStr := "ReadAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpInternalTimeStampConversionError, errStr
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
		return nil, hccerr.HarpGrpcArgumentError, "ReadAdaptiveIPServerList(): please insert row and page arguments or leave arguments as empty state"
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
		stmt, err = mysql.Db.Query(sql, row, row*(page-1))
	} else {
		sql += " order by created_at desc"
		stmt, err = mysql.Db.Query(sql)
	}

	if err != nil {
		errStr := "ReadAdaptiveIPServerList(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&serverUUID, &publicIP, &privateIP, &privateGateway, &createdAt)
		if err != nil {
			errStr := "ReadAdaptiveIPServerList(): " + err.Error()
			logger.Logger.Println(errStr)
			if strings.Contains(err.Error(), "no rows in result set") {
				return nil, hccerr.HarpSQLNoResult, errStr
			}
			return nil, hccerr.HarpSQLOperationFail, errStr
		}

		_createdAt, err := ptypes.TimestampProto(createdAt)
		if err != nil {
			logger.Logger.Println(err)
			errStr := "ReadAdaptiveIPServerList(): " + err.Error()
			logger.Logger.Println(errStr)
			return nil, hccerr.HarpInternalTimeStampConversionError, errStr
		}

		adaptiveIPServers = append(adaptiveIPServers, pb.AdaptiveIPServer{
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
		return nil, hccerr.HarpInternalIPAddressError, errStr
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
	err := mysql.Db.QueryRow(sql).Scan(&adaptiveIPServerNr)
	if err != nil {
		errStr := "ReadAdaptiveIPServerNum(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}
	adaptiveIPServerNum.Num = adaptiveIPServerNr

	return &adaptiveIPServerNum, 0, ""
}

// CreateAdaptiveIPServer : Create AdaptiveIP of server
func CreateAdaptiveIPServer(in *pb.ReqCreateAdaptiveIPServer) (*pb.AdaptiveIPServer, uint64, string) {
	serverUUID := in.ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	publicIP := in.PublicIP
	publicIPOk := len(publicIP) != 0

	if !serverUUIDOk || !publicIPOk {
		return nil, hccerr.HarpGrpcArgumentError, "CreateAdaptiveIPServer(): need ServerUUID and PublicIP arguments"
	}

	oldAdaptiveIPServer, _, _ := ReadAdaptiveIPServer(serverUUID)
	if oldAdaptiveIPServer != nil {
		return nil, hccerr.HarpInternalAdaptiveIPAllocatedError, "CreateAdaptiveIPServer(): provided ServerUUID is already allocated to one of adaptiveIP"
	}

	subnet, errCode, _ := ReadSubnetByServer(serverUUID)
	if errCode != 0 {
		return nil, hccerr.HarpInternalSubnetNotAllocatedError, "CreateAdaptiveIPServer(): provided ServerUUID is not allocated to one of private subnet"
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
		return nil, hccerr.HarpInternalIPAddressError, "CreateAdaptiveIPServer(): " + err.Error()
	}

	adaptiveIPServer := pb.AdaptiveIPServer{
		ServerUUID: serverUUID,
		PublicIP:   publicIP,
	}

	firstIP, _, err := iputil.GetFirstAndLastIPs(subnet.NetworkIP, subnet.Netmask)
	if err != nil {
		return nil, hccerr.HarpInternalIPAddressError, "CreateAdaptiveIPServer(): " + err.Error()
	}

	adaptiveIPServer.PrivateIP = firstIP.String()
	adaptiveIPServer.PrivateGateway = subnet.Gateway

	if syscheck.OS == "freebsd" {
		err = pf.CreateAndLoadAnchorConfig(adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP)
	} else {
		err = iptablesext.CreateIPTABLESRulesAndExtIface(adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP)
	}

	if err != nil {
		return nil, hccerr.HarpInternalOperationFail, "CreateAdaptiveIPServer(): " + err.Error()
	}

	sql := "insert into adaptiveip_server(server_uuid, public_ip, private_ip, private_gateway, created_at) values (?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		errStr := "CreateAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(adaptiveIPServer.ServerUUID, adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP,
		adaptiveIPServer.PrivateGateway)
	if err != nil {
		errStr := "CreateAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.HarpSQLOperationFail, errStr
	}

	return &adaptiveIPServer, 0, ""
}

// DeleteAdaptiveIPServer : Delete AdaptiveIP of the server
func DeleteAdaptiveIPServer(in *pb.ReqDeleteAdaptiveIPServer) (string, uint64, string) {
	var err error

	serverUUID := in.ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	if !serverUUIDOk {
		return "", hccerr.HarpGrpcArgumentError, "DeleteAdaptiveIPServer(): need a server_uuid argument"
	}

	adaptiveIPServer, _, _ := ReadAdaptiveIPServer(serverUUID)
	if adaptiveIPServer == nil {
		return "", hccerr.HarpGrpcArgumentError, "DeleteAdaptiveIPServer(): adaptiveIPServer is nil"
	}

	if syscheck.OS == "freebsd" {
		err = pf.DeleteAndUnloadAnchorConfig(adaptiveIPServer.PublicIP)
	} else {
		err = iptablesext.DeleteIPTABLESRulesAndExtIface(adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP)
	}
	if err != nil {
		errStr := "DeleteAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hccerr.HarpInternalOperationFail, errStr
	}

	sql := "delete from adaptiveip_server where server_uuid = ?"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		errStr := "DeleteAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hccerr.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err2 := stmt.Exec(serverUUID)
	if err2 != nil {
		errStr := "DeleteAdaptiveIPServer(): " + err2.Error()
		logger.Logger.Println(errStr)
		return "", hccerr.HarpSQLOperationFail, errStr
	}

	return serverUUID, 0, ""
}
