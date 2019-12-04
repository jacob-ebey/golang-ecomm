package dataloaders

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"

	core "github.com/jacob-ebey/graphql-core"

	"github.com/jacob-ebey/golang-ecomm/db"
)

func LoadTransactions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)
	transactionLoader := ctx.Value("transaction").(*dataloader.Loader)

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
		results := []*db.Transaction{}

		if err := database.
			Model(&results).
			OrderExpr("id DESC").
			Offset(page.Skip).
			Limit(page.Limit).
			Select(); err != nil {
			fmt.Println(err)
			pages[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load transaction page.",
					InternalError: err,
				},
			}
			continue
		}

		for _, result := range results {
			transactionLoader.Prime(ctx, IntKey(result.ID), result)
		}

		pages[index] = &dataloader.Result{
			Data: results,
		}
	}

	return pages
}

func LoadTransaction(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
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
		WhereIn("transaction.id IN (?)", ids).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load transaction.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int]*dataloader.Result{}
	for _, transaction := range dbResults {
		resultMap[transaction.ID] = &dataloader.Result{
			Data: transaction,
		}
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load transaction `" + key.String() + "`.",
				},
			}
		}

		results[index] = result
	}

	return results
}

func LoadTransactionLineItems(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {

		}
		ids[index] = id
	}

	dbResults := []*db.TransactionLineItem{}
	if err := database.
		Model(&dbResults).
		WhereIn("transaction_line_item.transaction_id IN (?)", ids).
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load transaction line items.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int][]*db.TransactionLineItem{}
	for _, lineItem := range dbResults {
		if resultMap[lineItem.TransactionID] == nil {
			resultMap[lineItem.TransactionID] = []*db.TransactionLineItem{}
		}

		resultMap[lineItem.TransactionID] = append(resultMap[lineItem.TransactionID], lineItem)
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load line items for `" + key.String() + "`.",
				},
			}
		}

		results[index] = &dataloader.Result{
			Data: result,
		}
	}

	return results
}

func LoadTransactionAddresses(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	ids := make([]int, len(keys))
	for index, key := range keys {
		id, ok := key.Raw().(int)
		if !ok {
			// TODO: Handle error
		}
		ids[index] = id
	}

	dbResults := []*db.TransactionAddressInfo{}
	if err := database.
		Model(&dbResults).
		Column("transaction_address_info.transaction_id").
		Column("transaction_address_info.billing_address_id").
		Column("transaction_address_info.shipping_address_id").
		WhereIn("transaction_address_info.transaction_id IN (?)", ids).
		Relation("BillingAddress").
		Relation("ShippingAddress").
		Select(); err != nil {
		results := make([]*dataloader.Result, len(keys))
		for index, _ := range keys {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message:       "Failed to load transaction addresses.",
					InternalError: err,
				},
			}
		}

		return results
	}

	resultMap := map[int]*dataloader.Result{}
	for _, addressInfo := range dbResults {
		resultMap[addressInfo.TransactionID] = &dataloader.Result{
			Data: addressInfo,
		}
	}

	results := make([]*dataloader.Result, len(keys))
	for index, key := range keys {
		result, ok := resultMap[key.Raw().(int)]

		if !ok {
			results[index] = &dataloader.Result{
				Error: &core.WrappedError{
					Message: "Failed to load transaction addresses for `" + key.String() + "`.",
				},
			}
		}

		results[index] = result
	}

	return results
}
