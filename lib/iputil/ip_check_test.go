package iputil

import (
	"net"
	"testing"
)

func Test_CheckIP(t *testing.T) {
	netIP := CheckValidIP("Vaild IP")
	if netIP == nil {
		t.Fatal("wrong network IP")
	}

	mask, err := CheckNetmask("Vaild Netmask")
	if err != nil {
		t.Fatal(err)
	}

	ipNet := net.IPNet{
		IP:   netIP,
		Mask: mask,
	}

	err = CheckGateway(ipNet, "Valid Gateway")
	if err != nil {
		t.Fatal(err)
	}
}
