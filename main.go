package main

import (
	"fmt"
	"github.com/mactsouk/go/simpleGitHub"
	"hcc/harp/checkroot"
	"hcc/harp/config"
	"hcc/harp/graphql"
	"hcc/harp/logger"
	"hcc/harp/mysql"
	"net/http"
)

func main() {

	fmt.Println(simpleGitHub.AddTwo(5, 6))

	if !checkroot.CheckRoot() {
		return
	}

	if !logger.Prepare() {
		return
	}
	defer logger.FpLog.Close()

	err := mysql.Prepare()
	if err != nil {
		return
	}
	defer mysql.Db.Close()

	http.Handle("/graphql", graphql.GraphqlHandler)

	logger.Logger.Println("Server is running on port " + config.HTTPPort)
	err = http.ListenAndServe(":"+config.HTTPPort, nil)
	if err != nil {
		logger.Logger.Println("Failed to prepare http server!")
	}
}
