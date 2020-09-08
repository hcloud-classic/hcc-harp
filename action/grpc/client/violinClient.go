package client

import (
	"context"
	"google.golang.org/grpc"
	"hcc/harp/action/grpc/errconv"
	"hcc/harp/action/grpc/pb/rpcviolin"
	"hcc/harp/lib/config"
	"hcc/harp/lib/errors"
	"hcc/harp/lib/logger"
	"strconv"
	"time"
)

var violinConn *grpc.ClientConn

func initViolin() *errors.HccError {
	var err error

	addr := config.Violin.ServerAddress + ":" + strconv.FormatInt(config.Violin.ServerPort, 10)
	logger.Logger.Println("Trying to connect to violin module (" + addr + ")")

	for i := 0; i < int(config.Violin.ConnectionRetryCount); i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Violin.ConnectionTimeOutMs)*time.Millisecond)
		violinConn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logger.Logger.Println("Failed to connect violin module (" + addr + "): " + err.Error())
			logger.Logger.Println("Re-trying to connect to violin module (" +
				strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Violin.ConnectionRetryCount)) + ")")
			continue
		}

		RC.violin = rpcviolin.NewViolinClient(violinConn)
		logger.Logger.Println("gRPC client connected to violin module")

		return nil
	}

	return errors.NewHccError(errors.HarpInternalInitFail, "retry count exceeded to connect violin module")
}

func closeViolin() {
	_ = violinConn.Close()
}

// AllServerUUID : Get all of server UUIDs
func (rc *RPCClient) AllServerUUID() ([]string, *errors.HccErrorStack) {
	var serverUUIDs []string
	var errStack *errors.HccErrorStack = nil

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resServerList, err := rc.violin.GetServerList(ctx, &rpcviolin.ReqGetServerList{})
	if err != nil {
		return nil, errors.NewHccErrorStack(errors.NewHccError(errors.HarpGrpcRequestError, "GetServerList "+err.Error()))
	}

	if pserverList := resServerList.GetServer(); pserverList != nil {
		for i := range pserverList {
			serverUUIDs = append(serverUUIDs, pserverList[i].UUID)
		}
	}

	if es := resServerList.GetHccErrorStack(); es != nil {
		errStack = errconv.GrpcStackToHcc(&es)
	}

	return serverUUIDs, errStack
}
