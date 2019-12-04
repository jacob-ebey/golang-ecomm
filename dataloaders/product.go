package dataloaders

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"

	core "github.com/jacob-ebey/graphql-core"

	"github.com/jacob-ebey/golang-ecomm/db"
)

func LoadProducts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)
	productLoader := ctx.Value("product").(*dataloader.Loader)
	productBySlug := ctx.Value("productBySlug").(*dataloader.Loader)

	pagination := make([]PaginationKey, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(PaginationKey)
		if !ok {
			continue
		}
		pagination[index] = id
	}

	pages := make([]*dataloader.Result, len(pagination))

	for index, page := range pagination {
		results := []*db.Product{}

		if err := database.
			Model(&results).
			OrderExpr("id DESC").
			Offset(page.Skip).
			Limit(page.Limit).
			Where("product.published IS TRUE").
			Select(); err != nil {
			pages[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load product page.",
					InternalError: err,
				},
			}
			continue
		}

		for _, result := range results {
			productLoader.Prime(ctx, IntKey(result.ID), result)
			productBySlug.Prime(ctx, dataloader.StringKey(result.Slug), result)
		}

		pages[index] = &dataloader.Result{
			Data: results,
		}
	}

	return pages
}

func LoadAdminProducts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)
	productLoader := ctx.Value("product").(*dataloader.Loader)
	productBySlug := ctx.Value("productBySlug").(*dataloader.Loader)

	pagination := make([]PaginationKey, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(PaginationKey)
		if !ok {
			continue
		}
		pagination[index] = id
	}

	pages := make([]*dataloader.Result, len(pagination))

	for index, page := range pagination {
		results := []*db.Product{}

		if err := database.
			Model(&results).
			OrderExpr("id DESC").
			Offset(page.Skip).
			Limit(page.Limit).
			Select(); err != nil {
			pages[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load admin product page.",
					InternalError: err,
				},
			}
			continue
		}

		for _, result := range results {
			productLoader.Prime(ctx, IntKey(result.ID), result)
			productBySlug.Prime(ctx, dataloader.StringKey(result.Slug), result)
		}

		pages[index] = &dataloader.Result{
			Data: results,
		}
	}

	return pages
}

func LoadProduct(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)
	productBySlug := ctx.Value("productBySlug").(*dataloader.Loader)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {
			continue
		}
		ids[index] = id
	}

	dbResults := []*db.Product{}
	if err := database.
		Model(&dbResults).
		WhereIn("product.id IN (?)", ids).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load product.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int]*dataloader.Result{}
	for _, product := range dbResults {
		productBySlug.Prime(ctx, dataloader.StringKey(product.Slug), product)

		resultMap[product.ID] = &dataloader.Result{
			Data: product,
		}
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load product `" + key.String() + "`.",
				},
			}
		}

		results[index] = result
	}

	return results
}

func LoadProductBySlug(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)
	productLoader := ctx.Value("product").(*dataloader.Loader)

	slugs := make([]string, len(keys))
	for index, key := range keys {
		slugs[index] = key.String()
	}

	dbResults := []*db.Product{}
	if err := database.
		Model(&dbResults).
		WhereIn("product.slug IN (?)", slugs).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load product.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[string]*dataloader.Result{}
	for _, product := range dbResults {
		productLoader.Prime(ctx, IntKey(product.ID), product)

		resultMap[product.Slug] = &dataloader.Result{
			Data: product,
		}
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.String()]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load product `" + key.String() + "`.",
				},
			}
		}

		results[index] = result
	}

	return results
}

func LoadProductOptions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {

		}
		ids[index] = id
	}

	dbResults := []*db.ProductOption{}
	if err := database.
		Model(&dbResults).
		Column("product_option.*").
		WhereIn("product_option.product_id IN (?)", ids).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load product options.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int][]*db.ProductOption{}
	for _, option := range dbResults {
		if resultMap[option.ProductID] == nil {
			resultMap[option.ProductID] = []*db.ProductOption{}
		}

		resultMap[option.ProductID] = append(resultMap[option.ProductID], option)
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load product options for `" + key.String() + "`.",
				},
			}
		}

		results[index] = &dataloader.Result{
			Data: result,
		}
	}

	return results
}

func LoadProductVariants(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {

		}
		ids[index] = id
	}

	dbResults := []*db.ProductVariant{}
	if err := database.
		Model(&dbResults).
		Column("product_variant.*").
		WhereIn("product_variant.product_id IN (?)", ids).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load product variants.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int][]*db.ProductVariant{}
	for _, variant := range dbResults {
		if resultMap[variant.ProductID] == nil {
			resultMap[variant.ProductID] = []*db.ProductVariant{}
		}

		resultMap[variant.ProductID] = append(resultMap[variant.ProductID], variant)
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load product variants for `" + key.String() + "`.",
				},
			}
		}

		results[index] = &dataloader.Result{
			Data: result,
		}
	}

	return results
}

func LoadProductImages(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {
			continue
		}
		ids[index] = id
	}

	dbResults := []*db.ProductImage{}
	if err := database.
		Model(&dbResults).
		Column("product_image.product_id").
		Column("product_image.image_id").
		WhereIn("product_image.product_id IN (?)", ids).
		Relation("Image").
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load product images.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int][]*db.Image{}
	for _, option := range dbResults {
		if resultMap[option.ProductID] == nil {
			resultMap[option.ProductID] = []*db.Image{}
		}

		resultMap[option.ProductID] = append(resultMap[option.ProductID], option.Image)
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, _ := resultMap[key.Raw().(int)]

		results[index] = &dataloader.Result{
			Data: result,
		}
	}

	return results
}
