package schema

import (
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	"github.com/jacob-ebey/golang-ecomm/db"
)

var ReceiptLineItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ReceiptLineItem",
	Fields: graphql.Fields{
		"price": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"quantity": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"variant": &graphql.Field{
			Type: ProductVariantType,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				productVariant := params.Context.Value("productVariant").(*dataloader.Loader)

				lineItem := params.Source.(*db.TransactionLineItem)

				thunk := productVariant.Load(params.Context, dataloaders.IntKey(lineItem.ProductVariantID))

				return func() (interface{}, error) {
					return thunk()
				}, nil
			},
		},
	},
})

var ReceiptType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Receipt",
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
