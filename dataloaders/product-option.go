package dataloaders

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"

	core "github.com/jacob-ebey/graphql-core"

	"github.com/jacob-ebey/golang-ecomm/db"
)

func LoadProductOptionValues(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {

		}
		ids[index] = id
	}

	dbResults := []*db.ProductOptionValue{}
	if err := database.
		Model(&dbResults).
		Column("product_option_value.*").
		WhereIn("product_option_value.product_option_id IN (?)", ids).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load product option values.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int][]*db.ProductOptionValue{}
	for _, value := range dbResults {
		if resultMap[value.ProductOptionID] == nil {
			resultMap[value.ProductOptionID] = []*db.ProductOptionValue{}
		}

		resultMap[value.ProductOptionID] = append(resultMap[value.ProductOptionID], value)
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load product option values for `" + key.String() + "`.",
				},
			}
		}

		results[index] = &dataloader.Result{
			Data: result,
		}
	}

	return results
}
