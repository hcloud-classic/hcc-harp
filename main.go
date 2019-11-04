package main

import (
	"hcc/harp/action/graphql"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/syscheck"
	"net/http"
	"strconv"
)

func main() {
	if !syscheck.CheckRoot() {
		return
	}

	if !syscheck.CheckArpingCommand() {
		return
	}

	if !logger.Prepare() {
		return
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()
	err := dhcpd.CheckLocalDHCPDConfig()
	if err != nil {
		logger.Logger.Panicln(err)
	}

	if !syscheck.CheckIfaceExist(config.ARPING.IfaceName) {
		return
	}

	err = mysql.Prepare()
	if err != nil {
		return
	}
	defer func() {
		_ = mysql.Db.Close()
	}()

	http.Handle("/graphql", graphql.GraphqlHandler)

	logger.Logger.Println("Server is running on port " + strconv.Itoa(int(config.HTTP.Port)))
	err = http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println("Failed to prepare http server!")
	}
}
