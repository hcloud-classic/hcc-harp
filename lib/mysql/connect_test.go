package mysql

import (
	"fmt"
	"hcc/harp/lib/config"
	"hcc/harp/lib/errors"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/syscheck"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	err := syscheck.CheckOS()
	if err != nil {
		fmt.Println("Please run harp module on Linux or FreeBSD machine.")
		panic(err)
	}

	err = logger.Init()
	if err != nil {
		errors.SetErrLogger(logger.Logger)
		errors.NewHccError(errors.HarpInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	errors.SetErrLogger(logger.Logger)
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Init()

	err = Init()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = Db.Close()
	}()
}
