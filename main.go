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

	dhcpd.ConfParser("192.168.110.0", "255.255.255.252", nil, "", "", "")

	http.Handle("/graphql", graphql.GraphqlHandler)

	logger.Logger.Println("Server is running on port " + strconv.Itoa(int(config.HTTP.Port)))
	err = http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println("Failed to prepare http server!")
	}
}
