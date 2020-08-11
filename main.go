package main

import (
	"hcc/harp/driver/grpccli"
	"hcc/harp/driver/grpcsrv"
	"hcc/harp/lib/adaptiveip"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/logger"
	"hcc/harp/lib/mysql"
	"hcc/harp/lib/pf"
	"hcc/harp/lib/syscheck"
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

func main() {
	defer func() {
		grpccli.CleanGRPCClient()
		mysql.End()
		logger.End()
	}()

	grpcsrv.Init()
}
