package dao
//
//import (
//	dbsql "database/sql"
//	"errors"
//	"hcc/harp/lib/configext"
//	"hcc/harp/lib/iputil"
//	"hcc/harp/lib/logger"
//	"hcc/harp/lib/mysql"
//	"hcc/harp/lib/pf"
//	"hcc/harp/model"
//	"time"
//)
//
//// ReadAdaptiveIPServer - ish
//func ReadAdaptiveIPServer(args map[string]interface{}) (interface{}, error) {
//	var adaptiveipServer model.AdaptiveIPServer
//
//	serverUUID := args["server_uuid"].(string)
//	var publicIP string
//	var privateIP string
//	var privateGateway string
//	var createdAt time.Time
//
//	sql := "select server_uuid, public_ip, private_ip, private_gateway, created_at from adaptiveip_server where server_uuid = ?"
//	err := mysql.Db.QueryRow(sql, serverUUID).Scan(
//		&serverUUID,
//		&publicIP,
//		&privateIP,
//		&privateGateway,
//		&createdAt)
//	if err != nil {
//		logger.Logger.Println(err)
//		return nil, err
//	}
//
//	adaptiveipServer.PublicIP = publicIP
//	adaptiveipServer.PrivateIP = privateIP
//	adaptiveipServer.PrivateGateway = privateGateway
//	adaptiveipServer.CreatedAt = createdAt
//
//	return adaptiveipServer, nil
//}
//
//func checkReadAdaptiveIPServerListPageRow(args map[string]interface{}) bool {
//	_, rowOk := args["row"].(int)
//	_, pageOk := args["page"].(int)
//
//	return !rowOk || !pageOk
//}
//
//// ReadAdaptiveIPServerList - ish
//func ReadAdaptiveIPServerList(args map[string]interface{}) (interface{}, error) {
//	var adaptiveipServers []model.AdaptiveIPServer
//
//	serverUUID, serverUUIDOk := args["server_uuid"].(string)
//	publicIP, publicIPOk := args["public_ip"].(string)
//	privateIP, privateIPOk := args["private_ip"].(string)
//	privateGateway, privateGatewayOk := args["private_gateway"].(string)
//	var createdAt time.Time
//
//	row, _ := args["row"].(int)
//	page, _ := args["page"].(int)
//	if checkReadAdaptiveIPServerListPageRow(args) {
//		return nil, errors.New("need row and page arguments")
//	}
//
//	sql := "select server_uuid, public_ip, private_ip, private_gateway, created_at from adaptiveip_server where 1=1"
//
//	if serverUUIDOk {
//		sql += " and server_uuid = '" + serverUUID + "'"
//	}
//	if publicIPOk {
//		sql += " and public_ip = '" + publicIP + "'"
//	}
//	if privateIPOk {
//		sql += " and private_ip = '" + privateIP + "'"
//	}
//	if privateGatewayOk {
//		sql += " and private_gateway = '" + privateGateway + "'"
//	}
//
//	sql += " order by created_at desc limit ? offset ?"
//
//	stmt, err := mysql.Db.Query(sql, row, row*(page-1))
//	if err != nil {
//		logger.Logger.Println(err.Error())
//		return nil, err
//	}
//	defer func() {
//		_ = stmt.Close()
//	}()
//
//	for stmt.Next() {
//		err := stmt.Scan(&serverUUID, &publicIP, &privateIP, &privateGateway, &createdAt)
//		if err != nil {
//			logger.Logger.Println(err.Error())
//			return nil, err
//		}
//		adaptiveipServer := model.AdaptiveIPServer{ServerUUID: serverUUID, PublicIP: publicIP, PrivateIP: privateIP,
//			PrivateGateway: privateGateway, CreatedAt: createdAt}
//		adaptiveipServers = append(adaptiveipServers, adaptiveipServer)
//	}
//
//	adaptiveIP := configext.GetAdaptiveIPNetwork()
//	netNetwork, err := iputil.CheckNetwork(adaptiveIP.ExtIfaceIPAddress, adaptiveIP.Netmask)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, adaptiveIPServer := range adaptiveipServers {
//		netIP := iputil.CheckValidIP(adaptiveIPServer.PublicIP)
//		if !netNetwork.Contains(netIP) {
//			adaptiveIPServer.Status = "Invalid"
//			continue
//		}
//		adaptiveIPServer.Status = "Using"
//	}
//
//	return adaptiveipServers, nil
//}
//
//// ReadAdaptiveIPServerAll - ish
//func ReadAdaptiveIPServerAll(args map[string]interface{}) (interface{}, error) {
//	var adaptiveipServers []model.AdaptiveIPServer
//	var serverUUID string
//	var publicIP string
//	var privateIP string
//	var privateGateway string
//	var createdAt time.Time
//
//	row, rowOk := args["row"].(int)
//	page, pageOk := args["page"].(int)
//	var sql string
//	var stmt *dbsql.Rows
//	var err error
//
//	if !rowOk && !pageOk {
//		sql = "select server_uuid, public_ip, private_ip, private_gateway, created_at from adaptiveip_server"
//		stmt, err = mysql.Db.Query(sql)
//	} else if rowOk && pageOk {
//		sql = "select server_uuid, public_ip, private_ip, private_gateway, created_at from adaptiveip_server order by created_at desc limit ? offset ?"
//		stmt, err = mysql.Db.Query(sql, row, row*(page-1))
//	} else {
//		return nil, errors.New("please insert row and page arguments or leave arguments as empty state")
//	}
//
//	if err != nil {
//		logger.Logger.Println(err.Error())
//		return nil, err
//	}
//	defer func() {
//		_ = stmt.Close()
//	}()
//
//	for stmt.Next() {
//		err := stmt.Scan(&serverUUID, &publicIP, &privateIP, &privateGateway, &createdAt)
//		if err != nil {
//			logger.Logger.Println(err)
//			return nil, err
//		}
//		adaptiveipServer := model.AdaptiveIPServer{ServerUUID: serverUUID, PublicIP: publicIP, PrivateIP: privateIP,
//			PrivateGateway: privateGateway, CreatedAt: createdAt}
//		adaptiveipServers = append(adaptiveipServers, adaptiveipServer)
//	}
//
//	adaptiveIP := configext.GetAdaptiveIPNetwork()
//	netNetwork, err := iputil.CheckNetwork(adaptiveIP.ExtIfaceIPAddress, adaptiveIP.Netmask)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, adaptiveIPServer := range adaptiveipServers {
//		netIP := iputil.CheckValidIP(adaptiveIPServer.PublicIP)
//		if !netNetwork.Contains(netIP) {
//			adaptiveIPServer.Status = "Invalid"
//			continue
//		}
//		adaptiveIPServer.Status = "Using"
//	}
//
//	return adaptiveipServers, nil
//}
//
//// ReadAdaptiveIPServerNum - ish
//func ReadAdaptiveIPServerNum(args map[string]interface{}) (model.AdaptiveIPServerNum, error) {
//	var adaptiveIPServerNum model.AdaptiveIPServerNum
//
//	serverUUID, serverUUIDOk := args["server_uuid"].(string)
//	if !serverUUIDOk {
//		return adaptiveIPServerNum, errors.New("need a server_uuid argument")
//	}
//
//	var adaptiveIPServerNr int
//	var err error
//
//	sql := "select count(*) from adaptiveip_server where server_uuid = '" + serverUUID + "'"
//	err = mysql.Db.QueryRow(sql).Scan(&adaptiveIPServerNr)
//	if err != nil {
//		logger.Logger.Println(err)
//		return adaptiveIPServerNum, err
//	}
//	adaptiveIPServerNum.Number = adaptiveIPServerNr
//
//	return adaptiveIPServerNum, nil
//}
//
//// CreateAdaptiveIPServer - ish
//func CreateAdaptiveIPServer(args map[string]interface{}) (interface{}, error) {
//	adaptiveipServer := model.AdaptiveIPServer{
//		ServerUUID: args["server_uuid"].(string),
//		PublicIP:   args["public_ip"].(string),
//	}
//
//	subnet, err := ReadSubnetByServer(adaptiveipServer.ServerUUID)
//	if err != nil {
//		return nil, errors.New("provided server_uuid is not allocated to one of private subnet")
//	}
//
//	firstIP, _, err := iputil.GetFirstAndLastIPs(subnet.(model.Subnet).NetworkIP, subnet.(model.Subnet).Netmask)
//	if err != nil {
//		return nil, err
//	}
//
//	adaptiveipServer.PrivateIP = firstIP.String()
//	adaptiveipServer.PrivateGateway = subnet.(model.Subnet).Gateway
//
//	err = pf.CreateAndLoadAnchorConfig(adaptiveipServer.PublicIP, adaptiveipServer.PrivateIP)
//	if err != nil {
//		return nil, err
//	}
//
//	sql := "insert into adaptiveip_server(server_uuid, public_ip, private_ip, private_gateway, created_at) values (?, ?, ?, ?, now())"
//	stmt, err := mysql.Db.Prepare(sql)
//	if err != nil {
//		logger.Logger.Println(err.Error())
//		return nil, err
//	}
//	defer func() {
//		_ = stmt.Close()
//	}()
//	result, err := stmt.Exec(adaptiveipServer.ServerUUID, adaptiveipServer.PublicIP, adaptiveipServer.PrivateIP,
//		adaptiveipServer.PrivateGateway)
//	if err != nil {
//		logger.Logger.Println(err)
//		return nil, err
//	}
//	logger.Logger.Println(result.LastInsertId())
//
//	return adaptiveipServer, nil
//}
//
//// DeleteAdaptiveIPServer - ish
//func DeleteAdaptiveIPServer(args map[string]interface{}) (interface{}, error) {
//	var err error
//
//	serverUUID, serverUUIDOk := args["server_uuid"].(string)
//	if !serverUUIDOk {
//		return nil, errors.New("need a server_uuid argument")
//	}
//
//	serverUUIDArg := make(map[string]interface{})
//	serverUUIDArg["server_uuid"] = serverUUID
//	adaptiveipServer, err := ReadAdaptiveIPServer(serverUUIDArg)
//
//	err = pf.DeleteAndUnloadAnchorConfig(adaptiveipServer.(model.AdaptiveIPServer).PublicIP)
//	if err != nil {
//		logger.Logger.Println(err.Error())
//		return nil, err
//	}
//
//	sql := "delete from adaptiveip_server where server_uuid = ?"
//	stmt, err := mysql.Db.Prepare(sql)
//	if err != nil {
//		logger.Logger.Println(err.Error())
//		return nil, err
//	}
//	defer func() {
//		_ = stmt.Close()
//	}()
//	result, err2 := stmt.Exec(serverUUID)
//	if err2 != nil {
//		logger.Logger.Println(err2)
//		return nil, err
//	}
//	logger.Logger.Println(result.RowsAffected())
//
//	return serverUUID, nil
//}
