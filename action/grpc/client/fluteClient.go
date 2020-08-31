package client

import (
	"context"
	"google.golang.org/grpc"
	"hcc/harp/action/grpc/errconv"
	"hcc/harp/action/grpc/pb/rpcflute"
	"hcc/harp/lib/config"
	"hcc/harp/lib/errors"
	"hcc/harp/lib/logger"
	"strconv"
	"time"
)

var fluteConn *grpc.ClientConn

func initFlute() *errors.HccError {
	var err error

	addr := config.Flute.ServerAddress + ":" + strconv.FormatInt(config.Flute.ServerPort, 10)
	logger.Logger.Println("Trying to connect to flute module (" + addr + ")")

	for i := 0; i < int(config.Flute.ConnectionRetryCount); i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Flute.ConnectionTimeOutMs)*time.Millisecond)
		fluteConn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logger.Logger.Println("Failed to connect flute module (" + addr + "): " + err.Error())
			logger.Logger.Println("Re-trying to connect to flute module (" +
				strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Flute.ConnectionRetryCount)) + ")")
			continue
		}

		RC.flute = rpcflute.NewFluteClient(fluteConn)
		logger.Logger.Println("gRPC client connected to flute module")

		return nil
	}

	return errors.NewHccError(errors.HarpInternalInitFail, "retry count exceeded to connect flute module")
}

func closeFlute() {
	_ = fluteConn.Close()
}

// GetNode : Get infos of the node
func (rc *RPCClient) GetNode(uuid string) (*rpcflute.Node, *errors.HccErrorStack) {
	var errStack *errors.HccErrorStack = nil

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	node, err := rc.flute.GetNode(ctx, &rpcflute.ReqGetNode{UUID: uuid})
	if err != nil {
		return nil, errors.NewHccErrorStack(errors.NewHccError(errors.HarpGrpcRequestError, "GetNode "+err.Error()))
	}
	if es := node.GetHccErrorStack(); es != nil {
		errStack = errconv.GrpcStackToHcc(&es)
	}

	return node.Node, errStack
}

// GetNodeList : Get the list of nodes by server UUID.
func (rc *RPCClient) GetNodeList(serverUUID string) ([]rpcflute.Node, *errors.HccErrorStack) {
	var nodeList []rpcflute.Node
	var errStack *errors.HccErrorStack = nil

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	pnodeList, err := rc.flute.GetNodeList(ctx, &rpcflute.ReqGetNodeList{Node: &rpcflute.Node{ServerUUID: serverUUID}})
	if err != nil {
		return nil, errors.NewHccErrorStack(errors.NewHccError(errors.HarpGrpcRequestError, "GetNodeList "+err.Error()))
	}

	for _, pnode := range pnodeList.Node {
		nodeList = append(nodeList, rpcflute.Node{
			UUID:        pnode.UUID,
			ServerUUID:  pnode.ServerUUID,
			BmcMacAddr:  pnode.BmcMacAddr,
			BmcIP:       pnode.BmcIP,
			PXEMacAddr:  pnode.PXEMacAddr,
			Status:      pnode.Status,
			CPUCores:    pnode.CPUCores,
			Memory:      pnode.Memory,
			Description: pnode.Description,
			Active:      pnode.Active,
			CreatedAt:   pnode.CreatedAt,
		})
	}
	if es := pnodeList.GetHccErrorStack(); es != nil {
		errStack = errconv.GrpcStackToHcc(&es)
	}

	return nodeList, errStack
}
