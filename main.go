package main

import (
	"fmt"
	"hcc/harp/driver/grpccli"
	"hcc/harp/driver/grpcsrv"
	"hcc/harp/lib/adaptiveip"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
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

	err = syscheck.CheckArpingCommand()
	if err != nil {
		panic(err)
	}

	err = logger.Init()
	if err != nil {
		panic(err)
	}

	config.Parser()

	err = mysql.Init()
	if err != nil {
		panic(err)
	}

	go grpcsrv.Init()
	err = grpccli.InitGRPCClient()
	if err != nil {
		panic(err)
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.ExternalIfaceName)
	if err != nil {
		panic(err)
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.InternalIfaceName)
	if err != nil {
		panic(err)
	}

	err = dhcpd.CheckLocalDHCPDConfig()
	if err != nil {
		panic(err)
	}

	err = pf.PreparePFConfigFiles()
	if err != nil {
		panic(err)
	}

	err = adaptiveip.LoadHarpPFRules()
	if err != nil {
		panic(err)
	}
}

func end(){
	grpccli.CleanGRPCClient()
	mysql.End()
	logger.End()
}

func main() {
	// Catch the exit signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func(){
		<- sigChan
		end()
		fmt.Println("Exiting harp module...")
		os.Exit(0)
	}()

	// Prevent to exit main thread
	mainChan := make(chan bool)
	<-mainChan
}
