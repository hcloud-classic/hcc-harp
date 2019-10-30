package graphql

import (
	"github.com/graphql-go/graphql"
	"hcc/harp/dao"
	"hcc/harp/lib/config"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/logger"
	"strings"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		// subnet DB
		"create_subnet": &graphql.Field{
			Type:        subnetType,
			Description: "Create new subnet",
			Args: graphql.FieldConfigArgument{
				"network_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"netmask": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"gateway": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"next_server": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"name_server": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"domain_name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"server_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"leader_node_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"os": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"subnet_name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return dao.CreateSubnet(params.Args)
			},
		},
		"update_subnet": &graphql.Field{
			Type:        subnetType,
			Description: "Update subnet",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"network_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"netmask": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"gateway": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"next_server": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"name_server": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"domain_name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"server_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"leader_node_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"os": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"subnet_name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_subnet")
				subnet, err := dao.UpdateSubnet(params.Args)
				if err != nil {
					logger.Logger.Print(err)
					return nil, err
				}

				return subnet, nil
			},
		},
		"delete_subnet": &graphql.Field{
			Type:        subnetType,
			Description: "Delete subnet by uuid",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: delete_subnet")
				return dao.DeleteSubnet(params.Args)
			},
		},

		"create_dhcpd_conf": &graphql.Field{
			Type:        graphql.String,
			Description: "Create new dhcpd config",
			Args: graphql.FieldConfigArgument{
				"subnet_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"node_uuids": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				subnetUUID := params.Args["subnet_uuid"].(string)
				nodeUUIDs := params.Args["node_uuids"].(string)

				nodeUUIDsParts := strings.Split(nodeUUIDs, ",")

				err := dhcpd.CreateConfig(subnetUUID, nodeUUIDsParts)
				if err != nil {
					return nil, err
				}

				err = dhcpd.RestartDHCPDServer()
				if err != nil {
					logger.Logger.Println("Failed to restart dhcpd server (" + config.DHCPD.LocalDHCPDServiceName +")")
					return nil, err
				}

				return "CreateDHCPDConfig: succeed", nil
			},
		},
	},
})
