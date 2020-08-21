package main

import (
	"fmt"
	"hcc/harp/action/grpc/client"
	"hcc/harp/action/grpc/server"
	"hcc/harp/lib/adaptiveip"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/pf"
	"hcc/harp/lib/syscheck"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	err := syscheck.CheckRoot()
	if err != nil {
		log.Fatalf("syscheck.CheckRoot(): %v", err.Error())
	}

	err = syscheck.CheckArpingCommand()
	if err != nil {
		log.Fatalf("syscheck.CheckArpingCommand(): %v", err.Error())
	}

	err = logger.Init()
	if err != nil {
		log.Fatalf("logger.Init(): %v", err.Error())
	}

	config.Init()

	err = mysql.Init()
	if err != nil {
		logger.Logger.Fatalf("mysql.Init(): %v", err.Error())
	}

	err = client.Init()
	if err != nil {
		logger.Logger.Fatalf("client.Init(): %v", err.Error())
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.ExternalIfaceName)
	if err != nil {
		logger.Logger.Fatalf("syscheck.CheckIfaceExist(): %v", err.Error())
	}

	_, err = syscheck.CheckIfaceExist(config.AdaptiveIP.InternalIfaceName)
	if err != nil {
		logger.Logger.Fatalf("syscheck.CheckIfaceExist(): %v", err.Error())
	}

	err = dhcpd.CheckLocalDHCPDConfig()
	if err != nil {
		logger.Logger.Fatalf("dhcpd.CheckLocalDHCPDConfig(): %v", err.Error())
	}

	err = pf.PreparePFConfigFiles()
	if err != nil {
		logger.Logger.Fatalf("pf.PreparePFConfigFiles(): %v", err.Error())
	}

	err = adaptiveip.LoadHarpPFRules()
	if err != nil {
		logger.Logger.Fatalf("adaptiveip.LoadHarpPFRules(): %v", err.Error())
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
