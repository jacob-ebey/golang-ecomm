package schema

import (
	"github.com/graphql-go/graphql"
)

var AddressType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Address",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"line1": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"line2": &graphql.Field{
			Type: graphql.String,
		},
		"line3": &graphql.Field{
			Type: graphql.String,
		},
		"city": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"region": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"postalCode": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"country": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var AddressInputSchema = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "AddressInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"line1": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"line2": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"line3": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"city": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"region": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"postalCode": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"country": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})
