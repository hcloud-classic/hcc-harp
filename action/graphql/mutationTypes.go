package graphql

import (
	"github.com/graphql-go/graphql"
	graphqlType "hcc/harp/action/graphql/type"
	"hcc/harp/dao"
	"hcc/harp/driver"
	"hcc/harp/lib/logger"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		// subnet DB
		"create_subnet": &graphql.Field{
			Type:        graphqlType.SubnetType,
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
				logger.Logger.Println("Resolving: create_subnet")
				return dao.CreateSubnet(params.Args)
			},
		},
		"update_subnet": &graphql.Field{
			Type:        graphqlType.SubnetType,
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
				return dao.UpdateSubnet(params.Args)
			},
		},
		"delete_subnet": &graphql.Field{
			Type:        graphqlType.SubnetType,
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
		// dhcpd
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
				logger.Logger.Println("Resolving: create_dhcpd_conf")
				return driver.CreateDHCPDConfig(params)
			},
		},
		// adaptive IP
		"create_adaptiveip": &graphql.Field{
			Type:        graphqlType.SubnetType,
			Description: "Create new adaptiveip",
			Args: graphql.FieldConfigArgument{
				"network_address": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"netmask": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"gateway": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"start_ip_address": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"end_ip_address": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: create_adaptiveip")
				return dao.CreateAdaptiveIP(params.Args)
			},
		},
		"update_adaptiveip": &graphql.Field{
			Type:        graphqlType.SubnetType,
			Description: "Update adaptiveip",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"network_address": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"netmask": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"gateway": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"start_ip_address": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"end_ip_address": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_adaptiveip")
				return dao.UpdateAdaptiveIP(params.Args)
			},
		},
		"delete_adaptiveip": &graphql.Field{
			Type:        graphqlType.SubnetType,
			Description: "Delete adaptiveip by uuid",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: delete_subnet")
				return dao.DeleteAdaptiveIP(params.Args)
			},
		},
		"create_adaptiveip_server": &graphql.Field{
			Type:        graphqlType.SubnetType,
			Description: "Create new adaptiveip_server",
			Args: graphql.FieldConfigArgument{
				"server_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"public_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"private_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"private_gateway": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: create_adaptiveip_server")
				return dao.CreateAdaptiveIPServer(params.Args)
			},
		},
		"delete_adaptiveip_server": &graphql.Field{
			Type:        graphqlType.SubnetType,
			Description: "Delete adaptiveip_server by server_uuid",
			Args: graphql.FieldConfigArgument{
				"server_uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: delete_adaptiveip_server")
				return dao.DeleteAdaptiveIPServer(params.Args)
			},
		},
	},
})
