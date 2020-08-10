package grpccli

import (
	"google.golang.org/grpc"
	"hcc/harp/action/grpc/rpcflute"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"strconv"
)

var fluteConn *grpc.ClientConn
var err error

func initFlute() error {
	addr := config.Flute.ServerAddress + ":" + strconv.FormatInt(config.Flute.ServerPort, 10)
	logger.Logger.Println("Trying to connect to flute module (" + addr + ")")
	fluteConn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Logger.Fatalf("Failed to connect flute module ("+addr+"): %v", err)
		return err
	}

	RC.flute = rpcflute.NewFluteClient(fluteConn)
	logger.Logger.Println("gRPC client connected to flute module")

	return nil
}

func cleanFlute() {
	_ = fluteConn.Close()
}
