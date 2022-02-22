package mysql

import (
	"hcc/harp/action/grpc/client"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/syscheck"
	"innogrid.com/hcloud-classic/hcc_errors"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	err := syscheck.CheckOS()
	if err != nil {
		t.Fatal("Please run harp module on Linux machine.")
	}

	err = logger.Init()
	if err != nil {
		hcc_errors.SetErrLogger(logger.Logger)
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	hcc_errors.SetErrLogger(logger.Logger)
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Init()

	err = client.Init()
	if err != nil {
		t.Fatal(err)
	}

	err = Init()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = Db.Close()
	}()
}
