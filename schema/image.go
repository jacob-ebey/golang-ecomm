package schema

import "github.com/graphql-go/graphql"

var ImageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Image",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"raw": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"thumbnail": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"height600": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})
