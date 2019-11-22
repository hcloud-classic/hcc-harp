package graphql

import (
	"github.com/graphql-go/graphql"
	graphqlType "hcc/harp/action/graphql/type"
	"hcc/harp/dao"
	"hcc/harp/lib/logger"
	"hcc/harp/model"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// subnet DB
			"subnet": &graphql.Field{
				Type:        graphqlType.SubnetType,
				Description: "Get subnet by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: subnet")
					return dao.ReadSubnet(params.Args)
				},
			},
			"list_subnet": &graphql.Field{
				Type:        graphql.NewList(graphqlType.SubnetType),
				Description: "Get subnet list",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
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
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"row": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: list_subnet")
					return dao.ReadSubnetList(params.Args)
				},
			},
			"all_subnet": &graphql.Field{
				Type:        graphql.NewList(graphqlType.SubnetType),
				Description: "Get all subnet list",
				Args: graphql.FieldConfigArgument{
					"row": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: all_subnet")
					return dao.ReadSubnetAll(params.Args)
				},
			},
			"num_subnet": &graphql.Field{
				Type:        graphqlType.SubnetNumType,
				Description: "Get the number of subnet",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: num_subnet")
					var subnetNum model.SubnetNum
					var err error
					subnetNum, err = dao.ReadSubnetNum()

					return subnetNum, err
				},
			},
			// adaptive IP
			"adaptiveip": &graphql.Field{
				Type:        graphqlType.AdaptiveIPType,
				Description: "Get adaptiveip by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: adaptiveip")
					return dao.ReadAdaptiveIP(params.Args["uuid"].(string))
				},
			},
			"list_adaptiveip": &graphql.Field{
				Type:        graphql.NewList(graphqlType.AdaptiveIPType),
				Description: "Get adaptiveip list",
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
					logger.Logger.Println("Resolving: list_adaptiveip")
					return dao.ReadAdaptiveIPList(params.Args)
				},
			},
			"all_adaptiveip": &graphql.Field{
				Type:        graphql.NewList(graphqlType.AdaptiveIPType),
				Description: "Get all adaptiveip list",
				Args: graphql.FieldConfigArgument{
					"row": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: all_adaptiveip")
					return dao.ReadAdaptiveIPAll(params.Args)
				},
			},
			"num_adaptiveip": &graphql.Field{
				Type:        graphqlType.AdaptiveIPNumType,
				Description: "Get the number of adaptiveip",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: num_adaptiveip")
					var adaptiveIPNum model.AdaptiveIPNum
					var err error
					adaptiveIPNum, err = dao.ReadAdaptiveIPNum()

					return adaptiveIPNum, err
				},
			},
			"adaptiveip_server": &graphql.Field{
				Type:        graphqlType.AdaptiveIPServerType,
				Description: "Get adaptiveip by uuid",
				Args: graphql.FieldConfigArgument{
					"adaptiveip_uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"server_uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: adaptiveip_server")
					return dao.ReadAdaptiveIPServer(params.Args)
				},
			},
			"list_adaptiveip_server": &graphql.Field{
				Type:        graphql.NewList(graphqlType.AdaptiveIPServerType),
				Description: "Get adaptiveip_server list",
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
					logger.Logger.Println("Resolving: list_adaptiveip_server")
					return dao.ReadAdaptiveIPServerList(params.Args)
				},
			},
			"all_adaptiveip_server": &graphql.Field{
				Type:        graphql.NewList(graphqlType.AdaptiveIPServerType),
				Description: "Get all adaptiveip_server list",
				Args: graphql.FieldConfigArgument{
					"row": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: all_adaptiveip_server")
					return dao.ReadAdaptiveIPServerAll(params.Args)
				},
			},
			"num_adaptiveip_server": &graphql.Field{
				Type:        graphqlType.AdaptiveIPServerNumType,
				Description: "Get the number of adaptiveip_server",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: num_adaptiveip_server")
					var adaptiveIPServerNum model.AdaptiveIPServerNum
					var err error
					adaptiveIPServerNum, err = dao.ReadAdaptiveIPServerNum(params.Args)

					return adaptiveIPServerNum, err
				},
			},
		},
	})
