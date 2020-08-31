package dao

import (
	dbsql "database/sql"
	"errors"
	"github.com/golang/protobuf/ptypes"
	pb "hcc/harp/action/grpc/pb/rpcharp"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/pf"
	"net"
	"time"
)

// ReadAdaptiveIPServer : Get a information of AdaptiveIP server setting
func ReadAdaptiveIPServer(serverUUID string) (*pb.AdaptiveIPServer, error) {
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
		logger.Logger.Println(err)
		return nil, err
	}

	adaptiveIPServer.PublicIP = publicIP
	adaptiveIPServer.PrivateIP = privateIP
	adaptiveIPServer.PrivateGateway = privateGateway

	_createdAt, err := ptypes.TimestampProto(createdAt)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	adaptiveIPServer.CreatedAt = _createdAt

	return &adaptiveIPServer, nil
}

// ReadAdaptiveIPServerList : Get the list of AdaptiveIP server settings
func ReadAdaptiveIPServerList(in *pb.ReqGetAdaptiveIPServerList) (*pb.ResGetAdaptiveIPServerList, error) {
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
		return nil, errors.New("please insert row and page arguments or leave arguments as empty state")
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
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&serverUUID, &publicIP, &privateIP, &privateGateway, &createdAt)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		_createdAt, err := ptypes.TimestampProto(createdAt)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
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
		return nil, err
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

	return &adaptiveIPList, nil
}

// ReadAdaptiveIPServerNum : Get the number of AdaptiveIPServer
func ReadAdaptiveIPServerNum() (*pb.ResGetAdaptiveIPServerNum, error) {
	var adaptiveIPServerNum pb.ResGetAdaptiveIPServerNum
	var adaptiveIPServerNr int64

	sql := "select count(*) from adaptiveip_server"
	err := mysql.Db.QueryRow(sql).Scan(&adaptiveIPServerNr)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	adaptiveIPServerNum.Num = adaptiveIPServerNr

	return &adaptiveIPServerNum, nil
}

// CreateAdaptiveIPServer : Create AdaptiveIP of server
func CreateAdaptiveIPServer(in *pb.ReqCreateAdaptiveIPServer) (*pb.AdaptiveIPServer, error) {
	serverUUID := in.ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	publicIP := in.PublicIP
	publicIPOk := len(publicIP) != 0

	if !serverUUIDOk || !publicIPOk {
		return nil, errors.New("need ServerUUID and PublicIP arguments")
	}

	oldAdaptiveIPServer, _ := ReadAdaptiveIPServer(serverUUID)
	if oldAdaptiveIPServer != nil {
		return nil, errors.New("provided ServerUUID is already allocated to one of adaptiveIP")
	}

	subnet, err := ReadSubnetByServer(serverUUID)
	if err != nil {
		return nil, errors.New("provided ServerUUID is not allocated to one of private subnet")
	}

	adaptiveIP := configext.GetAdaptiveIPNetwork()
	netNetwork, _ := iputil.CheckNetwork(adaptiveIP.ExtIfaceIPAddress, adaptiveIP.Netmask)
	mask, _ := iputil.CheckNetmask(adaptiveIP.Netmask)
	netIP := net.IPNet{
		IP:   netNetwork.IP,
		Mask: mask,
	}

	err = iputil.CheckIPisInSubnet(netIP, publicIP)
	if err != nil {
		return nil, err
	}

	adaptiveIPServer := pb.AdaptiveIPServer{
		ServerUUID: serverUUID,
		PublicIP:   publicIP,
	}

	firstIP, _, err := iputil.GetFirstAndLastIPs(subnet.NetworkIP, subnet.Netmask)
	if err != nil {
		return nil, err
	}

	adaptiveIPServer.PrivateIP = firstIP.String()
	adaptiveIPServer.PrivateGateway = subnet.Gateway

	err = pf.CreateAndLoadAnchorConfig(adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP)
	if err != nil {
		return nil, err
	}

	sql := "insert into adaptiveip_server(server_uuid, public_ip, private_ip, private_gateway, created_at) values (?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err := stmt.Exec(adaptiveIPServer.ServerUUID, adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP,
		adaptiveIPServer.PrivateGateway)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	logger.Logger.Println(result.LastInsertId())

	return &adaptiveIPServer, nil
}

// DeleteAdaptiveIPServer : Delete AdaptiveIP of the server
func DeleteAdaptiveIPServer(in *pb.ReqDeleteAdaptiveIPServer) (string, error) {
	var err error

	serverUUID := in.ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	if !serverUUIDOk {
		return "", errors.New("need a server_uuid argument")
	}

	adaptiveIPServer, err := ReadAdaptiveIPServer(serverUUID)
	if adaptiveIPServer == nil {
		return "", errors.New("adaptiveIPServer is nil")
	}

	err = pf.DeleteAndUnloadAnchorConfig(adaptiveIPServer.PublicIP)
	if err != nil {
		logger.Logger.Println(err.Error())
		return "", err
	}

	sql := "delete from adaptiveip_server where server_uuid = ?"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return "", err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err2 := stmt.Exec(serverUUID)
	if err2 != nil {
		logger.Logger.Println(err2)
		return "", err
	}
	logger.Logger.Println(result.RowsAffected())

	return serverUUID, nil
}
