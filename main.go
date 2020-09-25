package main

import (
	"fmt"
	"hcc/harp/action/grpc/client"
	"hcc/harp/action/grpc/server"
	"hcc/harp/lib/adaptiveip"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/errors"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/pf"
	"hcc/harp/lib/syscheck"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	err := syscheck.CheckRoot()
	if err != nil {
		panic(err)
	}

	err = logger.Init()
	if err != nil {
		errors.SetErrLogger(logger.Logger)
		errors.NewHccError(errors.HarpInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	errors.SetErrLogger(logger.Logger)

	err = syscheck.CheckArpingCommand()
	if err != nil {
		errors.NewHccError(errors.HarpInternalInitFail, "syscheck.CheckArpingCommand(): "+err.Error()).Fatal()
	}

	config.Init()

	err = mysql.Init()
	if err != nil {
		errors.NewHccError(errors.PiccoloInternalInitFail, "mysql.Init(): "+err.Error()).Fatal()
	}

	err = client.Init()
	if err != nil {
		errors.NewHccError(errors.PiccoloInternalInitFail, "client.Init(): "+err.Error()).Fatal()
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.ExternalIfaceName)
	if err != nil {
		errors.NewHccError(errors.PiccoloInternalInitFail, "syscheck.CheckIfaceExist(): "+err.Error()).Fatal()
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.InternalIfaceName)
	if err != nil {
		errors.NewHccError(errors.PiccoloInternalInitFail, "syscheck.CheckIfaceExist(): "+err.Error()).Fatal()
	}

	err = dhcpd.CheckLocalDHCPDConfig()
	if err != nil {
		errors.NewHccError(errors.PiccoloInternalInitFail, "dhcpd.CheckLocalDHCPDConfig(): "+err.Error()).Fatal()
	}

	err = pf.PreparePFConfigFiles()
	if err != nil {
		errors.NewHccError(errors.PiccoloInternalInitFail, "pf.PreparePFConfigFiles(): "+err.Error()).Fatal()
	}

	err = adaptiveip.LoadHarpPFRules()
	if err != nil {
		errors.NewHccError(errors.PiccoloInternalInitFail, "adaptiveip.LoadHarpPFRules(): "+err.Error()).Fatal()
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
