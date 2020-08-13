package grpccli

import (
	"context"
	"google.golang.org/grpc"
	"hcc/harp/action/grpc/rpcviolin"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"strconv"
	"time"
)

var violinConn *grpc.ClientConn

func initViolin() error {
	var err error

	addr := config.Violin.ServerAddress + ":" + strconv.FormatInt(config.Violin.ServerPort, 10)
	logger.Logger.Println("Trying to connect to violin module (" + addr + ")")
	violinConn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Logger.Fatalf("Failed to connect violin module ("+addr+"): %v", err)
		return err
	}

	RC.violin = rpcviolin.NewViolinClient(violinConn)
	logger.Logger.Println("gRPC client connected to violin module")

	return nil
}

func cleanViolin() {
	_ = violinConn.Close()
}

// AllServerUUID : Get all of server UUIDs
func (rc *RPCClient) AllServerUUID() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resServerList, err := rc.violin.GetServerList(ctx, &rpcviolin.ReqGetServerList{})
	if err != nil {
		return nil, err
	}

	var serverUUIDs []string
	pserverList := resServerList.Server

	for i := range pserverList {
		serverUUIDs = append(serverUUIDs, pserverList[i].UUID)
	}

	return serverUUIDs, nil
}
