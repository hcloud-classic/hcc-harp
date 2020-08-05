package graphql

import (
	"hcc/harp/dao"
	"hcc/harp/lib/logger"

	"github.com/graphql-go/graphql"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"subnet": &graphql.Field{
				Type:        subnetType,
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
				Type:        graphql.NewList(subnetType),
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
		},
	})
