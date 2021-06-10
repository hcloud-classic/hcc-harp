package dao

import (
	dbsql "database/sql"
	"errors"
	daoext2 "hcc/harp/daoext"
	"hcc/harp/lib/iptablesext"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"strings"
)

// ReadPortForwardingList : Get the AdaptiveIP's PortForwarding list of the server
func ReadPortForwardingList(in *pb.ReqGetPortForwardingList) (*pb.ResGetPortForwardingList, uint64, string) {
	var portForwardingList pb.ResGetPortForwardingList
	var portForwardings []pb.PortForwarding
	var pPortForwardings []*pb.PortForwarding

	if in.PortForwarding == nil {
		return nil, hcc_errors.HarpGrpcArgumentError, "ReadPortForwardingList(): PortForwarding is nil"
	}

	var serverUUID string
	var forwardTCP bool
	var forwardUDP bool
	var protocol string
	var externalPort int64
	var internalPort int64
	var description string

	reqPortForwarding := in.PortForwarding

	serverUUID = reqPortForwarding.ServerUUID
	serverUUIDOk := len(serverUUID) != 0

	if !serverUUIDOk {
		return nil, hcc_errors.HarpGrpcArgumentError, "ReadPortForwardingList(): Need a serverUUID argument"
	}

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
		return nil, hcc_errors.HarpGrpcArgumentError, "ReadPortForwardingList(): please insert row and page arguments or leave arguments as empty state"
	}

	sql := "select * from port_forwarding where server_uuid = '" + serverUUID + "'"

	forwardTCP = reqPortForwarding.ForwardTCP
	forwardTCPOk := forwardTCP
	forwardUDP = reqPortForwarding.ForwardUDP
	forwardUDPOk := forwardUDP
	externalPort = reqPortForwarding.ExternalPort
	externalPortOk := externalPort != 0
	internalPort = reqPortForwarding.InternalPort
	internalPortOk := internalPort != 0
	description = reqPortForwarding.Description
	descriptionOk := len(description) != 0

	if forwardTCPOk || forwardUDPOk {
		if forwardTCPOk && forwardUDPOk {
			protocol = "all"
		} else if forwardTCPOk {
			protocol = "tcp"
		} else if forwardUDPOk {
			protocol = "udp"
		}

		sql += " and protocol = '" + protocol + "'"
	}

	if externalPortOk {
		sql += " and external_port = " + strconv.Itoa(int(externalPort))
	}
	if internalPortOk {
		sql += " and internal_port = " + strconv.Itoa(int(internalPort))
	}
	if descriptionOk {
		sql += " and description like '%" + description + "%'"
	}

	var stmt *dbsql.Rows
	var err error
	if isLimit {
		sql += " order by external_port asc limit ? offset ?"
		stmt, err = mysql.Query(sql, row, row*(page-1))
	} else {
		sql += " order by external_port asc"
		stmt, err = mysql.Query(sql)
	}

	if err != nil {
		errStr := "ReadPortForwardingList(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&serverUUID, &protocol, &externalPort, &internalPort, &description)
		if err != nil {
			errStr := "ReadPortForwardingList(): " + err.Error()
			logger.Logger.Println(errStr)
			if strings.Contains(err.Error(), "no rows in result set") {
				return nil, hcc_errors.HarpSQLNoResult, errStr
			}
			return nil, hcc_errors.HarpSQLOperationFail, errStr
		}

		if protocol == "all" {
			forwardTCP = true
			forwardUDP = true
		} else if protocol == "tcp" {
			forwardTCP = true
			forwardUDP = false
		} else if protocol == "udp" {
			forwardTCP = false
			forwardUDP = true
		} else {
			return nil, hcc_errors.HarpSQLOperationFail, "ReadPortForwardingList(): Unknown protocol (serverUUID=" + serverUUID + ")"
		}

		portForwardings = append(portForwardings, pb.PortForwarding{
			ServerUUID:   serverUUID,
			ForwardTCP:   forwardTCP,
			ForwardUDP:   forwardUDP,
			ExternalPort: externalPort,
			InternalPort: internalPort,
			Description:  description,
		})

	}

	for i := range portForwardings {
		pPortForwardings = append(pPortForwardings, &portForwardings[i])
	}

	portForwardingList.PortForwarding = pPortForwardings

	return &portForwardingList, 0, ""
}

func checkPortRange(port int64) error {
	if port < 0 || port > 65535 {
		return errors.New("port number is out of range")
	}

	return nil
}

// CreatePortForwarding : Create AdaptiveIP's PortForwarding of server
func CreatePortForwarding(in *pb.ReqCreatePortForwarding) (*pb.PortForwarding, uint64, string) {
	if in.PortForwarding == nil {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreatePortForwarding(): PortForwarding is nil"
	}

	var serverUUID string
	var forwardTCP bool
	var forwardUDP bool
	var protocol string
	var externalPort int64
	var internalPort int64
	var description string

	reqPortForwarding := in.PortForwarding

	serverUUID = reqPortForwarding.ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	forwardTCP = reqPortForwarding.ForwardTCP
	forwardTCPOk := forwardTCP
	forwardUDP = reqPortForwarding.ForwardUDP
	forwardUDPOk := forwardUDP
	externalPort = reqPortForwarding.ExternalPort
	externalPortOk := externalPort != 0
	internalPort = reqPortForwarding.InternalPort
	internalPortOk := internalPort != 0
	description = reqPortForwarding.Description
	descriptionOk := len(description) != 0

	if !serverUUIDOk || (!forwardTCPOk && !forwardUDPOk) || !externalPortOk || !internalPortOk || !descriptionOk {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreatePortForwarding(): need ServerUUID and ForwardTCP/ForwardUDP," +
			"ExternalPort, InternalPort, Description arguments"
	}

	adaptiveIPServer, _, _ := daoext2.ReadAdaptiveIPServer(serverUUID)
	if adaptiveIPServer == nil {
		return nil, hcc_errors.HarpInternalAdaptiveIPAllocatedError,
			"CreatePortForwarding(): AdaptiveIP is not allocated for the server (serverUUID=" + serverUUID + ")"
	}

	err := checkPortRange(externalPort)
	if err != nil {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreatePortForwarding(): External " + err.Error()
	}

	err = checkPortRange(internalPort)
	if err != nil {
		return nil, hcc_errors.HarpGrpcArgumentError, "CreatePortForwarding(): Internal " + err.Error()
	}

	oldPortForwarding, _, _ := ReadPortForwardingList(&pb.ReqGetPortForwardingList{
		PortForwarding: &pb.PortForwarding{
			ServerUUID: serverUUID,
		},
	})
	if oldPortForwarding != nil {
		for _, portForwarding := range oldPortForwarding.PortForwarding {
			if portForwarding.ExternalPort == externalPort {
				return nil, hcc_errors.HarpInternalAdaptiveIPAllocatedError, "CreatePortForwarding(): Provided external port is already allocated to the server"
			}
		}
	}

	err = iptablesext.PortForwarding(true, forwardTCP, forwardUDP,
		adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP, int(externalPort), int(internalPort))
	if err != nil {
		return nil, hcc_errors.HarpInternalAdaptiveIPAllocatedError, "CreatePortForwarding(): " + err.Error()
	}

	if forwardTCPOk || forwardUDPOk {
		if forwardTCPOk && forwardUDPOk {
			protocol = "all"
		} else if forwardTCPOk {
			protocol = "tcp"
		} else if forwardUDPOk {
			protocol = "udp"
		}
	}

	portForwarding := pb.PortForwarding{
		ServerUUID:   serverUUID,
		ForwardTCP:   forwardTCP,
		ForwardUDP:   forwardUDP,
		ExternalPort: externalPort,
		InternalPort: internalPort,
		Description:  description,
	}

	sql := "insert into port_forwarding(server_uuid, protocol, external_port, internal_port, description) values (?, ?, ?, ?, ?)"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "CreatePortForwarding(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(portForwarding.ServerUUID, protocol, portForwarding.ExternalPort, portForwarding.InternalPort, portForwarding.Description)
	if err != nil {
		errStr := "CreatePortForwarding(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}

	return &portForwarding, 0, ""
}

// DeletePortForwarding : Delete AdaptiveIP's PortForwarding of the server
func DeletePortForwarding(in *pb.ReqDeletePortForwarding) (string, uint64, string) {
	if in.PortForwarding == nil {
		return "", hcc_errors.HarpGrpcArgumentError, "DeletePortForwarding(): PortForwarding is nil"
	}

	reqPortForwarding := in.PortForwarding

	serverUUID := reqPortForwarding.ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	externalPort := reqPortForwarding.ExternalPort
	externalPortOk := externalPort != 0
	if !serverUUIDOk || !externalPortOk {
		return "", hcc_errors.HarpGrpcArgumentError, "DeletePortForwarding(): need ServerUUID and ExternalPort arguments"
	}

	adaptiveIPServer, _, _ := daoext2.ReadAdaptiveIPServer(serverUUID)
	if adaptiveIPServer == nil {
		return "", hcc_errors.HarpInternalAdaptiveIPAllocatedError,
			"DeletePortForwarding(): AdaptiveIP is not found with the provided ServerUUID"
	}

	err := checkPortRange(externalPort)
	if err != nil {
		return "", hcc_errors.HarpGrpcArgumentError, "DeletePortForwarding(): External " + err.Error()
	}

	portForwarding, _, _ := ReadPortForwardingList(&pb.ReqGetPortForwardingList{
		PortForwarding: &pb.PortForwarding{
			ServerUUID:   serverUUID,
			ExternalPort: externalPort,
		},
	})
	if portForwarding == nil || len(portForwarding.PortForwarding) != 1 {
		return "", hcc_errors.HarpInternalAdaptiveIPAllocatedError,
			"DeletePortForwarding(): Failed to get port forwarding info" +
				"(serverUUID=" + serverUUID + ", ExternalPort=" + strconv.Itoa(int(externalPort)) + ")"
	}

	err = iptablesext.PortForwarding(false, portForwarding.PortForwarding[0].ForwardTCP, portForwarding.PortForwarding[0].ForwardUDP,
		adaptiveIPServer.PublicIP, adaptiveIPServer.PrivateIP, int(externalPort), int(portForwarding.PortForwarding[0].InternalPort))
	if err != nil {
		return "", hcc_errors.HarpInternalAdaptiveIPAllocatedError, "DeletePortForwarding(): " + err.Error()
	}

	sql := "delete from port_forwarding where server_uuid = ? and external_port = ?"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "DeletePortForwarding(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hcc_errors.HarpSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(serverUUID, externalPort)
	if err != nil {
		errStr := "DeletePortForwarding(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hcc_errors.HarpSQLOperationFail, errStr
	}

	return serverUUID, 0, ""
}
