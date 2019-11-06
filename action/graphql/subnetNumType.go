package graphql

import "github.com/graphql-go/graphql"

var subnetNum = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "SubnetNum",
		Fields: graphql.Fields{
			"number": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
