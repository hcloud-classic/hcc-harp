package pid

import (
	"hcc/harp/lib/fileutil"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

var harpPIDFileLocation = "/var/run"
var harpPIDFile = "/var/run/harp.pid"

// IsHarpRunning : Check if harp is running
func IsHarpRunning() (err error, running bool, pid int) {
	if _, err := os.Stat(harpPIDFile); os.IsNotExist(err) {
		return nil, false, 0
	}

	pidStr, err := ioutil.ReadFile(harpPIDFile)
	if err != nil {
		return err, false, 0
	}

	harpPID, _ := strconv.Atoi(string(pidStr))

	proc, err := os.FindProcess(harpPID)
	if err != nil {
		return err, false, 0
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return nil, true, harpPID
	}

	return nil, false, 0
}

// WriteHarpPID : Write harp PID to "/var/run/harp.pid"
func WriteHarpPID() error {
	pid := os.Getpid()

	err := fileutil.CreateDirIfNotExist(harpPIDFileLocation)
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(harpPIDFile, strconv.Itoa(pid))
	if err != nil {
		return err
	}

	return nil
}

// DeleteHarpPID : Delete the harp PID file
func DeleteHarpPID() error {
	err := fileutil.DeleteFile(harpPIDFile)
	if err != nil {
		return err
	}

	return nil
}
