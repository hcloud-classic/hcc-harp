package graphqltype

import "github.com/graphql-go/graphql"

// AdaptiveIPNumType : Graphql object type of AdaptiveIPNumType
var AdaptiveIPNumType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AdaptiveIPNum",
		Fields: graphql.Fields{
			"number": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
