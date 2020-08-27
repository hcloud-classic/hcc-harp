package logger

import (
	"hcc/harp/lib/syscheck"
	"testing"
)

func Test_Logger_Prepare(t *testing.T) {
	err := syscheck.CheckRoot()
	if err != nil {
		t.Fatal("Failed to get root permission!")
	}

	err = Init()
	if err != nil {
		t.Fatal("Failed to prepare logger!")
	}
	defer func() {
		_ = FpLog.Close()
	}()
}
