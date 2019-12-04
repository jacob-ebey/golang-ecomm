package schema

import (
	"github.com/graphql-go/graphql"
)

var MutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"signUp":       SignUpField,
		"signIn":       SignInField,
		"refreshToken": RefreshTokenField,

		"createAddress": CreateAddressField,

		"createProductDraft": CreateProductDraftField,
		"updateProduct":      UpdateProductField,
		"publishProduct":     PublishProductField,
		"addProductImage":    AddProductImageField,
		"removeProductImage": RemoveProductImageField,

		"createProductOption": CreateProductOptionField,
		"removeProductOption": RemoveProductOptionField,

		"createProductVariant":      CreateProductVariantField,
		"updateProductVariant":      UpdateProductVariantField,
		"removeProductVariant":      RemoveProductVariantField,
		"addProductVariantImage":    AddProductVariantImageField,
		"removeProductVariantImage": RemoveProductVariantImageField,
		"createProductPermutations": CreateProductPermutationsField,

		"submitBraintreeTransaction": SubmitBraintreeTransactionField,

		"purchaseShippoLabel": PurchaseShippoLabelField,
	},
})
