package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/braintree-go/braintree-go"
	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/jacob-ebey/go-shippo/client"
	"github.com/jacob-ebey/golang-ecomm/auth"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	"github.com/jacob-ebey/golang-ecomm/db"
	"github.com/jacob-ebey/golang-ecomm/email"
	"github.com/jacob-ebey/golang-ecomm/services"
	core "github.com/jacob-ebey/graphql-core"
)

var BraintreeClientTokenField = &graphql.Field{
	Type:        graphql.NewNonNull(graphql.String),
	Description: "Get a braintree client token.",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		braintreeClient := params.Context.Value("braintree").(*braintree.Braintree)

		token, err := braintreeClient.ClientToken().Generate(params.Context)

		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create braintree client token.",
				InternalError: err,
			}
		}

		return token, nil
	},
}

var SubmitBraintreeTransactionField = &graphql.Field{
	Type: ReceiptType,
	Args: graphql.FieldConfigArgument{
		"braintreeNonce": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"billingAddressId": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"billingAddress": &graphql.ArgumentConfig{
			Type: AddressInputSchema,
		},
		"saveBillingAddress": &graphql.ArgumentConfig{
			Type: graphql.Boolean,
		},
		"shippingAddressId": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"shippingAddress": &graphql.ArgumentConfig{
			Type: AddressInputSchema,
		},
		"saveShippingAddress": &graphql.ArgumentConfig{
			Type: graphql.Boolean,
		},
		"total": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"shippingRateId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"variants": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(CartInputSchema))),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		database := params.Context.Value("database").(*pg.DB)
		shippoClient := params.Context.Value("shippo").(*client.Client)
		subtotalLoader := params.Context.Value("subtotal").(*dataloader.Loader)
		taxesLoader := params.Context.Value("taxes").(*dataloader.Loader)
		productLoader := params.Context.Value("product").(*dataloader.Loader)
		productVariantLoader := params.Context.Value("productVariant").(*dataloader.Loader)
		braintreeClient := params.Context.Value("braintree").(*braintree.Braintree)
		emailClient := params.Context.Value("email").(email.Client)
		baseUrl := params.Context.Value("baseUrl").(string)

		claims := params.Context.Value("claims").(*auth.Claims)
		userID := 0
		if claims != nil {
			userID = claims.ID
		}

		braintreeNonce := params.Args["braintreeNonce"].(string)

		billingAddressID, _ := params.Args["billingAddressId"].(int)
		var billingAddress *db.Address = nil
		if billingAddressTemp, ok := params.Args["billingAddress"]; ok {
			billingAddress = &db.Address{}
			if err := ConvertObject(billingAddressTemp, billingAddress); err != nil {
				return nil, &core.WrappedError{
					Message:       "Could not decode billingAddress.",
					InternalError: err,
				}
			}
		}
		saveBillingAddress, _ := params.Args["saveBillingAddress"].(bool)

		shippingAddressID, _ := params.Args["shippingAddressId"].(int)
		var shippingAddress *db.Address = nil
		if shippingAddressTemp, ok := params.Args["shippingAddress"]; ok {
			shippingAddress = &db.Address{}
			if err := ConvertObject(shippingAddressTemp, shippingAddress); err != nil {
				return nil, &core.WrappedError{
					Message:       "Could not decode shippingAddress.",
					InternalError: err,
				}
			}
		}
		saveShippingAddress, _ := params.Args["saveShippingAddress"].(bool)

		total := params.Args["total"].(int)

		shippingRateID := params.Args["shippingRateId"].(string)

		cart := dataloaders.CartKey{}
		if err := ConvertObject(params.Args["variants"], &cart); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not convert variants argument.",
				InternalError: err,
			}
		}

		shippingAddress, err := getAddress(params.Context, shippingAddress, shippingAddressID, saveShippingAddress)
		if err != nil {
			return nil, err
		}
		billingAddress, err = getAddress(params.Context, billingAddress, billingAddressID, saveBillingAddress)
		if err != nil {
			return nil, err
		}

		subtotalCalculatedTemp, err := subtotalLoader.Load(params.Context, cart)()
		if err != nil {
			return nil, err
		}
		subtotalCalculated := subtotalCalculatedTemp.(int)

		taxesCalculatedTemp, err := taxesLoader.Load(params.Context, *billingAddress)()
		if err != nil || taxesCalculatedTemp == nil {
			return nil, &core.WrappedError{
				Message:       "Could not load tax information",
				InternalError: err,
			}
		}
		taxRates := taxesCalculatedTemp.(*dataloaders.Taxes)
		taxesCalculated := int(math.Round(float64(subtotalCalculated) * taxRates.TotalRate))

		rate, err := shippoClient.RetrieveRate(shippingRateID)
		if err != nil || rate == nil {
			return nil, &core.WrappedError{
				Message:       "Could not retrieve shipping rate.",
				InternalError: err,
			}
		}

		amount, err := strconv.ParseFloat(rate.Amount, 64)
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not convert shipping price.",
				InternalError: err,
			}
		}
		shippingCalculated := int(math.Round(amount * 100))

		totalCalculated := subtotalCalculated + taxesCalculated + shippingCalculated

		if totalCalculated != total {
			return nil, &core.WrappedError{
				Message: "The provided total does not match the calculated ones.",
			}
		}

		variantIDs := make(dataloader.Keys, len(cart))
		for index, lineItem := range cart {
			variantIDs[index] = dataloaders.IntKey(lineItem.VariantID)
		}
		variantsTemp, errs := productVariantLoader.LoadMany(params.Context, variantIDs)()
		if errs != nil && len(errs) > 0 {
			return nil, &core.WrappedError{
				Message:       "Could not loaad variants.",
				InternalError: dataloaders.HandleErrors(errs),
			}
		}
		variantMap := map[int]*db.ProductVariant{}
		for _, variantTemp := range variantsTemp {
			variant := variantTemp.(*db.ProductVariant)
			variantMap[variant.ID] = variant
		}

		result := db.Transaction{
			Subtotal:     subtotalCalculated,
			Taxes:        taxesCalculated,
			Shipping:     shippingCalculated,
			Total:        totalCalculated,
			ShippoRateID: rate.ObjectID,
			UserID:       userID,
		}
		if err := database.Insert(&result); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create transaction.",
				InternalError: err,
			}
		}

		status := db.TransactionStatus{
			CreatedAt:     time.Now(),
			TransactionID: result.ID,
			Status:        "RECEIVED",
		}
		if err := database.Insert(&status); err != nil {
			database.Delete(&result)
			database.ForceDelete(&result)
			return nil, &core.WrappedError{
				Message:       "Could not create transaction.",
				InternalError: err,
			}
		}

		addressInfo := db.TransactionAddressInfo{
			TransactionID:     result.ID,
			BillingAddressID:  billingAddress.ID,
			ShippingAddressID: shippingAddress.ID,
		}
		if err := database.Insert(&addressInfo); err != nil {
			database.Delete(&status)
			database.ForceDelete(&status)
			database.Delete(&result)
			database.ForceDelete(&result)

			return nil, &core.WrappedError{
				Message:       "Could not create transaction.",
				InternalError: err,
			}
		}
		addressInfo.BillingAddress = billingAddress
		addressInfo.ShippingAddress = shippingAddress

		createdLineItems := []*db.TransactionLineItem{}
		for _, lineItem := range cart {
			toCreate := db.TransactionLineItem{
				TransactionID:    result.ID,
				ProductVariantID: lineItem.VariantID,
				Quantity:         lineItem.Quantity,
				Price:            variantMap[lineItem.VariantID].Price,
			}
			if err = database.Insert(&toCreate); err != nil {
				break
			}

			createdLineItems = append(createdLineItems, &toCreate)
		}

		if err != nil {
			for _, lineItem := range createdLineItems {
				database.Delete(&lineItem)
				database.ForceDelete(&lineItem)
			}

			database.Delete(&addressInfo)
			database.ForceDelete(&addressInfo)

			database.Delete(&result)
			database.ForceDelete(&result)

			return nil, &core.WrappedError{
				Message:       "Could not create transaction.",
				InternalError: err,
			}
		}

		braintreeLineItems := make([]*braintree.TransactionLineItemRequest, len(createdLineItems))
		for index, lineItem := range createdLineItems {
			name := variantMap[lineItem.ProductVariantID].Name

			if name == "" {
				product, err := productLoader.Load(params.Context, dataloaders.IntKey(variantMap[lineItem.ProductVariantID].ProductID))()
				if err != nil || product == nil {
					return nil, &core.WrappedError{
						Message:       "Could not get product variant name.",
						InternalError: err,
					}
				}

				name = product.(*db.Product).Name
			}

			braintreeLineItems[index] = &braintree.TransactionLineItemRequest{
				Kind:        braintree.TransactionLineItemKindDebit,
				Quantity:    braintree.NewDecimal(int64(lineItem.Quantity), 0),
				Name:        name,
				UnitAmount:  braintree.NewDecimal(int64(lineItem.Price), 2),
				TotalAmount: braintree.NewDecimal(int64(lineItem.Price*lineItem.Quantity), 2),
			}
		}

		extendedAddress := shippingAddress.Line2 + "," + shippingAddress.Line3
		if extendedAddress == "," {
			extendedAddress = ""
		}

		braintreeAddress := braintree.Address{
			StreetAddress:   shippingAddress.Line1,
			ExtendedAddress: extendedAddress,
			Locality:        shippingAddress.City,
			Region:          shippingAddress.Region,
			PostalCode:      shippingAddress.PostalCode,
		}

		if _, err := strconv.Atoi(shippingAddress.Country); err == nil {
			braintreeAddress.CountryCodeNumeric = shippingAddress.Country
		} else if len(shippingAddress.Country) == 2 {
			braintreeAddress.CountryCodeAlpha2 = shippingAddress.Country
		} else if len(shippingAddress.Country) == 3 {
			braintreeAddress.CountryCodeAlpha3 = shippingAddress.Country
		} else {
			braintreeAddress.CountryName = shippingAddress.Country
		}

		braintreeTransaction, err := braintreeClient.Transaction().Create(params.Context, &braintree.TransactionRequest{
			Type:               "sale",
			PaymentMethodNonce: braintreeNonce,
			Options: &braintree.TransactionOptions{
				SubmitForSettlement: true,
			},
			OrderId:         strconv.Itoa(result.ID),
			Amount:          braintree.NewDecimal(int64(result.Total), 2),
			TaxAmount:       braintree.NewDecimal(int64(result.Taxes), 2),
			LineItems:       braintreeLineItems,
			ShippingAddress: &braintreeAddress,
		})

		if err != nil {
			j, _ := json.MarshalIndent(err, "", "\t")
			fmt.Println(string(j))
			database.Delete(&result)
			database.ForceDelete(&result)

			return nil, err
		}

		result.BraintreeID = braintreeTransaction.Id
		if err := database.Update(&result); err != nil {
			fmt.Println("Failed to update transaction with Braintree ID.")
			fmt.Println(err)
		}

		toSend, err := email.NewPurchaseEmail(baseUrl)
		if err != nil {
			fmt.Println("Failed create purchase email.")
			fmt.Println(err)
		} else if claims != nil {
			err = emailClient.SendMail(claims.Email, "Thanks for your purchase.", toSend)
			if err != nil {
				fmt.Println("Failed to send purchase email.")
				fmt.Println(err)
			}
		}

		return &result, nil
	},
}

func getAddress(ctx context.Context, address *db.Address, addressID int, saveAddress bool) (*db.Address, error) {
	database := ctx.Value("database").(*pg.DB)
	addressValidator := ctx.Value("addressValidator").(services.AddressValidator)

	if address != nil {
		claims := ctx.Value("claims").(*auth.Claims)
		if saveAddress && claims != nil {
			address.UserID = claims.ID

			_, err := addressValidator.ValidateAddress(ctx, *address)
			if err != nil {
				return nil, err
			}

			database.Insert(address)
		}

		return address, nil
	}

	newAddress := db.Address{}
	if err := database.Model(&newAddress).Where("id = ?", addressID).Select(); err != nil {
		return nil, &core.WrappedError{
			Message:       "Could not find address.",
			InternalError: err,
		}
	}

	return &newAddress, nil
}
