package client

import (
	"hcc/harp/action/grpc/pb/rpcflute"
	"hcc/harp/action/grpc/pb/rpcviolin"
)

// RPCClient : Struct type of gRPC clients
type RPCClient struct {
	flute  rpcflute.FluteClient
	violin rpcviolin.ViolinClient
}

// RC : Exported variable pointed to RPCClient
var RC = &RPCClient{}

// Init : Initialize clients of gRPC
func Init() error {
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

// End : Close connections of gRPC clients
func End() {
	closeFlute()
	closeViolin()
}
