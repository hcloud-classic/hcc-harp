package grpc

import (
	"context"
	pb "hcc/harp/action/grpc/rpcharp"
	"hcc/harp/dao"
)

// Server : Server type of Harp's Grpc
type Server struct {
	pb.UnimplementedHarpServer
}

func returnSubnet(subnet *pb.Subnet) *pb.Subnet {
	return &pb.Subnet{
		Uuid: subnet.Uuid,
		NetworkIp: subnet.NetworkIp,
		Netmask: subnet.Netmask,
		Gateway: subnet.Gateway,
		NextServer: subnet.NextServer,
		NameServer: subnet.NameServer,
		DomainName: subnet.DomainName,
		ServerUuid: subnet.ServerUuid,
		LeaderNodeUuid: subnet.LeaderNodeUuid,
		Os: subnet.Os,
		SubnetName: subnet.SubnetName,
		CreatedAt: subnet.CreatedAt,
	}
}

func (s *Server) CreateSubnet(_ context.Context, in *pb.ReqCreateSubnet) (*pb.ResCreateSubnet, error) {
	subnet, err := dao.CreateSubnet(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResCreateSubnet{Subnet: returnSubnet(subnet)}, nil
}

func (s *Server) GetSubnet(_ context.Context, in *pb.ReqGetSubnet) (*pb.ResGetSubnet, error) {
	subnet, err := dao.ReadSubnet(in.GetUuid())
	if err != nil {
		return nil, err
	}

	return &pb.ResGetSubnet{Subnet: returnSubnet(subnet)}, nil
}

func (s *Server) GetSubnetList(_ context.Context, in *pb.ReqGetSubnetList) (*pb.ResGetSubnetList, error) {
	subnetList, err := dao.ReadSubnetList(in)
	if err != nil {
		return nil, err
	}

	return subnetList, nil
}

//func (s *Server) GetSubnetNum(_ context.Context, _ *pb.Empty) (*pb.SubnetList, error) {}