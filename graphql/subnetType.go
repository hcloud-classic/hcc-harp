package graphql

import "github.com/graphql-go/graphql"

var subnetType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Subnet",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"ip": &graphql.Field{
				Type: graphql.String,
			},
			"netmask": &graphql.Field{
				Type: graphql.String,
			},
			"os": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
