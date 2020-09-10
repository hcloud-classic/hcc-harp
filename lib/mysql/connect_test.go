package mysql

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/errors"
	"hcc/harp/lib/logger"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	err := logger.Init()
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
