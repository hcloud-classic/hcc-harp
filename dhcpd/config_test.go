package dhcpd

import (
	"hcc/harp/checkroot"
	"hcc/harp/config"
	"hcc/harp/logger"
	"testing"
)

func Test_CreateConfig(t *testing.T) {
	if !checkroot.CheckRoot() {
		t.Fatal("Failed to get root permission!")
	}

	if !logger.Prepare() {
		t.Fatal("Failed to prepare logger!")
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	var nodeUUIDs = []string{
		"48d08a00-b652-11e8-906e-000ffee02d5c",
		"d4f3a900-b674-11e8-906e-000ffee02d5c",
		"b9e43600-b4c8-11e8-906e-000ffee02d5c",
		"18aada80-b696-11e8-906e-000ffee02d5c"}

	err := CreateConfig("172.18.0.160", "255.255.255.240", "172.18.0.161",
		"172.18.0.10", "8.8.8.8", "google.com",
		6, nodeUUIDs, "48d08a00-b652-11e8-906e-000ffee02d5c", "CentOS 6", "test1")
	if err != nil {
		logger.Logger.Panic(err)
	}

	err = CreateConfig("192.168.110.0", "255.255.255.0", "192.168.110.254",
		"192.168.110.240", "8.8.8.8", "google.com",
		10, nodeUUIDs, "48d08a00-b652-11e8-906e-000ffee02d5c", "CentOS 6", "test2")
	if err != nil {
		logger.Logger.Panic(err)
	}
}
