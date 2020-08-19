package grpcsrv

import (
	"context"
	pb "hcc/harp/action/grpc/rpcharp"
	"hcc/harp/dao"
	"hcc/harp/lib/adaptiveip"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/pf"
)

type harpServer struct {
	pb.UnimplementedHarpServer
}

func returnSubnet(subnet *pb.Subnet) *pb.Subnet {
	return &pb.Subnet{
		UUID:           subnet.UUID,
		NetworkIP:      subnet.NetworkIP,
		Netmask:        subnet.Netmask,
		Gateway:        subnet.Gateway,
		NextServer:     subnet.NextServer,
		NameServer:     subnet.NameServer,
		DomainName:     subnet.DomainName,
		ServerUUID:     subnet.ServerUUID,
		LeaderNodeUUID: subnet.LeaderNodeUUID,
		OS:             subnet.OS,
		SubnetName:     subnet.SubnetName,
		CreatedAt:      subnet.CreatedAt,
	}
}

func (s *harpServer) CreateSubnet(_ context.Context, in *pb.ReqCreateSubnet) (*pb.ResCreateSubnet, error) {
	logger.Logger.Println("Request received: CreateSubnet()")

	subnet, err := dao.CreateSubnet(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResCreateSubnet{Subnet: returnSubnet(subnet)}, nil
}

func (s *harpServer) GetSubnet(_ context.Context, in *pb.ReqGetSubnet) (*pb.ResGetSubnet, error) {
	logger.Logger.Println("Request received: GetSubnet()")

	subnet, err := dao.ReadSubnet(in.GetUUID())
	if err != nil {
		return nil, err
	}

	return &pb.ResGetSubnet{Subnet: returnSubnet(subnet)}, nil
}

func (s *harpServer) GetSubnetByServer(_ context.Context, in *pb.ReqGetSubnetByServer) (*pb.ResGetSubnetByServer, error) {
	logger.Logger.Println("Request received: GetSubnetByServer()")

	subnet, err := dao.ReadSubnetByServer(in.GetServerUUID())
	if err != nil {
		return nil, err
	}

	return &pb.ResGetSubnetByServer{Subnet: returnSubnet(subnet)}, nil
}

func (s *harpServer) GetSubnetList(_ context.Context, in *pb.ReqGetSubnetList) (*pb.ResGetSubnetList, error) {
	logger.Logger.Println("Request received: GetSubnetList()")

	subnetList, err := dao.ReadSubnetList(in)
	if err != nil {
		return nil, err
	}

	return subnetList, nil
}

func (s *harpServer) GetSubnetNum(_ context.Context, _ *pb.Empty) (*pb.ResGetSubnetNum, error) {
	logger.Logger.Println("Request received: GetSubnetNum()")

	subnetNum, err := dao.ReadSubnetNum()
	if err != nil {
		return nil, err
	}

	return subnetNum, nil
}

func (s *harpServer) UpdateSubnet(_ context.Context, in *pb.ReqUpdateSubnet) (*pb.ResUpdateSubnet, error) {
	logger.Logger.Println("Request received: UpdateSubnet()")

	updateSubnet, err := dao.UpdateSubnet(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResUpdateSubnet{Subnet: updateSubnet}, nil
}

func (s *harpServer) DeleteSubnet(_ context.Context, in *pb.ReqDeleteSubnet) (*pb.ResDeleteSubnet, error) {
	logger.Logger.Println("Request received: DeleteSubnet()")

	uuid, err := dao.DeleteSubnet(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResDeleteSubnet{UUID: uuid}, nil
}

func (s *harpServer) CreateAdaptiveIPSetting(_ context.Context, in *pb.ReqCreateAdaptiveIPSetting) (*pb.ResCreateAdaptiveIPSetting, error) {
	logger.Logger.Println("Request received: CreateAdaptiveIPSetting()")

	adaptiveIPSetting, err := adaptiveip.WriteNetworkConfigAndReloadHarpNetwork(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResCreateAdaptiveIPSetting{AdaptiveipSetting: adaptiveIPSetting}, nil
}

func (s *harpServer) GetAdaptiveIPSetting(_ context.Context, _ *pb.Empty) (*pb.ResGetAdaptiveIPSetting, error) {
	logger.Logger.Println("Request received: GetAdaptiveIPSetting()")

	adaptiveIPNetwork := configext.GetAdaptiveIPNetwork()

	return &pb.ResGetAdaptiveIPSetting{AdaptiveipSetting: adaptiveIPNetwork}, nil
}

func (s *harpServer) GetAdaptiveIPAvailableIPList(_ context.Context, _ *pb.Empty) (*pb.ResGetAdaptiveIPAvailableIPList, error) {
	logger.Logger.Println("Request received: GetAdaptiveIPAvailableIPList()")

	adaptiveIPAvailableIPList := pf.GetAvailableIPList()

	return &pb.ResGetAdaptiveIPAvailableIPList{AdaptiveipAvailableipList: adaptiveIPAvailableIPList}, nil
}

func (s *harpServer) CreateAdaptiveIPServer(_ context.Context, in *pb.ReqCreateAdaptiveIPServer) (*pb.ResCreateAdaptiveIPServer, error) {
	logger.Logger.Println("Request received: CreateAdaptiveIPServer()")

	adaptiveIPServer, err := dao.CreateAdaptiveIPServer(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResCreateAdaptiveIPServer{AdaptiveipServer: adaptiveIPServer}, nil
}

func (s *harpServer) GetAdaptiveIPServer(_ context.Context, in *pb.ReqGetAdaptiveIPServer) (*pb.ResGetAdaptiveIPServer, error) {
	logger.Logger.Println("Request received: GetAdaptiveIPServer()")

	adaptiveIPServer, err := dao.ReadAdaptiveIPServer(in.GetServerUUID())
	if err != nil {
		return nil, err
	}

	return &pb.ResGetAdaptiveIPServer{AdaptiveipServer: adaptiveIPServer}, nil
}

func (s *harpServer) GetAdaptiveIPServerList(_ context.Context, in *pb.ReqGetAdaptiveIPServerList) (*pb.ResGetAdaptiveIPServerList, error) {
	logger.Logger.Println("Request received: GetAdaptiveIPServerList()")

	adaptiveIPServerList, err := dao.ReadAdaptiveIPServerList(in)
	if err != nil {
		return nil, err
	}

	return adaptiveIPServerList, nil
}

func (s *harpServer) GetAdaptiveIPServerNum(_ context.Context, _ *pb.Empty) (*pb.ResGetAdaptiveIPServerNum, error) {
	logger.Logger.Println("Request received: GetAdaptiveIPServerNum()")

	adaptiveIPServerNum, err := dao.ReadAdaptiveIPServerNum()
	if err != nil {
		return nil, err
	}

	return adaptiveIPServerNum, nil
}

func (s *harpServer) DeleteAdaptiveIPServer(_ context.Context, in *pb.ReqDeleteAdaptiveIPServer) (*pb.ResDeleteAdaptiveIPServer, error) {
	logger.Logger.Println("Request received: DeleteAdaptiveIPServer()")

	serverUUID, err := dao.DeleteAdaptiveIPServer(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResDeleteAdaptiveIPServer{ServerUUID: serverUUID}, nil
}

func (s *harpServer) CreateDHPCDConf(_ context.Context, in *pb.ReqCreateDHPCDConf) (*pb.ResCreateDHPCDConf, error) {
	logger.Logger.Println("Request received: CreateDHPCDConf()")

	result, err := dhcpd.CreateDHCPDConfig(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResCreateDHPCDConf{Result: result}, nil
}

func (s *harpServer) DeleteDHPCDConf(_ context.Context, in *pb.ReqDeleteDHPCDConf) (*pb.ResDeleteDHPCDConf, error) {
	logger.Logger.Println("Request received: DeleteDHPCDConf()")

	result, err := dhcpd.DeleteDHCPDConfigFile(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResDeleteDHPCDConf{Result: result}, nil
}
