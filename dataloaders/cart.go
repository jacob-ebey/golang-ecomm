package dataloaders

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"
	core "github.com/jacob-ebey/graphql-core"

	"github.com/jacob-ebey/golang-ecomm/db"
)

type CartVariant struct {
	VariantID int
	Quantity  int
}

type CartKey []CartVariant

func (key CartKey) String() string {
	sort.SliceStable(key, func(i, j int) bool {
		return key[i].VariantID < key[j].VariantID
	})

	result := ""
	for index, variant := range key {
		if index > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d|%d", variant.VariantID, variant.Quantity)
	}

	return result
}

func (key CartKey) Raw() interface{} {
	return key
}

func LoadSubtotal(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	database := ctx.Value("database").(*pg.DB)

	results := make([]*dataloader.Result, len(keys))

	for index, key := range keys {
		key, ok := key.Raw().(CartKey)
		if !ok {
			results[index] = &dataloader.Result{
				Error: fmt.Errorf("Improper key type provided."),
			}

			continue
		}

		res, err := subtotal(database, key)
		if err != nil {
			results[index] = &dataloader.Result{
				Error: err,
			}

			continue
		}

		results[index] = &dataloader.Result{
			Data: res,
		}
	}

	return results
}

var InvalidQuantityError = fmt.Errorf("Quanity for each variant must be greater than 0.")
var MissmatchedVariantsError = fmt.Errorf("Failed to calculate subtotal. One or more variants is not avaliable for purchase.")

func subtotal(database *pg.DB, variants CartKey) (int, error) {
	if len(variants) == 0 {
		return 0, nil
	}

	quanties := map[int]int{}
	where := make([]int, len(variants))
	for index, variant := range variants {
		if variant.Quantity < 1 {
			return 0, InvalidQuantityError
		}

		quanties[variant.VariantID] = variant.Quantity
		where[index] = variant.VariantID
	}

	found := []db.ProductVariant{}
	if err := database.
		Model(&found).
		Column("product_variant.id").
		Column("product_variant.price").
		WhereIn("product_variant.id IN (?)", where).
		Select(); err != nil {
		return 0, &core.WrappedError{
			Message:       "Failed to calculate subtotal",
			InternalError: err,
		}
	}

	if len(found) != len(variants) {
		return 0, MissmatchedVariantsError
	}

	var subtotal int = 0
	for _, variant := range found {
		if quantity, ok := quanties[variant.ID]; ok {
			subtotal += variant.Price * int(quantity)
			continue
		}

		return 0, MissmatchedVariantsError
	}

	return subtotal, nil
}
