package logger

import (
	"hcc/harp/lib/errors"
	"testing"
)

func Test_Logger_Prepare(t *testing.T) {
	err := Init()
	if err != nil {
		errors.SetErrLogger(Logger)
		errors.NewHccError(errors.HarpInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	errors.SetErrLogger(Logger)
	defer func() {
		_ = FpLog.Close()
	}()
}
