package dao

import (
	"hcc/harp/action/grpc/client"
	daoext2 "hcc/harp/daoext"
	"hcc/harp/lib/configadapriveipnetwork"
	"hcc/harp/lib/hccweb"
	"hcc/harp/lib/iptablesext"
	"hcc/harp/lib/iputil"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"strings"
)

// ReadAdaptiveIPServerNum : Get the number of AdaptiveIPServer
func ReadAdaptiveIPServerNum(in *pb.ReqGetAdaptiveIPServerNum) (*pb.ResGetAdaptiveIPServerNum, uint64, string) {
	var adaptiveIPServerNum pb.ResGetAdaptiveIPServerNum
	var adaptiveIPServerNr int64

	var groupID = in.GetGroupID()
	if groupID == 0 {
		return nil, hcc_errors.HarpGrpcArgumentError, "ReadAdaptiveIPServerNum(): please insert a group_id argument"
	}

	sql := "select count(*) from adaptiveip_server where group_id = " + strconv.Itoa(int(groupID))
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
	publicIP := in.PublicIP
	publicIPOk := len(publicIP) != 0

	if !serverUUIDOk || !publicIPOk {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreateAdaptiveIPServer(): need ServerUUID and PublicIP arguments"
	}

	oldAdaptiveIPServer, _, _ := daoext2.ReadAdaptiveIPServer(serverUUID)
	if oldAdaptiveIPServer != nil {
		return nil, hcc_errors.HarpInternalAdaptiveIPAllocatedError, "CreateAdaptiveIPServer(): provided ServerUUID is already allocated to one of adaptiveIP"
	}

	subnet, errCode, _ := daoext2.ReadSubnetByServer(serverUUID)
	if errCode != 0 {
		return nil, hcc_errors.HarpInternalSubnetNotAllocatedError, "CreateAdaptiveIPServer(): provided ServerUUID is not allocated to one of private subnet"
	}

	server, err := client.RC.GetServer(serverUUID)
	if err != nil {
		return nil, hcc_errors.HarpGrpcRequestError, "CreateAdaptiveIPServer(): Failed to get server information"
	}
	if server.Server != nil &&
		(strings.ToLower(server.Server.Status) == "creating" || strings.ToLower(server.Server.Status) == "deleting") {
		return nil, hcc_errors.HarpInternalOperationFail, "CreateAdaptiveIPServer(): You can't create the AdaptiveIP when creating or deleting the server"
	}

	resGetQuota, errStack := client.RC.GetQuota(subnet.GroupID)
	if errStack != nil {
		return nil, hcc_errors.HarpGrpcRequestError, "CreateAdaptiveIPServer(): " + errStack.Pop().Text()
	}
	adaptiveIPServerNum, errCode, errText := ReadAdaptiveIPServerNum(&pb.ReqGetAdaptiveIPServerNum{GroupID: subnet.GroupID})
	if errCode != 0 {
		return nil, errCode, "CreateAdaptiveIPServer(): " + errText
	}
	if adaptiveIPServerNum.Num+1 > int64(resGetQuota.Quota.LimitAdaptiveIPCnt) {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreateAdaptiveIPServer(): AdaptiveIP count quota exceeded"
	}

	adaptiveIP := configadapriveipnetwork.GetAdaptiveIPNetwork()

	var startIPSum = 0
	var endIPSsum = 0
	var publicIPSum = 0

	startIPSplit := strings.Split(adaptiveIP.ExternalStartIPAddress, ".")
	endIPSplit := strings.Split(adaptiveIP.ExternalEndIPAddress, ".")
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
		GroupID:    subnet.GroupID,
		PublicIP:   publicIP,
	}

	firstIP, _, err := iputil.GetFirstAndLastIPs(subnet.NetworkIP, subnet.Netmask)
	if err != nil {
		return nil, hcc_errors.HarpInternalIPAddressError, "CreateAdaptiveIPServer(): " + err.Error()
	}

	adaptiveIPServer.PrivateIP = firstIP.String()
	adaptiveIPServer.PrivateGateway = subnet.Gateway

	err = iptablesext.ControlNetDevAndIPTABLES(true, adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP)
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

	err = hccweb.PortForwardDocker(true, serverUUID)
	if err != nil {
		errStr := "CreateAdaptiveIPServer(): Failed to create port forwarding for the hccweb docker container"
		logger.Logger.Println(errStr)
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

	adaptiveIPServer, errCode, _ := daoext2.ReadAdaptiveIPServer(serverUUID)
	if errCode == hcc_errors.HarpSQLNoResult {
		return serverUUID, 0, ""
	}
	if adaptiveIPServer == nil {
		return "", hcc_errors.HarpGrpcArgumentError, "DeleteAdaptiveIPServer(): adaptiveIPServer is nil"
	}

	server, err := client.RC.GetServer(serverUUID)
	if err != nil {
		return "", hcc_errors.HarpGrpcRequestError, "DeleteAdaptiveIPServer(): Failed to get server information"
	}
	if server.Server != nil &&
		(strings.ToLower(server.Server.Status) == "creating" || strings.ToLower(server.Server.Status) == "deleting") {
		return "", hcc_errors.HarpInternalOperationFail, "DeleteAdaptiveIPServer(): You can't delete the AdaptiveIP when creating or deleting the server"
	}

	err = hccweb.PortForwardDocker(false, serverUUID)
	if err != nil {
		errStr := "DeleteAdaptiveIPServer(): Failed to delete port forwarding for the hccweb docker container"
		logger.Logger.Println(errStr)
	}

	portForwardingList, errCode, errStr := ReadPortForwardingList(&pb.ReqGetPortForwardingList{
		PortForwarding: &pb.PortForwarding{
			ServerUUID: serverUUID,
		},
	})
	if errCode != 0 {
		return "", hcc_errors.HarpInternalOperationFail, "DeleteAdaptiveIPServer(): " + errStr
	}

	for _, portForward := range portForwardingList.PortForwarding {
		_, errCode, errStr := DeletePortForwarding(&pb.ReqDeletePortForwarding{
			PortForwarding: &pb.PortForwarding{
				ServerUUID:   serverUUID,
				ExternalPort: portForward.ExternalPort,
			},
		})
		if errCode != 0 {
			return "", hcc_errors.HarpInternalOperationFail, "DeleteAdaptiveIPServer(): " + errStr
		}
	}

	err = iptablesext.ControlNetDevAndIPTABLES(false, adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP)
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
	_, err = stmt.Exec(serverUUID)
	if err != nil {
		errStr := "DeleteAdaptiveIPServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hcc_errors.HarpSQLOperationFail, errStr
	}

	return serverUUID, 0, ""
}
