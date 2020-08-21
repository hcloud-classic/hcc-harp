package client

import (
	"hcc/harp/action/grpc/rpcflute"
	"hcc/harp/action/grpc/rpcviolin"
)

// RPCClient : Struct type of gRPC clients
type RPCClient struct {
	flute  rpcflute.FluteClient
	violin rpcviolin.ViolinClient
}

// RC : Exported variable pointed to RPCClient
var RC = &RPCClient{}

// InitGRPCClient : Initialize clients of gRPC
func InitGRPCClient() error {
	err := initFlute()
	if err != nil {
		return err
	}

	err = initViolin()
	if err != nil {
		return err
	}

	return nil
}

// CleanGRPCClient : Close connections of gRPC clients
func CleanGRPCClient() {
	cleanFlute()
	cleanViolin()
}
