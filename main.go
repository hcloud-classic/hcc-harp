package main

import (
	"hcc/harp/action/graphql"
	harpEnd "hcc/harp/end"
	harpInit "hcc/harp/init"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"net/http"
	"strconv"
)

func init() {
	err := harpInit.MainInit()
	if err != nil {
		panic(err)
	}
}

func main() {
	defer func() {
		harpEnd.MainEnd()
	}()

	http.Handle("/graphql", graphql.GraphqlHandler)
	logger.Logger.Println("Opening server on port " + strconv.Itoa(int(config.HTTP.Port)) + "...")
	err := http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println(err)
		logger.Logger.Println("Failed to prepare http server!")
	}
}
