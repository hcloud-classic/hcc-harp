package daoext

import (
	sql2 "database/sql"
	"google.golang.org/protobuf/types/known/timestamppb"
	"hcc/harp/lib/adaptiveipext"
	"hcc/harp/lib/configadapriveipnetwork"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
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
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hcc_errors.HarpSQLNoResult, errStr
		}
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}

	adaptiveIPServer.ServerUUID = serverUUID
	adaptiveIPServer.GroupID = groupID
	adaptiveIPServer.PublicIP = publicIP
	adaptiveIPServer.PrivateIP = privateIP
	adaptiveIPServer.PrivateGateway = privateGateway
	adaptiveIPServer.CreatedAt = timestamppb.New(createdAt)

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

		serverUUID = reqAdaptiveIPServer.ServerUUID
		serverUUIDOk := len(serverUUID) != 0
		groupID = reqAdaptiveIPServer.GroupID
		groupIDOk := groupID != 0
		publicIP = reqAdaptiveIPServer.PublicIP
		publicIPOk := len(publicIP) != 0
		privateIP = reqAdaptiveIPServer.PrivateIP
		privateIPOk := len(privateIP) != 0
		privateGateway = reqAdaptiveIPServer.PrivateGateway
		privateGatewayOk := len(privateGateway) != 0

		if serverUUIDOk {
			sql += " and server_uuid like '%" + serverUUID + "%'"
		}
		if groupIDOk {
			sql += " and group_id = " + strconv.Itoa(int(groupID))
		}
		if publicIPOk {
			sql += " and public_ip like '%" + publicIP + "%'"
		}
		if privateIPOk {
			sql += " and private_ip like '%" + privateIP + "%'"
		}
		if privateGatewayOk {
			sql += " and private_gateway like '%" + privateGateway + "%'"
		}
	}

	var stmt *sql2.Rows
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

		adaptiveIPServers = append(adaptiveIPServers, pb.AdaptiveIPServer{
			GroupID:        groupID,
			ServerUUID:     serverUUID,
			PublicIP:       publicIP,
			PrivateIP:      privateIP,
			PrivateGateway: privateGateway,
			CreatedAt:      timestamppb.New(createdAt),
		})
	}

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()
	netNetwork, err := iputil.CheckNetwork(adaptiveIP.ExtIfaceIPAddress, adaptiveIP.Netmask)
	if err != nil {
		logger.Logger.Println(err)
		errStr := "ReadAdaptiveIPServerList(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpInternalIPAddressError, errStr
	}

	for i := range adaptiveIPServers {
		internalIP, err := adaptiveipext.ExternalIPtoInternalIP(adaptiveIPServers[i].PublicIP)
		if err != nil {
			adaptiveIPServers[i].Status = "Invalid"
			continue
		}
		netIP := iputil.CheckValidIP(internalIP)
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
