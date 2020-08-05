package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	harpGrpc "hcc/harp/action/grpc"
	pb "hcc/harp/action/grpc/rpcharp"
	"hcc/harp/lib/adaptiveip"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/pf"
	"hcc/harp/lib/syscheck"
	"net"
	"strconv"
)

func init() {
	err := syscheck.CheckRoot()
	if err != nil {
		panic(err)
	}

	err = syscheck.CheckArpingCommand()
	if err != nil {
		panic(err)
	}

	err = logger.Init()
	if err != nil {
		panic(err)
	}

	config.Parser()

	err = mysql.Init()
	if err != nil {
		panic(err)
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.ExternalIfaceName)
	if err != nil {
		panic(err)
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.InternalIfaceName)
	if err != nil {
		panic(err)
	}

	err = dhcpd.CheckLocalDHCPDConfig()
	if err != nil {
		panic(err)
	}

	err = pf.PreparePFConfigFiles()
	if err != nil {
		panic(err)
	}

	err = adaptiveip.LoadHarpPFRules()
	if err != nil {
		panic(err)
	}
}

func main() {
	defer func() {
		mysql.End()
		logger.End()
	}()

	lis, err := net.Listen("tcp", ":" + strconv.Itoa(int(config.HTTP.Port)))
	if err != nil {
		logger.Logger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHarpServer(s, &harpGrpc.Server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		logger.Logger.Fatalf("failed to serve: %v", err)
	}
}
