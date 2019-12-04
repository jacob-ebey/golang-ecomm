package dataloaders

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"

	core "github.com/jacob-ebey/graphql-core"

	"github.com/jacob-ebey/golang-ecomm/db"
)

func LoadUser(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {

		}
		ids[index] = id
	}

	dbResults := []*db.User{}
	if err := database.
		Model(&dbResults).
		WhereIn("id IN (?)", ids).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load user.",
					InternalError: err,
				},
			}
		}
		return results
	}

	resultMap := map[int]*dataloader.Result{}
	for _, product := range dbResults {
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
					Message: "Failed to load user `" + key.String() + "`.",
				},
			}
		}

		results[index] = result
	}

	return results
}

func LoadUserAddresses(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {

		}
		ids[index] = id
	}

	dbResults := []*db.Address{}
	if err := database.
		Model(&dbResults).
		WhereIn("address.user_id IN (?)", ids).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load addresses for user.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int][]*db.Address{}
	for _, address := range dbResults {
		if resultMap[address.UserID] == nil {
			resultMap[address.UserID] = []*db.Address{}
		}

		resultMap[address.UserID] = append(resultMap[address.UserID], address)
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load addresses for user `" + key.String() + "`.",
				},
			}
		}

		results[index] = &dataloader.Result{
			Data: result,
		}
	}

	return results
}

func LoadUserTransactions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {

		}
		ids[index] = id
	}

	dbResults := []*db.Transaction{}
	if err := database.
		Model(&dbResults).
		OrderExpr("id DESC").
		WhereIn("transaction.user_id IN (?)", ids).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load transactions for user.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int][]*db.Transaction{}
	for _, transaction := range dbResults {
		if resultMap[transaction.UserID] == nil {
			resultMap[transaction.UserID] = []*db.Transaction{}
		}

		resultMap[transaction.UserID] = append(resultMap[transaction.UserID], transaction)
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load transactions for user `" + key.String() + "`.",
				},
			}
		}

		results[index] = &dataloader.Result{
			Data: result,
		}
	}

	return results
}
