package schema

import (
	"github.com/graphql-go/graphql"
)

var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"ping": &graphql.Field{
			Type:        graphql.String,
			Description: "Ping",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return "Pong", nil
			},
		},

		"me": MeField,

		"catalog": NewPaginationField(PaginationFieldOpts{
			Type:        ProductType,
			Dataloader:  "products",
			Description: "Paginate through the products.",
		}),
		"product":                         ProductField,
		"productBySlug":                   ProductBySlugField,
		"productVariantsByIds":            ProductVariantsByIdsField,
		"productVariantBySelectedOptions": ProductVariantBySelectedOptionsField,
		"products": NewPaginationField(PaginationFieldOpts{
			Type:        ProductType,
			Dataloader:  "adminProducts",
			Description: "Paginate through the products. This is the admin entry, use catalog for public access.",
			AuthRole:    "ADMIN",
		}),

		"subtotal":            SubtotalField,
		"taxes":               TaxesField,
		"shippingEstimations": ShippingEstimationsField,

		"braintreeClientToken": BraintreeClientTokenField,

		"transaction": TransactionField,
		"transactions": NewPaginationField(PaginationFieldOpts{
			Type:        TransactionType,
			Dataloader:  "transactions",
			Description: "Paginate through the transactions.",
			AuthRole:    "ADMIN",
		}),
	},
})
