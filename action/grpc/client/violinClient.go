package client

import (
	"context"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"google.golang.org/grpc"
	"hcc/harp/action/grpc/errconv"
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

	for i := 0; i < int(config.Violin.ConnectionRetryCount); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Violin.ConnectionTimeOutMs)*time.Millisecond)
		violinConn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logger.Logger.Println("Failed to connect violin module (" + addr + "): " + err.Error())
			logger.Logger.Println("Re-trying to connect to violin module (" +
				strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Violin.ConnectionRetryCount)) + ")")

			cancel()
			continue
		}

		RC.violin = pb.NewViolinClient(violinConn)
		logger.Logger.Println("gRPC client connected to violin module")

		cancel()
		return nil
	}

	hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "initViolin(): retry count exceeded to connect violin module")).Stack()
	return (*hccErrStack)[0].ToError()
}

func closeViolin() {
	_ = violinConn.Close()
}

// AllServerUUID : Get all of server UUIDs
func (rc *RPCClient) AllServerUUID() ([]string, *hcc_errors.HccErrorStack) {
	var serverUUIDs []string
	var errStack *hcc_errors.HccErrorStack = nil

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resServerList, err := rc.violin.GetServerList(ctx, &pb.ReqGetServerList{})
	if err != nil {
		hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpGrpcRequestError, "AllServerUUID(): "+err.Error()))
		return nil, hccErrStack
	}

	if pserverList := resServerList.GetServer(); pserverList != nil {
		for i := range pserverList {
			serverUUIDs = append(serverUUIDs, pserverList[i].UUID)
		}
	}

	if es := resServerList.GetHccErrorStack(); es != nil {
		errStack = errconv.GrpcStackToHcc(es)
	}

	return serverUUIDs, errStack
}
