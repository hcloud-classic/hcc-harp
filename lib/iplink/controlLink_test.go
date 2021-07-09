package iplink

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/syscheck"
	"innogrid.com/hcloud-classic/hcc_errors"
	"testing"
)

func Test_addHarpInternalDevice(t *testing.T) {
	err := syscheck.CheckRoot()
	if err != nil {
		t.Fatal(err)
	}

	err = syscheck.CheckOS()
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

	err = syscheck.CheckIPLink()
	if err != nil {
		t.Fatal(err)
	}

	err = syscheck.CheckVnStat()
	if err != nil {
		t.Fatal(err)
	}

	config.Init()

	err = mysql.Init()
	if err != nil {
		t.Fatal(err)
	}

	_ = UnsetHarpInternalDevice("192.168.192.168")
	_ = UnsetHarpInternalDevice("999.999.999.999")

	err = addHarpInternalDevice("192.168.192.168")
	if err != nil {
		t.Fatal("Failed to add harp internal device!")
	}

	err = addHarpInternalDevice("192.168.192.168")
	if err != nil {
		t.Log("Tried to add harp internal device with already created")
	}
}

func Test_isHarpInternalDeviceExist(t *testing.T) {
	exist := isHarpInternalDeviceExist("192.168.192.168")
	if !exist {
		t.Fatal("Failed to find harp internal device!")
	}

	exist = isHarpInternalDeviceExist("999.999.999.999")
	if !exist {
		t.Log("Tried to find harp internal device with wrong IP address")
	}
}

func Test_upHarpInternalDevice(t *testing.T) {
	err := upHarpInternalDevice("192.168.192.168")
	if err != nil {
		t.Fatal("Failed to up harp internal device!")
	}

	err = upHarpInternalDevice("255.255.255.255")
	if err != nil {
		t.Log("Tried to up non-exist harp internal device")
	}
}

func Test_downHarpInternalDevice(t *testing.T) {
	err := downHarpInternalDevice("192.168.192.168")
	if err != nil {
		t.Fatal("Failed to down harp internal device!")
	}

	err = downHarpInternalDevice("255.255.255.255")
	if err != nil {
		t.Log("Tried to down non-exist harp internal device")
	}
}

func Test_setIPtoHarpInternalDevice(t *testing.T) {
	err := setIPtoHarpInternalDevice("192.168.192.168", 30)
	if err != nil {
		t.Fatal("Failed to set IP of harp internal device!")
	}

	err = setIPtoHarpInternalDevice("192.168.192.168 aaa", 30)
	if err != nil {
		t.Log("Tried to set IP of harp internal device with wrong arguments")
	}
}

func Test_deleteHarpInternalDevice(t *testing.T) {
	err := deleteHarpInternalDevice("192.168.192.168")
	if err != nil {
		t.Fatal("Failed to delete harp internal device!")
	}

	err = deleteHarpInternalDevice("192.168.192.168")
	if err != nil {
		t.Log("Tried to delete harp internal device with already deleted")
	}
}

func Test_SetHarpInternalDevice(t *testing.T) {
	err := SetHarpInternalDevice("192.168.192.168", "255.255.255.252")
	if err != nil {
		t.Log(err.Error())
		t.Fatal("Failed to set harp internal device!")
	}

	t.Log("Trying to set harp internal device with already created")
	err = SetHarpInternalDevice("192.168.192.168", "255.255.255.252")
	if err != nil {
		t.Log(err.Error())
		t.Fatal("Failed to set harp internal device!")
	}

	err = SetHarpInternalDevice("999.999.999.999", "999.999.999.999")
	if err != nil {
		t.Log("Tried to set harp internal device with wrong arguments")
	}
	_ = UnsetHarpInternalDevice("999.999.999.999")

	err = SetHarpInternalDevice("999999999999.999999999999.999999999999.999999999999", "999.999.999.999")
	if err != nil {
		t.Log("Tried to set harp internal device with wrong arguments")
	}

	err = SetHarpInternalDevice("999.999.999.999", "255.255.255.252")
	if err != nil {
		t.Log("Tried to set harp internal device with wrong arguments")
	}
	_ = UnsetHarpInternalDevice("999.999.999.999")
}

func Test_UnsetHarpInternalDevice(t *testing.T) {
	err := UnsetHarpInternalDevice("192.168.192.168")
	if err != nil {
		t.Log(err.Error())
		t.Fatal("Failed to set harp internal device!")
	}

	err = UnsetHarpInternalDevice("192.168.192.168")
	if err != nil {
		t.Log(err.Error())
		t.Fatal("Tried to unset already deleted harp internal device")
	}
}

func Test_AddOrDeleteIPToHarpExternalDevice(t *testing.T) {
	err := AddOrDeleteIPToHarpExternalDevice("255.255.255.255", "255.255.255.252", true)
	if err != nil {
		t.Fatal("Failed to add harp external device")
	}

	err = AddOrDeleteIPToHarpExternalDevice("255.255.255.255", "255.255.255.252", true)
	if err != nil {
		t.Log("Tried to add harp external device with already created")
	}

	err = AddOrDeleteIPToHarpExternalDevice("255.255.255.255", "255.255.255.252", false)
	if err != nil {
		t.Fatal("Failed to delete harp external device")
	}

	err = AddOrDeleteIPToHarpExternalDevice("255.255.255.255", "999.999.999.999", true)
	if err != nil {
		t.Log("Tried to add harp external device with wrong arguments")
	}

	err = AddOrDeleteIPToHarpExternalDevice("999.999.999.999", "255.255.255.252", true)
	if err != nil {
		t.Log("Tried to add harp external device with wrong arguments")
	}
}
