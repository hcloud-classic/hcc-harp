package driver

import (
	"github.com/graphql-go/graphql"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/logger"
	"strings"
)

// CreateDHCPDConfig : Do dhcpd config file creation works
func CreateDHCPDConfig(params graphql.ResolveParams) (interface{}, error) {
	subnetUUID := params.Args["subnet_uuid"].(string)
	nodeUUIDs := params.Args["node_uuids"].(string)

	nodeUUIDsParts := strings.Split(nodeUUIDs, ",")

	err := dhcpd.CreateConfig(subnetUUID, nodeUUIDsParts)
	if err != nil {
		return nil, err
	}

	err = dhcpd.UpdateHarpDHCPDConfig()
	if err != nil {
		return nil, err
	}

	err = dhcpd.RestartDHCPDServer()
	if err != nil {
		logger.Logger.Println("Failed to restart dhcpd server (" + config.DHCPD.LocalDHCPDServiceName + ")")
		return nil, err
	}

	return "CreateDHCPDConfig: succeed", nil
}