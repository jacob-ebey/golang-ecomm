package services

import (
	"context"

	"github.com/jacob-ebey/go-shippo/client"
	"github.com/jacob-ebey/go-shippo/models"
	"github.com/jacob-ebey/golang-ecomm/db"
	core "github.com/jacob-ebey/graphql-core"
)

type AddressValidator interface {
	ValidateAddress(ctx context.Context, address db.Address) (bool, error)
}

type validateAddressFunc func(ctx context.Context, address db.Address) (bool, error)

func (validate validateAddressFunc) ValidateAddress(ctx context.Context, address db.Address) (bool, error) {
	return validate(ctx, address)
}

func (hook validateAddressFunc) PreExecute(ctx context.Context, req core.GraphQLRequest) context.Context {
	return context.WithValue(ctx, "addressValidator", hook)
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

var ValidateAddressWithShippo validateAddressFunc = func(ctx context.Context, address db.Address) (bool, error) {
	shippoClient := ctx.Value("shippo").(*client.Client)

	addr, err := createAddress(shippoClient, address)
	if err != nil || addr == nil {
		return false, &core.WrappedError{
			Message:       "Address is not valid.",
			InternalError: err,
		}
	}

	return true, nil
}
