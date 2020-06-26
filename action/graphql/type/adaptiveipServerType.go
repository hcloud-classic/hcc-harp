package graphqltype

import "github.com/graphql-go/graphql"

// AdaptiveIPServerType : Graphql object type of AdaptiveIPServer
var AdaptiveIPServerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AdaptiveIPServer",
		Fields: graphql.Fields{
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
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
