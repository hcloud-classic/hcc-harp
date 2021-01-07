package client

import (
	"context"
	"github.com/hcloud-classic/hcc_errors"
	"github.com/hcloud-classic/pb"
	"google.golang.org/grpc"
	"hcc/harp/action/grpc/errconv"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"strconv"
	"time"
)

var fluteConn *grpc.ClientConn

func initFlute() error {
	var err error

	addr := config.Flute.ServerAddress + ":" + strconv.FormatInt(config.Flute.ServerPort, 10)
	logger.Logger.Println("Trying to connect to flute module (" + addr + ")")

	for i := 0; i < int(config.Flute.ConnectionRetryCount); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Flute.ConnectionTimeOutMs)*time.Millisecond)
		fluteConn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logger.Logger.Println("Failed to connect flute module (" + addr + "): " + err.Error())
			logger.Logger.Println("Re-trying to connect to flute module (" +
				strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Flute.ConnectionRetryCount)) + ")")

			cancel()
			continue
		}

		RC.flute = pb.NewFluteClient(fluteConn)
		logger.Logger.Println("gRPC client connected to flute module")

		cancel()
		return nil
	}

	hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "initFlute(): retry count exceeded to connect flute module")).Stack()
	return (*hccErrStack)[0].ToError()
}

func closeFlute() {
	_ = fluteConn.Close()
}

// GetNode : Get infos of the node
func (rc *RPCClient) GetNode(uuid string) (*pb.Node, *hcc_errors.HccErrorStack) {
	var errStack *hcc_errors.HccErrorStack = nil

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	node, err := rc.flute.GetNode(ctx, &pb.ReqGetNode{UUID: uuid})
	if err != nil {
		hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpGrpcRequestError, "GetNode(): "+err.Error()))
		return nil, hccErrStack
	}
	if es := node.GetHccErrorStack(); es != nil {
		errStack = errconv.GrpcStackToHcc(&es)
	}

	return node.Node, errStack
}

// GetNodeList : Get the list of nodes by server UUID.
func (rc *RPCClient) GetNodeList(serverUUID string) ([]pb.Node, *hcc_errors.HccErrorStack) {
	var nodeList []pb.Node
	var errStack *hcc_errors.HccErrorStack = nil

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	pnodeList, err := rc.flute.GetNodeList(ctx, &pb.ReqGetNodeList{Node: &pb.Node{ServerUUID: serverUUID}})
	if err != nil {
		hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpGrpcRequestError, "GetNodeList(): "+err.Error()))
		return nil, hccErrStack
	}

	for _, pnode := range pnodeList.Node {
		nodeList = append(nodeList, pb.Node{
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
