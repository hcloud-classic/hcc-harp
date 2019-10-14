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



	err = dhcpd.RestartDHCPDServer()
	if err != nil {
		logger.Logger.Printf("Error occurred while restarting dhcpd service!\n" +
							"==> Error messages\n%s\n", err)
		logger.Logger.Panic(err)
	}

	http.Handle("/graphql", graphql.GraphqlHandler)

	logger.Logger.Println("Server is running on port " + strconv.Itoa(int(config.HTTP.Port)))
	err = http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println("Failed to prepare http server!")
	}
}
