package schema

import (
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	"github.com/jacob-ebey/golang-ecomm/db"
	core "github.com/jacob-ebey/graphql-core"
)

var CartInputSchema = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CartInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"variantId": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"quantity": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
})

var SubtotalField = &graphql.Field{
	Type:        graphql.Int,
	Description: "The subtotal for the provided variants and their quantities.",
	Args: graphql.FieldConfigArgument{
		"variants": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(
				graphql.NewNonNull(CartInputSchema),
			)),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		subtotal := params.Context.Value("subtotal").(*dataloader.Loader)

		cart := dataloaders.CartKey{}
		if err := ConvertObject(params.Args["variants"], &cart); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not convert variants argument.",
				InternalError: err,
			}
		}

		thunk := subtotal.Load(params.Context, cart)

		return func() (interface{}, error) {
			return thunk()
		}, nil
	},
}

var TaxesField = &graphql.Field{
	Type:        TaxesType,
	Description: "Gets tax rates for a given address.",
	Args: graphql.FieldConfigArgument{
		"address": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(AddressInputSchema),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		taxesLoader := params.Context.Value("taxes").(*dataloader.Loader)

		address := db.Address{}
		if err := ConvertObject(params.Args["address"], &address); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not convert address argument.",
				InternalError: err,
			}
		}

		thunk := taxesLoader.Load(params.Context, address)

		return func() (interface{}, error) {
			return thunk()
		}, nil
	},
}
