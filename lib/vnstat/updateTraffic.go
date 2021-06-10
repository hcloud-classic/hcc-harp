package vnstat

import (
	"errors"
	"hcc/harp/daoext"
	"hcc/harp/lib/iplinkext"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"time"
)

func getTodayByNumString() string {
	currentTime := time.Now()
	return currentTime.Format("060102")
}

func checkIfTodayTrafficExist(serverUUID string) bool {
	sql := "select server_uuid from traffic where server_uuid = ? and day = ?"
	row := mysql.Db.QueryRow(sql, serverUUID, getTodayByNumString())
	err := mysql.QueryRowScan(row, &serverUUID)
	if err != nil {
		return false
	}

	return true
}

func insertTodayTraffic(serverUUID string) error {
	subnet, errCode, errText := daoext.ReadSubnetByServer(serverUUID)
	if errCode != 0 {
		return errors.New(errText)
	}

	sql := "insert into traffic(server_uuid, group_id, tx_kb, rx_kb, day) values (?, ?, ?, ?, ?)"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "InsertTodayTraffic(): " + err.Error()
		logger.Logger.Println(errStr)
		return errors.New(errStr)
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(serverUUID, subnet.GroupID, 0, 0, getTodayByNumString())
	if err != nil {
		errStr := "InsertTodayTraffic(): " + err.Error()
		logger.Logger.Println(errStr)
		return errors.New(errStr)
	}

	return nil
}

func updateTodayTraffic(harpIface string) error {
	adaptiveIPServerList, errCode, errText := daoext.ReadAdaptiveIPServerList(&pb.ReqGetAdaptiveIPServerList{})
	if errCode != 0 {
		return errors.New(errText)
	}

	var serverUUID = ""
	for _, adaptiveIPServer := range adaptiveIPServerList.AdaptiveipServer {
		subnet, _, _ := daoext.ReadSubnetByServer(adaptiveIPServer.ServerUUID)
		if subnet == nil {
			continue
		}
		ifaceName := iplinkext.HarpInternalPrefix + strconv.Itoa(iplinkext.GetIfaceVNUM(subnet.Gateway))
		if harpIface == ifaceName {
			serverUUID = adaptiveIPServer.ServerUUID
			break
		}
	}

	if serverUUID == "" {
		return nil
	}

	if !checkIfTodayTrafficExist(serverUUID) {
		err := insertTodayTraffic(serverUUID)
		if err != nil {
			return err
		}
	}

	txKB, rxKB, err := GetTodayVnStatData(harpIface)
	if err != nil {
		return err
	}

	sql := "update traffic set tx_kb = ?, rx_kb = ? where server_uuid = ? and day = ?"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "updateTodayTraffic(): " + err.Error()
		logger.Logger.Println(errStr)
		return errors.New(errStr)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(txKB, rxKB, serverUUID, getTodayByNumString())
	if err != nil {
		errStr := "updateTodayTraffic(): " + err.Error()
		logger.Logger.Println(errStr)
		return errors.New(errStr)
	}

	return nil
}
