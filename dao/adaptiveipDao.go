package dao

import (
	dbsql "database/sql"
	"errors"
	gouuid "github.com/nu7hatch/gouuid"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/model"
	"time"
)

// ReadAdaptiveIP - ish
func ReadAdaptiveIP(args map[string]interface{}) (interface{}, error) {
	var adaptiveip model.AdaptiveIP

	uuid := args["uuid"].(string)
	var networkAddress string
	var netmask string
	var gateway string
	var startIPaddress string
	var endIPaddress string
	var createdAt time.Time

	sql := "select network_address, netmask, gateway, start_ip_address, end_ip_address, created_at from adaptiveip where uuid = ?"
	err := mysql.Db.QueryRow(sql, uuid).Scan(
		&networkAddress,
		&netmask,
		&gateway,
		&startIPaddress,
		&endIPaddress,
		&createdAt)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	adaptiveip.UUID = uuid
	adaptiveip.NetworkAddress = networkAddress
	adaptiveip.Netmask = netmask
	adaptiveip.Gateway = gateway
	adaptiveip.StartIPAddress = startIPaddress
	adaptiveip.EndIPAddress = endIPaddress
	adaptiveip.CreatedAt = createdAt

	return adaptiveip, nil
}

func checkReadAdaptiveIPListPageRow(args map[string]interface{}) bool {
	_, rowOk := args["row"].(int)
	_, pageOk := args["page"].(int)

	return !rowOk || !pageOk
}

// ReadAdaptiveIPList - ish
func ReadAdaptiveIPList(args map[string]interface{}) (interface{}, error) {
	var adaptiveips []model.AdaptiveIP
	var uuid string
	var createdAt time.Time

	networkAddress, networkAddressOk := args["network_address"].(string)
	netmask, netmaskOk := args["netmask"].(string)
	gateway, gatewayOk := args["gateway"].(string)
	startIPaddress, startIPaddressOk := args["start_ip_address"].(string)
	endIPaddress, endIPaddressOk := args["end_ip_address"].(string)

	row, _ := args["row"].(int)
	page, _ := args["page"].(int)
	if checkReadAdaptiveIPListPageRow(args) {
		return nil, errors.New("need row and page arguments")
	}

	sql := "select * from adaptiveip where 1=1"

	if networkAddressOk {
		sql += " and network_address = '" + networkAddress + "'"
	}
	if netmaskOk {
		sql += " and netmask = '" + netmask + "'"
	}
	if gatewayOk {
		sql += " and gateway = '" + gateway + "'"
	}
	if startIPaddressOk {
		sql += " and start_ip_address = '" + startIPaddress + "'"
	}
	if endIPaddressOk {
		sql += " and end_ip_address = '" + endIPaddress + "'"
	}

	sql += " order by created_at desc limit ? offset ?"

	stmt, err := mysql.Db.Query(sql, row, row*(page-1))
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &networkAddress, &netmask, &gateway, &startIPaddress, &endIPaddress, &createdAt)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		adaptiveip := model.AdaptiveIP{UUID: uuid, NetworkAddress: networkAddress, Netmask: netmask, Gateway: gateway, StartIPAddress: startIPaddress, EndIPAddress: endIPaddress, CreatedAt: createdAt}
		adaptiveips = append(adaptiveips, adaptiveip)
	}
	return adaptiveips, nil
}

// ReadAdaptiveIPAll - ish
func ReadAdaptiveIPAll(args map[string]interface{}) (interface{}, error) {
	var adaptiveips []model.AdaptiveIP
	var uuid string
	var networkAddress string
	var netmask string
	var gateway string
	var startIPaddress string
	var endIPaddress string
	var createdAt time.Time

	row, rowOk := args["row"].(int)
	page, pageOk := args["page"].(int)
	var sql string
	var stmt *dbsql.Rows
	var err error

	if !rowOk && !pageOk {
		sql = "select * from adaptiveip order by created_at desc"
		stmt, err = mysql.Db.Query(sql)
	} else if rowOk && pageOk {
		sql = "select * from adaptiveip order by created_at desc limit ? offset ?"
		stmt, err = mysql.Db.Query(sql, row, row*(page-1))
	} else {
		return nil, errors.New("please insert row and page arguments or leave arguments as empty state")
	}

	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &networkAddress, &netmask, &gateway, &startIPaddress, &endIPaddress, &createdAt)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		adaptiveip := model.AdaptiveIP{UUID: uuid, NetworkAddress: networkAddress, Netmask: netmask, Gateway: gateway, StartIPAddress: startIPaddress, EndIPAddress: endIPaddress, CreatedAt: createdAt}
		adaptiveips = append(adaptiveips, adaptiveip)
	}

	return adaptiveips, nil
}

// ReadAdaptiveIPNum - ish
func ReadAdaptiveIPNum() (model.AdaptiveIPNum, error) {
	var adaptiveIPNum model.AdaptiveIPNum
	var adaptiveIPNr int
	var err error

	sql := "select count(*) from adaptiveip"
	err = mysql.Db.QueryRow(sql).Scan(&adaptiveIPNr)
	if err != nil {
		logger.Logger.Println(err)
		return adaptiveIPNum, err
	}
	adaptiveIPNum.Number = adaptiveIPNr

	return adaptiveIPNum, nil
}

// CreateAdaptiveIP - ish
func CreateAdaptiveIP(args map[string]interface{}) (interface{}, error) {
	out, err := gouuid.NewV4()
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	uuid := out.String()

	adaptiveip := model.AdaptiveIP{
		UUID:           uuid,
		NetworkAddress: args["network_address"].(string),
		Netmask:        args["netmask"].(string),
		Gateway:        args["gateway"].(string),
		StartIPAddress: args["start_ip_address"].(string),
		EndIPAddress:   args["end_ip_address"].(string),
	}

	sql := "insert into adaptiveip(uuid, network_address, netmask, gateway, start_ip_address, end_ip_address, created_at) values (?, ?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err := stmt.Exec(adaptiveip.UUID, adaptiveip.NetworkAddress, adaptiveip.Netmask, adaptiveip.Gateway, adaptiveip.StartIPAddress, adaptiveip.EndIPAddress)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	logger.Logger.Println(result.LastInsertId())

	return adaptiveip, nil
}

func checkUpdateAdaptiveIPArgs(args map[string]interface{}) bool {
	_, networkAddressOk := args["network_address"].(string)
	_, netmaskOk := args["netmask"].(string)
	_, gatewayOk := args["gateway"].(string)
	_, startIPaddressOk := args["start_ip_address"].(string)
	_, endIPaddressOk := args["end_ip_address"].(string)

	return !networkAddressOk && !netmaskOk && !gatewayOk && !startIPaddressOk && !endIPaddressOk
}

// UpdateAdaptiveIP - ish
func UpdateAdaptiveIP(args map[string]interface{}) (interface{}, error) {
	requestedUUID, requestedUUIDOk := args["uuid"].(string)
	networkAddress, networkAddressOk := args["network_address"].(string)
	netmask, netmaskOk := args["netmask"].(string)
	gateway, gatewayOk := args["gateway"].(string)
	startIPaddress, startIPaddressOk := args["start_ip_address"].(string)
	endIPaddress, endIPaddressOk := args["end_ip_address"].(string)

	adaptiveip := new(model.AdaptiveIP)
	adaptiveip.UUID = requestedUUID
	adaptiveip.NetworkAddress = networkAddress
	adaptiveip.Netmask = netmask
	adaptiveip.Gateway = gateway
	adaptiveip.StartIPAddress = startIPaddress
	adaptiveip.EndIPAddress = endIPaddress

	if requestedUUIDOk {
		if checkUpdateAdaptiveIPArgs(args) {
			return nil, errors.New("need some arguments")
		}

		sql := "update adaptiveip set"
		var updateSet = ""
		if networkAddressOk {
			updateSet += " network_ip = '" + adaptiveip.NetworkAddress + "', "
		}
		if netmaskOk {
			updateSet += " netmask = '" + adaptiveip.Netmask + "', "
		}
		if gatewayOk {
			updateSet += " gateway = '" + adaptiveip.Gateway + "', "
		}
		if startIPaddressOk {
			updateSet += " start_ip_address = '" + adaptiveip.StartIPAddress + "', "
		}
		if endIPaddressOk {
			updateSet += " end_ip_address = '" + adaptiveip.EndIPAddress + "', "
		}
		sql += updateSet[0:len(updateSet)-2] + " where uuid = ?"

		logger.Logger.Println("update_adaptiveip sql : ", sql)

		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()

		result, err2 := stmt.Exec(adaptiveip.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err
		}
		logger.Logger.Println(result.LastInsertId())
		return adaptiveip, nil
	}

	return nil, errors.New("need uuid argument")
}

// DeleteAdaptiveIP - ish
func DeleteAdaptiveIP(args map[string]interface{}) (interface{}, error) {
	var err error

	requestedUUID, ok := args["uuid"].(string)
	if ok {
		sql := "delete from adaptiveip where uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()
		result, err2 := stmt.Exec(requestedUUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err
		}
		logger.Logger.Println(result.RowsAffected())

		return requestedUUID, nil
	}

	return requestedUUID, err
}
