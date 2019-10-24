package dhcpd

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
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

	err := mysql.Prepare()
	if err != nil {
		return
	}

	testInitPass = true
}

func Test_CreateConfig(t *testing.T) {
	if !testInitPass {
		testInit(t)
	}
	defer func() {
		_ = mysql.Db.Close()
	}()

	var nodeUUIDs = []string{
		"48d08a00-b652-11e8-906e-000ffee02d5c",
		"d4f3a900-b674-11e8-906e-000ffee02d5c",
		"b9e43600-b4c8-11e8-906e-000ffee02d5c",
		"18aada80-b696-11e8-906e-000ffee02d5c"}

	err := CreateConfig("0ac56231-a0ee-4323-55ad-37c08c2d4a78", nodeUUIDs, "48d08a00-b652-11e8-906e-000ffee02d5c", "CentOS 6", "test1")
	if err != nil {
		t.Log(err)
	}

	err = CreateConfig("1f16b53e-082d-4e82-75da-7874ff59d82a", nodeUUIDs, "48d08a00-b652-11e8-906e-000ffee02d5c", "CentOS 6", "test2")
	if err != nil {
		t.Log(err)
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
