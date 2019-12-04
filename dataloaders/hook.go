package dataloaders

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"

	core "github.com/jacob-ebey/graphql-core"
)

type hooksDataloaderFunc func(ctx context.Context, req core.GraphQLRequest) context.Context

func (hook hooksDataloaderFunc) PreExecute(ctx context.Context, req core.GraphQLRequest) context.Context {
	return hook(ctx, req)
}

func (hook hooksDataloaderFunc) PostExecute(ctx context.Context, req core.GraphQLRequest, res *graphql.Result) {
	loader := ctx.Value("user").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("userAddresses").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("userTransactions").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("products").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("adminProducts").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("product").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("productBySlug").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("productOptions").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("productImages").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("productOptionValues").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("productVariants").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("productVariantBySelectedOptions").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("productVariant").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("productVariantOptions").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("productVariantImages").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("transaction").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("transactions").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("transactionAddresses").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("transactionLineItems").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("subtotal").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("taxes").(*dataloader.Loader)
	loader.ClearAll()
	loader = ctx.Value("shippingEstimations").(*dataloader.Loader)
	loader.ClearAll()
}

var HooksDataloader hooksDataloaderFunc = func(ctx context.Context, req core.GraphQLRequest) context.Context {
	ctx = context.WithValue(ctx, "user", dataloader.NewBatchedLoader(LoadUser))
	ctx = context.WithValue(ctx, "userAddresses", dataloader.NewBatchedLoader(LoadUserAddresses))
	ctx = context.WithValue(ctx, "userTransactions", dataloader.NewBatchedLoader(LoadUserTransactions))
	ctx = context.WithValue(ctx, "adminProducts", dataloader.NewBatchedLoader(LoadAdminProducts))
	ctx = context.WithValue(ctx, "products", dataloader.NewBatchedLoader(LoadProducts))
	ctx = context.WithValue(ctx, "product", dataloader.NewBatchedLoader(LoadProduct))
	ctx = context.WithValue(ctx, "productBySlug", dataloader.NewBatchedLoader(LoadProductBySlug))
	ctx = context.WithValue(ctx, "productOptions", dataloader.NewBatchedLoader(LoadProductOptions))
	ctx = context.WithValue(ctx, "productImages", dataloader.NewBatchedLoader(LoadProductImages))
	ctx = context.WithValue(ctx, "productOptionValues", dataloader.NewBatchedLoader(LoadProductOptionValues))
	ctx = context.WithValue(ctx, "productVariants", dataloader.NewBatchedLoader(LoadProductVariants))
	ctx = context.WithValue(ctx, "productVariantBySelectedOptions", dataloader.NewBatchedLoader(LoadProductVariantBySelectedOptions))
	ctx = context.WithValue(ctx, "productVariant", dataloader.NewBatchedLoader(LoadProductVariant))
	ctx = context.WithValue(ctx, "productVariantOptions", dataloader.NewBatchedLoader(LoadProductVariantOptions))
	ctx = context.WithValue(ctx, "productVariantImages", dataloader.NewBatchedLoader(LoadProductVariantImages))
	ctx = context.WithValue(ctx, "transaction", dataloader.NewBatchedLoader(LoadTransaction))
	ctx = context.WithValue(ctx, "transactions", dataloader.NewBatchedLoader(LoadTransactions))
	ctx = context.WithValue(ctx, "transactionAddresses", dataloader.NewBatchedLoader(LoadTransactionAddresses))
	ctx = context.WithValue(ctx, "transactionLineItems", dataloader.NewBatchedLoader(LoadTransactionLineItems))
	ctx = context.WithValue(ctx, "subtotal", dataloader.NewBatchedLoader(LoadSubtotal))
	ctx = context.WithValue(ctx, "taxes", dataloader.NewBatchedLoader(LoadTaxes))
	ctx = context.WithValue(ctx, "shippingEstimations", dataloader.NewBatchedLoader(LoadShippingEstimations))

	return ctx
}
