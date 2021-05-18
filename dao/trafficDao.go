package dao

import (
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"strings"
)

// GetTraffic : Get the traffic info from the database
func GetTraffic(serverUUID string, day string) (*pb.Traffic, uint64, string) {
	var traffic pb.Traffic

	var groupID int64
	var txKB int64
	var rxKB int64

	sql := "select group_id, tx_kb, rx_kb from traffic where server_uuid = ? and day = ?"
	row := mysql.Db.QueryRow(sql, serverUUID, day)
	err := mysql.QueryRowScan(row,
		&groupID,
		&txKB,
		&rxKB)
	if err != nil {
		errStr := "GetTraffic(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hcc_errors.HarpSQLNoResult, errStr
		}
		return nil, hcc_errors.HarpSQLOperationFail, errStr
	}

	traffic.ServerUUID = serverUUID
	traffic.GroupID = groupID
	traffic.Day = day
	traffic.TxKB = txKB
	traffic.RxKB = rxKB

	return &traffic, 0, ""
}
