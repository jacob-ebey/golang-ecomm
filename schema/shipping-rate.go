package schema

import (
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	core "github.com/jacob-ebey/graphql-core"
)

var ShippingRateType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ShippingRate",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"price": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"service": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"carrier": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"durationTerms": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var ShippingEstimationsField = &graphql.Field{
	Type:        graphql.NewList(ShippingRateType),
	Description: "Get shipping estimations for the provided variants and quantities.",
	Args: graphql.FieldConfigArgument{
		"address": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(AddressInputSchema),
		},
		"variants": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(
				graphql.NewNonNull(CartInputSchema),
			)),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		shippingEstimations := params.Context.Value("shippingEstimations").(*dataloader.Loader)

		key := dataloaders.ShippingEstimationKey{}
		if err := ConvertObject(params.Args, &key); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not convert arguments.",
				InternalError: err,
			}
		}

		thunk := shippingEstimations.Load(params.Context, key)

		return func() (interface{}, error) {
			return thunk()
		}, nil
	},
}
