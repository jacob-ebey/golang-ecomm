package schema

import (
	"math"
	"strconv"

	"github.com/jacob-ebey/go-shippo/client"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/jacob-ebey/golang-ecomm/auth"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	"github.com/jacob-ebey/golang-ecomm/db"
	core "github.com/jacob-ebey/graphql-core"
)

var TransactionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Transaction",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"subtotal": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"taxes": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"shipping": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"total": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"shippoRateId": &graphql.Field{
			Type: graphql.String,
		},
		"shippingLabel": &graphql.Field{
			Type: ShippingLabelType,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				shippoClient := params.Context.Value("shippo").(*client.Client)

				transaction := params.Source.(*db.Transaction)

				if transaction.ShippoTransactionID == "" {
					return nil, nil
				}

				label, err := shippoClient.RetrieveTransaction(transaction.ShippoTransactionID)
				if err != nil {
					return nil, &core.WrappedError{
						Message:       "Could not retrieve shipping label.",
						InternalError: err,
					}
				}

				return map[string]interface{}{
					"id":       label.ObjectID,
					"labelUrl": label.LabelURL,
				}, nil
			},
		},
		"shippingEstimation": &graphql.Field{
			Type: ShippingRateType,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				shippoClient := params.Context.Value("shippo").(*client.Client)

				transaction := params.Source.(*db.Transaction)

				if transaction.ShippoRateID == "" {
					return nil, nil
				}

				rate, err := shippoClient.RetrieveRate(transaction.ShippoRateID)
				if err != nil || rate == nil {
					return nil, &core.WrappedError{
						Message:       "Could not retrieve shipping rate.",
						InternalError: err,
					}
				}

				amount, err := strconv.ParseFloat(rate.Amount, 64)
				if err != nil {
					return nil, &core.WrappedError{
						Message:       "Could not convert estimation price.",
						InternalError: err,
					}
				}
				price := int(math.Round(amount * 100))

				return &dataloaders.ShippingEstimation{
					ID:            rate.ObjectID,
					Price:         price,
					Carrier:       rate.Provider,
					Service:       rate.ServiceLevel.Name,
					DurationTerms: rate.DurationTerms,
				}, nil
			},
		},
		"lineItems": &graphql.Field{
			Type: graphql.NewList(graphql.NewNonNull(ReceiptLineItemType)),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				transactionLineItems := params.Context.Value("transactionLineItems").(*dataloader.Loader)

				transaction := params.Source.(*db.Transaction)

				thunk := transactionLineItems.Load(params.Context, dataloaders.IntKey(transaction.ID))

				return func() (interface{}, error) {
					return thunk()
				}, nil
			},
		},
		"billingAddress": &graphql.Field{
			Type: AddressType,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				transactionAddresses := params.Context.Value("transactionAddresses").(*dataloader.Loader)

				transaction := params.Source.(*db.Transaction)

				thunk := transactionAddresses.Load(params.Context, dataloaders.IntKey(transaction.ID))

				return func() (interface{}, error) {
					res, err := thunk()

					if err != nil {
						return nil, err
					}

					return res.(*db.TransactionAddressInfo).BillingAddress, nil
				}, nil
			},
		},
		"shippingAddress": &graphql.Field{
			Type: AddressType,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				transactionAddresses := params.Context.Value("transactionAddresses").(*dataloader.Loader)

				transaction := params.Source.(*db.Transaction)

				thunk := transactionAddresses.Load(params.Context, dataloaders.IntKey(transaction.ID))

				return func() (interface{}, error) {
					res, err := thunk()

					if err != nil {
						return nil, err
					}

					return res.(*db.TransactionAddressInfo).ShippingAddress, nil
				}, nil
			},
		},
	},
})

var TransactionField = &graphql.Field{
	Type:        TransactionType,
	Description: "Get a transaction by ID.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		transaction := params.Context.Value("transaction").(*dataloader.Loader)
		claims := params.Context.Value("claims").(*auth.Claims)

		if claims == nil {
			return nil, auth.NotAuthenticatedError
		}

		if claims.Role != "ADMIN" {
			return nil, auth.NotAuthorizedError
		}

		id := params.Args["id"].(int)

		thunk := transaction.Load(params.Context, dataloaders.IntKey(id))

		return func() (interface{}, error) {
			return thunk()
		}, nil
	},
}
