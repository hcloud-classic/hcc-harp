package dhcpd

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/syscheck"
	"testing"
)

var testInitPass = false

func testInit(t *testing.T) {
	if !syscheck.CheckRoot() {
		t.Fatal("Failed to get root permission!")
	}

	if !logger.Prepare() {
		t.Fatal("Failed to prepare logger!")
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	testInitPass = true
}

func Test_CreateConfig(t *testing.T) {
	if !testInitPass {
		testInit(t)
	}

	var nodeUUIDs = []string{
		"48d08a00-b652-11e8-906e-000ffee02d5c",
		"d4f3a900-b674-11e8-906e-000ffee02d5c",
		"b9e43600-b4c8-11e8-906e-000ffee02d5c",
		"18aada80-b696-11e8-906e-000ffee02d5c"}

	err := CreateConfig("8d3f22a8-4010-49d4-8728-bb47889b13a6", "172.18.0.160", "255.255.255.240", "172.18.0.161",
		"172.18.0.10", "8.8.8.8", "google.com",
		6, nodeUUIDs, "48d08a00-b652-11e8-906e-000ffee02d5c", "CentOS 6", "test1")
	if err != nil {
		t.Fatal(err)
	}

	err = CreateConfig("8d3f22a8-4010-49d4-8728-bb47889b13a6", "192.168.110.0", "255.255.255.0", "192.168.110.254",
		"192.168.110.240", "8.8.8.8", "google.com",
		10, nodeUUIDs, "48d08a00-b652-11e8-906e-000ffee02d5c", "CentOS 6", "test2")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_CheckLocalDHCPDConfig(t *testing.T) {
	if !testInitPass {
		testInit(t)
	}

	err := CheckLocalDHCPDConfig()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_UpdateHarpDHCPDConfig(t *testing.T) {
	if !testInitPass {
		testInit(t)
	}

	err := UpdateHarpDHCPDConfig()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_RestartDHCPDServer(t *testing.T) {
	if !testInitPass {
		testInit(t)
	}

	err := RestartDHCPDServer()
	if err != nil {
		logger.Logger.Printf("Error occurred while restarting dhcpd service!\n"+
			"==> Error messages\n%s\n", err)
		// Ignore this error because we just try for testing.
	}
}
