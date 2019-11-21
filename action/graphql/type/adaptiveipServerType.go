package graphqlType

import "github.com/graphql-go/graphql"

// AdaptiveIPServerType : Graphql object type of AdaptiveIPServer
var AdaptiveIPServerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AdaptiveIPServer",
		Fields: graphql.Fields{
			"adaptiveip_uuid": &graphql.Field{
				Type: graphql.String,
			},
			"server_uuid": &graphql.Field{
				Type: graphql.String,
			},
			"public_ip": &graphql.Field{
				Type: graphql.String,
			},
			"private_ip": &graphql.Field{
				Type: graphql.String,
			},
			"private_gateway": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
