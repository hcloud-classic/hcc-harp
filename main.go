package main

import (
	"hcc/harp/checkroot"
	"hcc/harp/config"
	"hcc/harp/dhcpd"
	"hcc/harp/graphql"
	"hcc/harp/logger"
	"hcc/harp/mysql"
	"net/http"
	"strconv"
)

func main() {
	if !checkroot.CheckRoot() {
		return
	}

	if !logger.Prepare() {
		return
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	err := mysql.Prepare()
	if err != nil {
		return
	}
	defer func() {
		_ = mysql.Db.Close()
	}()

	err = dhcpd.CheckLocalDHCPDConfig()
	if err != nil {
		logger.Logger.Panic(err)
	}

	err = dhcpd.UpdateHarpDHCPDConfig()
	if err != nil {
		logger.Logger.Panic(err)
	}

	var nodeUUIDs = []string{
		"48d08a00-b652-11e8-906e-000ffee02d5c",
		"d4f3a900-b674-11e8-906e-000ffee02d5c",
		"b9e43600-b4c8-11e8-906e-000ffee02d5c",
		"18aada80-b696-11e8-906e-000ffee02d5c"}

	err = dhcpd.CreateConfig("172.18.0.160", "255.255.255.240", "172.18.0.161",
		"172.18.0.10", "8.8.8.8", "google.com",
		6, nodeUUIDs, "48d08a00-b652-11e8-906e-000ffee02d5c", "CentOS 6", "hcc")
	if err != nil {
		logger.Logger.Panic(err)
	}

	err = dhcpd.CreateConfig("192.168.110.0", "255.255.255.0", "192.168.110.254",
		"192.168.110.240", "8.8.8.8", "google.com",
		10, nodeUUIDs, "48d08a00-b652-11e8-906e-000ffee02d5c", "CentOS 6", "jolla")
	if err != nil {
		logger.Logger.Panic(err)
	}

	err = dhcpd.RestartDHCPDServer()
	if err != nil {
		logger.Logger.Panic(err)
	}

	http.Handle("/graphql", graphql.GraphqlHandler)

	logger.Logger.Println("Server is running on port " + strconv.Itoa(int(config.HTTP.Port)))
	err = http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println("Failed to prepare http server!")
	}
}
