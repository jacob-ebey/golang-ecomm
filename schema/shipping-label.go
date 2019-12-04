package schema

import (
	"github.com/graphql-go/graphql"
)

var ShippingLabelType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ShippingLabel",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"labelUrl": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)
