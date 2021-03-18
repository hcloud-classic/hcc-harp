package logger

import (
	"innogrid.com/hcloud-classic/hcc_errors"
	"testing"
)

func Test_Logger_Prepare(t *testing.T) {
	err := Init()
	if err != nil {
		hcc_errors.SetErrLogger(Logger)
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	hcc_errors.SetErrLogger(Logger)
	defer func() {
		_ = FpLog.Close()
	}()
}
