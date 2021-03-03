package main

import (
	"fmt"
	"github.com/hcloud-classic/hcc_errors"
	"hcc/harp/action/grpc/client"
	"hcc/harp/action/grpc/server"
	"hcc/harp/lib/adaptiveip"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/syscheck"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	err := syscheck.CheckOS()
	if err != nil {
		fmt.Println("Please run harp module on Linux or FreeBSD machine.")
		panic(err)
	}

	err = syscheck.CheckRoot()
	if err != nil {
		panic(err)
	}

	err = logger.Init()
	if err != nil {
		hcc_errors.SetErrLogger(logger.Logger)
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	hcc_errors.SetErrLogger(logger.Logger)

	err = syscheck.CheckFirewall()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "syscheck.CheckFirewall(): "+err.Error()).Fatal()
	}

	config.Init()

	err = mysql.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "mysql.Init(): "+err.Error()).Fatal()
	}

	err = client.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "client.Init(): "+err.Error()).Fatal()
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.ExternalIfaceName)
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "syscheck.CheckIfaceExist(): "+err.Error()).Fatal()
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.InternalIfaceName)
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "syscheck.CheckIfaceExist(): "+err.Error()).Fatal()
	}

	err = dhcpd.CheckLocalDHCPDConfig()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "dhcpd.CheckLocalDHCPDConfig(): "+err.Error()).Fatal()
	}

	err = adaptiveip.LoadFirewall()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "adaptiveip.LoadFirewall(): "+err.Error()).Fatal()
	}
}

func end() {
	client.End()
	mysql.End()
	logger.End()
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
