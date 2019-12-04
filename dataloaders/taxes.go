package dataloaders

import (
	"context"
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader"
	"github.com/jacob-ebey/golang-ecomm/apis"
	core "github.com/jacob-ebey/graphql-core"

	"github.com/jacob-ebey/golang-ecomm/db"
)

type Taxes struct {
	TotalRate float64
	Rates     []Rate
}

type Rate struct {
	Rate float64
	Name string
	Type string
}

func LoadTaxes(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	avatax := ctx.Value("avatax").(*apis.Avatax)

	results := make([]*dataloader.Result, len(keys))

	for index, key := range keys {
		address, ok := key.Raw().(db.Address)
		if !ok {
			results[index] = &dataloader.Result{
				Error: fmt.Errorf("Improper key type provided."),
			}

			continue
		}

		taxes, err := loadTaxes(avatax, apis.AvataxAddress{
			Line1:      address.Line1,
			Line2:      address.Line2,
			Line3:      address.Line3,
			City:       address.City,
			Region:     address.Region,
			Country:    address.Country,
			PostalCode: address.PostalCode,
		})

		if err != nil {
			results[index] = &dataloader.Result{
				Error: err,
			}

			continue
		}

		results[index] = &dataloader.Result{
			Data: taxes,
		}
	}

	return results
}

func loadTaxes(avatax *apis.Avatax, address apis.AvataxAddress) (*Taxes, error) {
	rates, err := avatax.TaxRatesByAddress(address)
	if err != nil {
		switch err.(type) {
		case *apis.AvataxError:
			if strings.Contains(err.Error(), "CreateTransaction()") {
				return nil, &core.WrappedError{
					Message:       "Could not get taxes for address.",
					InternalError: err,
				}
			}

			return nil, err
		default:
			return nil, &core.WrappedError{
				Message:       "Could not get taxes for address.",
				InternalError: err,
			}
		}
	}

	subRates := make([]Rate, len(rates.Rates))
	for index, rate := range rates.Rates {
		subRates[index] = Rate{
			Rate: rate.Rate,
			Name: rate.Name,
			Type: rate.Type,
		}
	}

	return &Taxes{
		TotalRate: rates.TotalRate,
		Rates:     subRates,
	}, nil
}
