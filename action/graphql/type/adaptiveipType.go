package graphqltype

import "github.com/graphql-go/graphql"

// AdaptiveIPType : Graphql object type of AdaptiveIP
var AdaptiveIPType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AdaptiveIP",
		Fields: graphql.Fields{
			"ext_iface_ip_address": &graphql.Field{
				Type: graphql.String,
			},
			"netmask": &graphql.Field{
				Type: graphql.String,
			},
			"gateway": &graphql.Field{
				Type: graphql.String,
			},
			"start_ip_address": &graphql.Field{
				Type: graphql.String,
			},
			"end_ip_address": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// AdaptiveIPAvailableIPListType : Graphql object type of AdaptiveIPAvailableIPList
var AdaptiveIPAvailableIPListType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AdaptiveIP",
		Fields: graphql.Fields{
			"available_ip_list": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)
