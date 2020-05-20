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

	var nodeUUIDs = []string{}

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
	}
}
