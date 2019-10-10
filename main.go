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
	"strconv"
)

func main() {

	fmt.Println(simpleGitHub.AddTwo(5, 6))

	if !checkroot.CheckRoot() {
		return
	}

	if !logger.Prepare() {
		return
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	err := mysql.Prepare()
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
