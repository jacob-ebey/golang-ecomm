package schema

import (
	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/jacob-ebey/golang-ecomm/auth"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	"github.com/jacob-ebey/golang-ecomm/db"
	"github.com/jacob-ebey/golang-ecomm/services"
	core "github.com/jacob-ebey/graphql-core"
)

var MeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Me",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"email": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"addresses": &graphql.Field{
			Type: graphql.NewList(graphql.NewNonNull(AddressType)),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userAddresses := params.Context.Value("userAddresses").(*dataloader.Loader)

				me := params.Source.(*auth.Claims)

				thunk := userAddresses.Load(params.Context, dataloaders.IntKey(me.ID))

				return func() (interface{}, error) {
					return thunk()
				}, nil
			},
		},
		"receipts": &graphql.Field{
			Type: graphql.NewList(graphql.NewNonNull(ReceiptType)),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userTransactions := params.Context.Value("userTransactions").(*dataloader.Loader)

				me := params.Source.(*auth.Claims)

				thunk := userTransactions.Load(params.Context, dataloaders.IntKey(me.ID))

				return func() (interface{}, error) {
					return thunk()
				}, nil
			},
		},
	},
})

var MeField = &graphql.Field{
	Type:        MeType,
	Description: "Your information.",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		claims := params.Context.Value("claims").(*auth.Claims)

		if claims == nil {
			return nil, auth.NotAuthenticatedError
		}

		return claims, nil
	},
}

var CreateAddressField = &graphql.Field{
	Type:        AddressType,
	Description: "Create a new address for the logged in user.",
	Args: graphql.FieldConfigArgument{
		"address": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(AddressInputSchema),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		database := params.Context.Value("database").(*pg.DB)
		addressValidator := params.Context.Value("addressValidator").(services.AddressValidator)

		claims := params.Context.Value("claims").(*auth.Claims)
		if claims == nil {
			return nil, auth.NotAuthenticatedError
		}

		address := db.Address{}
		if err := ConvertObject(params.Args["address"], &address); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not convert address argument.",
				InternalError: err,
			}
		}

		address.UserID = claims.ID

		_, err := addressValidator.ValidateAddress(params.Context, address)
		if err != nil {
			return nil, err
		}

		if err := database.Insert(&address); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create address.",
				InternalError: err,
			}
		}

		return &address, nil
	},
}

var DeleteAddressField = &graphql.Field{
	Type:        AddressType,
	Description: "Delete an address for the logged in user.",
	Args: graphql.FieldConfigArgument{
		"addressId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		database := params.Context.Value("database").(*pg.DB)
		claims := params.Context.Value("claims").(*auth.Claims)

		addressID := params.Args["addressId"].(int)

		if claims == nil {
			return nil, auth.NotAuthenticatedError
		}

		address := db.Address{}
		if err := database.
			Model(&address).
			Where("address.id = ?", addressID).
			Where("address.user_id = ?", claims.ID).
			Select(); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not find address for user to delete.",
				InternalError: err,
			}
		}

		if result, err := database.
			Model(&address).
			Where("address.id = ?", addressID).
			Where("address.user_id = ?", claims.ID).
			Delete(); err != nil || result.RowsAffected() == 0 {
			return nil, &core.WrappedError{
				Message:       "Could not delete address.",
				InternalError: err,
			}
		}

		return &address, nil
	},
}
