package main

import (
	"GraphQL_harp/harpcheckroot"
	"GraphQL_harp/harpconfig"
	"GraphQL_harp/harpgraphql"
	"GraphQL_harp/harplogger"
	"GraphQL_harp/harpmysql"
	"net/http"
)

func main() {
	if !harpcheckroot.CheckRoot() {
		return
	}

	if !harplogger.Prepare() {
		return
	}
	defer harplogger.FpLog.Close()

	err := harpmysql.Prepare()
	if err != nil {
		return
	}
	defer harpmysql.Db.Close()

	http.Handle("/graphql", harpgraphql.GraphqlHandler)

	harplogger.Logger.Println("Server is running on port " + harpconfig.HTTPPort)
	err = http.ListenAndServe(":"+harpconfig.HTTPPort, nil)
	if err != nil {
		harplogger.Logger.Println("Failed to prepare http server!")
	}
}
