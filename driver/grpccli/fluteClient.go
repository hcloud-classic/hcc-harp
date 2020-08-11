package grpccli

import (
	"context"
	"google.golang.org/grpc"
	"hcc/harp/action/grpc/rpcflute"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"strconv"
	"time"
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

// GetNode : Get infos of the node
func (rc *RpcClient) GetNode(uuid string) (*rpcflute.Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	node, err := rc.flute.GetNode(ctx, &rpcflute.ReqGetNode{UUID: uuid})
	if err != nil {
		return nil, err
	}

	return node.Node, nil
}

// GetNodeList : Get the list of nodes by server UUID.
func (rc *RpcClient) GetNodeList(serverUUID string) ([]rpcflute.Node, error) {
	var nodeList []rpcflute.Node

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	pnodeList, err := rc.flute.GetNodeList(ctx, &rpcflute.ReqGetNodeList{Node: &rpcflute.Node{ServerUUID: serverUUID}})
	if err != nil {
		return nil, err
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

	return nodeList, nil
}
