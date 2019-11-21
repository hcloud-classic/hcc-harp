package graphqlType

import "github.com/graphql-go/graphql"

// AdaptiveIPType : Graphql object type of AdaptiveIP
var AdaptiveIPType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AdaptiveIP",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.String,
			},
			"network_address": &graphql.Field{
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
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
