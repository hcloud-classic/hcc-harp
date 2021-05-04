package server

import (
	"context"
	"hcc/harp/action/grpc/errconv"
	"hcc/harp/dao"
	"hcc/harp/lib/adaptiveip"
	"hcc/harp/lib/configext"
	"hcc/harp/lib/dhcpd"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
)

type harpServer struct {
	pb.UnimplementedHarpServer
}

func returnSubnet(subnet *pb.Subnet) *pb.Subnet {
	return &pb.Subnet{
		UUID:           subnet.UUID,
		GroupID:        subnet.GroupID,
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
	subnet, errCode, errStr := dao.CreateSubnet(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResCreateSubnet{Subnet: &pb.Subnet{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResCreateSubnet{Subnet: returnSubnet(subnet)}, nil
}

func (s *harpServer) GetSubnet(_ context.Context, in *pb.ReqGetSubnet) (*pb.ResGetSubnet, error) {
	subnet, errCode, errStr := dao.ReadSubnet(in.GetUUID())
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetSubnet{Subnet: &pb.Subnet{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetSubnet{Subnet: returnSubnet(subnet)}, nil
}

func (s *harpServer) GetSubnetByServer(_ context.Context, in *pb.ReqGetSubnetByServer) (*pb.ResGetSubnetByServer, error) {
	subnet, errCode, errStr := dao.ReadSubnetByServer(in.GetServerUUID())
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetSubnetByServer{Subnet: &pb.Subnet{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetSubnetByServer{Subnet: returnSubnet(subnet)}, nil
}

func (s *harpServer) GetSubnetList(_ context.Context, in *pb.ReqGetSubnetList) (*pb.ResGetSubnetList, error) {
	resGetSubnetList, errCode, errStr := dao.ReadSubnetList(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetSubnetList{Subnet: []*pb.Subnet{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetSubnetList{Subnet: resGetSubnetList.Subnet}, nil
}

func (s *harpServer) GetAvailableSubnetList(_ context.Context, in *pb.ReqGetAvailableSubnetList) (*pb.ResGetAvailableSubnetList, error) {
	resGetSubnetList, errCode, errStr := dao.ReadAvailableSubnetList(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetAvailableSubnetList{Subnet: []*pb.Subnet{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetAvailableSubnetList{Subnet: resGetSubnetList.Subnet}, nil
}

func (s *harpServer) GetSubnetNum(_ context.Context, in *pb.ReqGetSubnetNum) (*pb.ResGetSubnetNum, error) {
	resGetSubnetNum, errCode, errStr := dao.ReadSubnetNum(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetSubnetNum{Num: 0, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetSubnetNum{Num: resGetSubnetNum.Num}, nil
}

func (s *harpServer) UpdateSubnet(_ context.Context, in *pb.ReqUpdateSubnet) (*pb.ResUpdateSubnet, error) {
	updateSubnet, errCode, errStr := dao.UpdateSubnet(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResUpdateSubnet{Subnet: &pb.Subnet{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResUpdateSubnet{Subnet: updateSubnet}, nil
}

func (s *harpServer) DeleteSubnet(_ context.Context, in *pb.ReqDeleteSubnet) (*pb.ResDeleteSubnet, error) {
	deleteSubnet, errCode, errStr := dao.DeleteSubnet(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResDeleteSubnet{Subnet: &pb.Subnet{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResDeleteSubnet{Subnet: deleteSubnet}, nil
}

func (s *harpServer) CreateAdaptiveIPSetting(_ context.Context, in *pb.ReqCreateAdaptiveIPSetting) (*pb.ResCreateAdaptiveIPSetting, error) {
	adaptiveIPSetting, err := adaptiveip.WriteNetworkConfigAndReloadHarpNetwork(in)
	if err != nil {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpInternalOperationFail, "WriteNetworkConfigAndReloadHarpNetwork(): "+err.Error()))
		return &pb.ResCreateAdaptiveIPSetting{AdaptiveipSetting: &pb.AdaptiveIPSetting{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResCreateAdaptiveIPSetting{AdaptiveipSetting: adaptiveIPSetting}, nil
}

func (s *harpServer) GetAdaptiveIPSetting(_ context.Context, _ *pb.Empty) (*pb.ResGetAdaptiveIPSetting, error) {
	adaptiveIPNetwork := configext.GetAdaptiveIPNetwork()

	return &pb.ResGetAdaptiveIPSetting{AdaptiveipSetting: adaptiveIPNetwork}, nil
}

func (s *harpServer) GetAdaptiveIPAvailableIPList(_ context.Context, _ *pb.Empty) (*pb.ResGetAdaptiveIPAvailableIPList, error) {
	adaptiveIPAvailableIPList, err := adaptiveip.GetAvailableIPList()
	if err != nil {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpInternalOperationFail, "GetAdaptiveIPAvailableIPList(): "+err.Error()))
		return &pb.ResGetAdaptiveIPAvailableIPList{AdaptiveipAvailableipList: &pb.AdaptiveIPAvailableIPList{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetAdaptiveIPAvailableIPList{AdaptiveipAvailableipList: adaptiveIPAvailableIPList}, nil
}

func (s *harpServer) CreateAdaptiveIPServer(_ context.Context, in *pb.ReqCreateAdaptiveIPServer) (*pb.ResCreateAdaptiveIPServer, error) {
	adaptiveIPServer, errCode, errStr := dao.CreateAdaptiveIPServer(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResCreateAdaptiveIPServer{AdaptiveipServer: &pb.AdaptiveIPServer{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResCreateAdaptiveIPServer{AdaptiveipServer: adaptiveIPServer}, nil
}

func (s *harpServer) GetAdaptiveIPServer(_ context.Context, in *pb.ReqGetAdaptiveIPServer) (*pb.ResGetAdaptiveIPServer, error) {
	adaptiveIPServer, errCode, errStr := dao.ReadAdaptiveIPServer(in.GetServerUUID())
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetAdaptiveIPServer{AdaptiveipServer: &pb.AdaptiveIPServer{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetAdaptiveIPServer{AdaptiveipServer: adaptiveIPServer}, nil
}

func (s *harpServer) GetAdaptiveIPServerList(_ context.Context, in *pb.ReqGetAdaptiveIPServerList) (*pb.ResGetAdaptiveIPServerList, error) {
	adaptiveIPServerList, errCode, errStr := dao.ReadAdaptiveIPServerList(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetAdaptiveIPServerList{AdaptiveipServer: []*pb.AdaptiveIPServer{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return adaptiveIPServerList, nil
}

func (s *harpServer) GetAdaptiveIPServerNum(_ context.Context, in *pb.ReqGetAdaptiveIPServerNum) (*pb.ResGetAdaptiveIPServerNum, error) {
	adaptiveIPServerNum, errCode, errStr := dao.ReadAdaptiveIPServerNum(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetAdaptiveIPServerNum{Num: 0, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return adaptiveIPServerNum, nil
}

func (s *harpServer) DeleteAdaptiveIPServer(_ context.Context, in *pb.ReqDeleteAdaptiveIPServer) (*pb.ResDeleteAdaptiveIPServer, error) {
	serverUUID, errCode, errStr := dao.DeleteAdaptiveIPServer(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResDeleteAdaptiveIPServer{ServerUUID: "", HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResDeleteAdaptiveIPServer{ServerUUID: serverUUID}, nil
}

func (s *harpServer) CreateDHCPDConf(_ context.Context, in *pb.ReqCreateDHCPDConf) (*pb.ResCreateDHCPDConf, error) {
	result, err := dhcpd.CreateDHCPDConfig(in)
	if err != nil {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpInternalDHCPDError, "CreateDHCPDConfig(): "+err.Error()))
		return &pb.ResCreateDHCPDConf{Result: "", HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResCreateDHCPDConf{Result: result}, nil
}

func (s *harpServer) DeleteDHCPDConf(_ context.Context, in *pb.ReqDeleteDHCPDConf) (*pb.ResDeleteDHCPDConf, error) {
	result, err := dhcpd.DeleteDHCPDConfig(in)
	if err != nil {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpInternalDHCPDError, "DeleteDHCPDConfig(): "+err.Error()))
		return &pb.ResDeleteDHCPDConf{Result: "", HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResDeleteDHCPDConf{Result: result}, nil
}

func (s *harpServer) GetPortForwardingList(_ context.Context, in *pb.ReqGetPortForwardingList) (*pb.ResGetPortForwardingList, error) {
	portForwardingList, errCode, errStr := dao.ReadPortForwardingList(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetPortForwardingList{PortForwarding: []*pb.PortForwarding{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return portForwardingList, nil
}

func (s *harpServer) CreatePortForwarding(_ context.Context, in *pb.ReqCreatePortForwarding) (*pb.ResCreatePortForwarding, error) {
	portForwarding, errCode, errStr := dao.CreatePortForwarding(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResCreatePortForwarding{PortForwarding: &pb.PortForwarding{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResCreatePortForwarding{PortForwarding: portForwarding}, nil
}

func (s *harpServer) DeletePortForwarding(_ context.Context, in *pb.ReqDeletePortForwarding) (*pb.ResDeletePortForwarding, error) {
	serverUUID, errCode, errStr := dao.DeletePortForwarding(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResDeletePortForwarding{ServerUUID: "", HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResDeletePortForwarding{ServerUUID: serverUUID}, nil
}
