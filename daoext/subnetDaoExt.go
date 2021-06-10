package daoext

import (
	"github.com/golang/protobuf/ptypes"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"strings"
	"time"
)

// ReadSubnetByServer : Get infos of a subnet by server UUID
func ReadSubnetByServer(serverUUID string) (*pb.Subnet, uint64, string) {
	var subnet pb.Subnet

	var uuid string
	var groupID int64
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

	sql := "select uuid, group_id, network_ip, netmask, gateway, next_server, name_server, domain_name, leader_node_uuid, os, subnet_name, created_at from subnet where server_uuid = ?"
	row := mysql.Db.QueryRow(sql, serverUUID)
	err := mysql.QueryRowScan(row,
		&uuid,
		&groupID,
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
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hcc_errors.HarpSQLNoResult, errStr
		}
		logger.Logger.Println(errStr)
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

	subnet.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		errStr := "ReadSubnetByServer(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.HarpInternalTimeStampConversionError, errStr
	}

	return &subnet, 0, ""
}
