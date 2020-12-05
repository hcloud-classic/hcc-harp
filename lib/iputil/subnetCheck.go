package iputil

import (
	"errors"
	pb "hcc/harp/action/grpc/pb/rpcharp"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"net"
)

func checkAClassPrivate(IP net.IP) bool {
	if IP[0] == 10 {
		return true
	}

	return false
}

func checkBClassPrivate(IP net.IP) bool {
	if IP[0] == 172 &&
		(IP[1] >= 16 && IP[1] <= 31) {
		return true
	}

	return false
}

func checkCClassPrivate(IP net.IP) bool {
	if IP[0] == 192 && IP[1] == 168 {
		return true
	}

	return false
}

// CheckPrivateSubnet : Check if given network address is private network address.
// Return error if given IP address is invalid or is not a network address.
// Return true if it is private address, return false otherwise.
func CheckPrivateSubnet(IP string, Netmask string) (bool, error) {
	netNetwork, err := CheckNetwork(IP, Netmask)
	if err != nil {
		return false, err
	}

	if netNetwork.IP.String() != IP {
		return false, errors.New("CheckPrivateSubnet(): invalid network address")
	}

	if checkAClassPrivate(netNetwork.IP) ||
		checkBClassPrivate(netNetwork.IP) ||
		checkCClassPrivate(netNetwork.IP) {
		return true, nil
	}

	return false, nil
}

func getSubnetList() ([]pb.Subnet, error) {
	var subnets []pb.Subnet
	var uuid string
	var networkIP string
	var netmask string

	sql := "select uuid, network_ip, netmask from subnet"
	stmt, err := mysql.Query(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &networkIP, &netmask)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}

		subnets = append(subnets, pb.Subnet{
			UUID:      uuid,
			NetworkIP: networkIP,
			Netmask:   netmask})
	}
	return subnets, nil
}

// CheckSubnetConflict : Check if given network address is conflict with one of subnet that stored in the database.
// Return true if conflicted, return false otherwise.
func CheckSubnetConflict(IP string, Netmask string, skipMine bool, oldSubnet *pb.Subnet) (bool, error) {
	netNetwork, err := CheckNetwork(IP, Netmask)
	if err != nil {
		return false, err
	}

	if netNetwork.IP.String() != IP {
		return false, errors.New("CheckPrivateSubnet(): invalid network address")
	}

	netmaskSize, _ := netNetwork.Mask.Size()

	subnetList, err := getSubnetList()
	if err != nil {
		return false, nil
	}

	for i := range subnetList {
		if skipMine && subnetList[i].UUID == oldSubnet.UUID {
			continue
		}

		var givenSubnetUpperNet *net.IPNet
		var subnetUpperNet *net.IPNet

		mask, _ := CheckNetmask(subnetList[i].Netmask)
		maskSize, _ := mask.Size()

		if netmaskSize >= maskSize {
			givenSubnetUpperNet, _ = CheckNetwork(IP, subnetList[i].Netmask)
			subnetUpperNet, _ = CheckNetwork(subnetList[i].NetworkIP, subnetList[i].Netmask)
		} else {
			givenSubnetUpperNet, _ = CheckNetwork(IP, Netmask)
			subnetUpperNet, _ = CheckNetwork(subnetList[i].NetworkIP, Netmask)
		}

		if subnetUpperNet.IP.Equal(givenSubnetUpperNet.IP) {
			return true, nil
		}
	}

	return false, nil
}
