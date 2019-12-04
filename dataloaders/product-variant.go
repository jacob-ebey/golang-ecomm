package dataloaders

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"

	core "github.com/jacob-ebey/graphql-core"

	"github.com/jacob-ebey/golang-ecomm/db"
)

func LoadProductVariant(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
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
		AllWithDeleted().
		WhereIn("product_variant.id IN (?)", ids).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load product variant.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int]*dataloader.Result{}
	for _, variant := range dbResults {
		resultMap[variant.ID] = &dataloader.Result{
			Data: variant,
		}
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok || result == nil {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load product variant `" + key.String() + "`.",
				},
			}
			continue
		}

		results[index] = result
	}

	return results
}

func LoadProductVariantOptions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {
			continue
		}
		ids[index] = id
	}

	dbResults := []*db.ProductVariantOption{}
	if err := database.
		Model(&dbResults).
		Column("product_variant_option.product_variant_id").
		Column("product_variant_option.product_option_value_id").
		WhereIn("product_variant_option.product_variant_id IN (?)", ids).
		Relation("ProductOptionValue").
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load product variant options.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int][]*db.ProductOptionValue{}
	for _, option := range dbResults {
		if resultMap[option.ProductVariantID] == nil {
			resultMap[option.ProductVariantID] = []*db.ProductOptionValue{}
		}

		resultMap[option.ProductVariantID] = append(resultMap[option.ProductVariantID], option.ProductOptionValue)
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load product variant options for `" + key.String() + "`.",
				},
			}
			continue
		}

		results[index] = &dataloader.Result{
			Data: result,
		}
	}

	return results
}

func LoadProductVariantBySelectedOptions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)
	productVariantLoader := ctx.Value("productVariant").(*dataloader.Loader)

	toLookup := make([]SelectedOptionsKey, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(SelectedOptionsKey)
		if !ok {
			continue
		}
		toLookup[index] = id
	}

	results := make([]*dataloader.Result, len(toLookup))

	for index, lookup := range toLookup {
		variant, err := productVariantBySelectedOptions(database, lookup.ProductID, lookup.SelectedOptions)

		if err != nil {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load product with provided selected options",
					InternalError: err,
				},
			}

			continue
		}

		if variant != nil {
			productVariantLoader.Prime(ctx, IntKey(variant.ID), variant)
		}

		results[index] = &dataloader.Result{
			Data: variant,
		}
	}

	return results
}

type SelectedOptionsKey struct {
	ProductID       int
	SelectedOptions []int
}

func (key SelectedOptionsKey) String() string {
	result := strconv.Itoa(key.ProductID) + "|"

	for index, option := range key.SelectedOptions {
		if index > 0 {
			result += ","
		}
		result += strconv.Itoa(option)
	}

	return result
}

func (key SelectedOptionsKey) Raw() interface{} {
	return key
}

func productVariantBySelectedOptions(database *pg.DB, productID int, selectedOptions []int) (*db.ProductVariant, error) {
	results := []db.ProductVariant{}

	query := `
SELECT a.*
FROM product_variants a
WHERE deleted_at IS NULL AND NOT EXISTS (
SELECT *
FROM product_variant_options b
WHERE b.product_id = ? AND a.id = b.product_variant_id
)
LIMIT 2
	`

	if len(selectedOptions) > 0 {
		where := ""
		for i, option := range selectedOptions {
			if i > 0 {
				where = where + ","
			}
			where = where + fmt.Sprintf("%d", option)
		}

		query = fmt.Sprintf(`
SELECT a.*
FROM product_variants a
INNER JOIN product_variant_options b
ON a.id = b.product_variant_id AND b.product_id = ?
WHERE a.deleted_at IS NULL AND b.product_option_value_id IN (%v)
GROUP BY a.id
HAVING COUNT(DISTINCT b.product_option_value_id) = ?
		`, where)
	}

	found, err := database.
		Query(&results, query, productID, len(selectedOptions))

	if err != nil {
		return nil, &core.WrappedError{
			Message:       "Failed to find product variant by selected options.",
			InternalError: err,
		}
	}

	if found.RowsReturned() != 1 {
		return nil, nil
	}

	return &results[0], nil
}

func LoadProductVariantImages(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {
			continue
		}
		ids[index] = id
	}

	dbResults := []*db.ProductVariantImage{}
	if err := database.
		Model(&dbResults).
		Column("product_variant_image.product_variant_id").
		Column("product_variant_image.image_id").
		WhereIn("product_variant_image.product_variant_id IN (?)", ids).
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
		if resultMap[option.ProductVariantID] == nil {
			resultMap[option.ProductVariantID] = []*db.Image{}
		}

		resultMap[option.ProductVariantID] = append(resultMap[option.ProductVariantID], option.Image)
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
