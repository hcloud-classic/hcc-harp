package main

import (
	"fmt"
	"innogrid.com/hcloud-classic/hcc_errors"
	"hcc/harp/action/grpc/client"
	"hcc/harp/action/grpc/server"
	"hcc/harp/lib/adaptiveip"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/pid"
	"hcc/harp/lib/syscheck"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func init() {
	err, harpRunning, harpPID := pid.IsHarpRunning()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if harpRunning {
		fmt.Println("harp is already running. (PID: " + strconv.Itoa(harpPID) + ")")
		os.Exit(1)
	}
	err = pid.WriteHarpPID()
	if err != nil {
		_ = pid.DeleteHarpPID()
		fmt.Println(err)
		panic(err)
	}

	err = syscheck.CheckOS()
	if err != nil {
		fmt.Println("Please run harp module on Linux or FreeBSD machine.")
		_ = pid.DeleteHarpPID()
		panic(err)
	}

	err = syscheck.CheckRoot()
	if err != nil {
		_ = pid.DeleteHarpPID()
		panic(err)
	}

	err = logger.Init()
	if err != nil {
		hcc_errors.SetErrLogger(logger.Logger)
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
		_ = pid.DeleteHarpPID()
	}
	hcc_errors.SetErrLogger(logger.Logger)

	err = syscheck.CheckFirewall()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "syscheck.CheckFirewall(): "+err.Error()).Fatal()
		_ = pid.DeleteHarpPID()
	}

	config.Init()

	err = mysql.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "mysql.Init(): "+err.Error()).Fatal()
		_ = pid.DeleteHarpPID()
	}

	err = client.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "client.Init(): "+err.Error()).Fatal()
		_ = pid.DeleteHarpPID()
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.ExternalIfaceName)
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "syscheck.CheckIfaceExist(): "+err.Error()).Fatal()
		_ = pid.DeleteHarpPID()
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.InternalIfaceName)
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "syscheck.CheckIfaceExist(): "+err.Error()).Fatal()
		_ = pid.DeleteHarpPID()
	}

	err = dhcpd.CheckLocalDHCPDConfig()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "dhcpd.CheckLocalDHCPDConfig(): "+err.Error()).Fatal()
		_ = pid.DeleteHarpPID()
	}

	err = adaptiveip.LoadFirewall()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "adaptiveip.LoadFirewall(): "+err.Error()).Fatal()
		_ = pid.DeleteHarpPID()
	}
}

func end() {
	client.End()
	mysql.End()
	logger.End()
	_ = pid.DeleteHarpPID()
}

func main() {
	// Catch the exit signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		end()
		fmt.Println("Exiting harp module...")
		os.Exit(0)
	}()

	server.Init()
}
