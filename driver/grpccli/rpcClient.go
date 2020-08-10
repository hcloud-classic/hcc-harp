package grpccli

import (
	"hcc/harp/action/grpc/rpcflute"
)

type RpcClient struct {
	flute rpcflute.FluteClient
}

var RC = &RpcClient{}

func InitGRPCClient() error {
	err := initFlute()
	if err != nil {
		return err
	}

	return nil
}

func CleanGRPCClient() {
	cleanFlute()
}
