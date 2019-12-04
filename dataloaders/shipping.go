package dataloaders

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/jacob-ebey/go-shippo/client"
	"github.com/jacob-ebey/go-shippo/models"
	"github.com/graph-gophers/dataloader"
	"github.com/jacob-ebey/golang-ecomm/db"
	core "github.com/jacob-ebey/graphql-core"
)

type ShippingEstimation struct {
	ID            string
	Price         int
	Service       string
	Carrier       string
	DurationTerms string
}

type ShippingParcel struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type ShippingEstimationKey struct {
	Address  db.Address
	Variants CartKey
}

func (key ShippingEstimationKey) String() string {
	return key.Address.String() + "***" + key.Variants.String()
}

func (key ShippingEstimationKey) Raw() interface{} {
	return key
}

func LoadShippingEstimations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	results := make([]*dataloader.Result, len(keys))

	for index, key := range keys {
		toEstimate, ok := key.Raw().(ShippingEstimationKey)
		if !ok {
			results[index] = &dataloader.Result{
				Error: fmt.Errorf("Improper key type provided."),
			}

			continue
		}

		estimations, err := loadShippingEstimations(ctx, toEstimate.Address, db.Address{
			Name:       "Space Needle",
			Line1:      "400 Broad St",
			City:       "Seattle",
			Region:     "WA",
			Country:    "USA",
			PostalCode: "98109",
		}, toEstimate.Variants)

		if err != nil {
			results[index] = &dataloader.Result{
				Error: err,
			}

			continue
		}

		results[index] = &dataloader.Result{
			Data: estimations,
		}
	}

	return results
}

func createAddress(shippoClient *client.Client, address db.Address) (*models.Address, error) {
	return shippoClient.CreateAddress(&models.AddressInput{
		Name:     address.Name,
		Street1:  address.Line1,
		Street2:  address.Line2,
		Street3:  address.Line3,
		City:     address.City,
		State:    address.Region,
		Zip:      address.PostalCode,
		Country:  address.Country,
		Validate: true,
	})
}

func loadShippingEstimations(
	ctx context.Context,
	toAddr db.Address,
	fromAddr db.Address,
	toEstimate []CartVariant) ([]*ShippingEstimation, error) {
	shippoClient := ctx.Value("shippo").(*client.Client)
	productVariant := ctx.Value("productVariant").(*dataloader.Loader)

	addressFrom, err := createAddress(shippoClient, fromAddr)
	if err != nil {
		return nil, err // TODO: Lookover the type of messages this error has
	}

	addressTo, err := createAddress(shippoClient, toAddr)
	if err != nil {
		return nil, err
	}

	ids := make([]dataloader.Key, len(toEstimate))
	for index, toEst := range toEstimate {
		ids[index] = IntKey(toEst.VariantID)
	}

	variants, errs := productVariant.LoadMany(ctx, ids)()
	if errs != nil {
		return nil, &core.WrappedError{
			Message:       "Could not get variants for estimation.",
			InternalError: HandleErrors(errs),
		}
	}

	parcels := map[int]*models.Parcel{}
	for _, tempVariant := range variants {
		variant := tempVariant.(*db.ProductVariant)

		parcelInput := &models.ParcelInput{
			Length:       fmt.Sprintf("%.2f", variant.Length),
			Width:        fmt.Sprintf("%.2f", variant.Width),
			Height:       fmt.Sprintf("%.2f", variant.Height),
			DistanceUnit: models.DistanceUnitInch,
			Weight:       fmt.Sprintf("%.2f", variant.Weight),
			MassUnit:     models.MassUnitOunce,
		}
		parcel, err := shippoClient.CreateParcel(parcelInput)
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create parcel.",
				InternalError: err,
			}
		}

		parcels[variant.ID] = parcel
	}

	parcelsToEstimate := []string{}
	for _, item := range toEstimate {
		for i := 0; i < item.Quantity; i++ {
			parcelsToEstimate = append(parcelsToEstimate, parcels[item.VariantID].ObjectID)
		}
	}

	shipment, err := shippoClient.CreateShipment(&models.ShipmentInput{
		AddressFrom: addressFrom.ObjectID,
		AddressTo:   addressTo.ObjectID,
		Parcels:     parcelsToEstimate,
		Async:       false,
	})
	if err != nil {
		return nil, &core.WrappedError{
			Message:       "Could not create shipping estimation.",
			InternalError: err,
		}
	}

	estimations := make([]*ShippingEstimation, len(shipment.Rates))
	for index, rate := range shipment.Rates {
		amount, err := strconv.ParseFloat(rate.Amount, 64)
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not convert estimation price.",
				InternalError: err,
			}
		}
		price := int(math.Round(amount * 100))

		estimations[index] = &ShippingEstimation{
			ID:            rate.ObjectID,
			Price:         price,
			Carrier:       rate.Provider,
			Service:       rate.ServiceLevel.Name,
			DurationTerms: rate.DurationTerms,
		}
	}

	return estimations, nil
}
