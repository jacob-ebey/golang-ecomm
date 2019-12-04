package schema

import (
	"github.com/graphql-go/graphql"
)

var TaxRateType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TaxRate",
		Fields: graphql.Fields{
			"rate": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var TaxesType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TaxRates",
		Fields: graphql.Fields{
			"totalRate": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"rates": &graphql.Field{
				Type: graphql.NewList(TaxRateType),
			},
		},
	},
)
