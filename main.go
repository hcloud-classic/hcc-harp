package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	harpGrpc "hcc/harp/action/grpc"
	harpEnd "hcc/harp/end"
	harpInit "hcc/harp/init"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"hcc/harp/pb"
	"net"
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

	//http.Handle("/graphql", graphql.GraphqlHandler)
	//logger.Logger.Println("Opening server on port " + strconv.Itoa(int(config.HTTP.Port)) + "...")
	//err := http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	//if err != nil {
	//	logger.Logger.Println(err)
	//	logger.Logger.Println("Failed to prepare http server!")
	//}

	lis, err := net.Listen("tcp", ":" + strconv.Itoa(int(config.HTTP.Port)))
	if err != nil {
		logger.Logger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHarpServer(s, &harpGrpc.Server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		logger.Logger.Fatalf("failed to serve: %v", err)
	}
}
