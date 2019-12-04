package schema

import (
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/jacob-ebey/go-shippo/client"
	"github.com/jacob-ebey/go-shippo/models"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	"github.com/jacob-ebey/golang-ecomm/db"
	"github.com/jacob-ebey/golang-ecomm/email"
	core "github.com/jacob-ebey/graphql-core"
)

var PurchaseShippoLabelField = &graphql.Field{
	Type:        ShippingLabelType,
	Description: "Purchase a shippo label for a transaction.",
	Args: graphql.FieldConfigArgument{
		"transactionId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"shippoRateId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		database := params.Context.Value("database").(*pg.DB)
		transactionLoader := params.Context.Value("transaction").(*dataloader.Loader)
		userLoader := params.Context.Value("user").(*dataloader.Loader)
		shippoClient := params.Context.Value("shippo").(*client.Client)
		emailClient := params.Context.Value("email").(email.Client)

		transactionId := params.Args["transactionId"].(int)
		shippoRateID := params.Args["shippoRateId"].(string)

		tempTransaction, err := transactionLoader.Load(params.Context, dataloaders.IntKey(transactionId))()
		if err != nil {
			return nil, err
		}
		transaction := tempTransaction.(*db.Transaction)

		rate, err := shippoClient.RetrieveRate(shippoRateID)
		if err != nil || rate == nil {
			return nil, &core.WrappedError{
				Message:       "Could not retrieve shipping rate.",
				InternalError: err,
			}
		}

		label, err := shippoClient.PurchaseShippingLabel(&models.TransactionInput{
			Rate:          rate.ObjectID,
			LabelFileType: models.LabelFileTypePDF,
			Async:         false,
		})
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not purchase shipping label.",
				InternalError: err,
			}
		}

		if label.Status == "ERROR" {
			message := "Could not purchase shipping label."
			if len(label.Messages) > 0 {
				message = label.Messages[0].Text
			}

			return nil, &core.WrappedError{
				Message: message,
			}
		}

		transaction.ShippoTransactionID = label.ObjectID

		status := &db.TransactionStatus{
			TransactionID: transaction.ID,
			Status:        "SHIPPED",
			Carrier:       rate.Provider,
			TrackingID:    label.TrackingNumber,
		}

		if err := database.Update(transaction); err != nil {
			fmt.Println("Failed to update transaction with shippo transaction id.")
			fmt.Println(err)
		}

		if err := database.Insert(status); err != nil {
			fmt.Println("Failed to create transaction status with tracking number.")
			fmt.Println(err)
		}

		if transaction.UserID > 0 {
			toSend, err := email.NewShippedEmail(label.TrackingURLProvider)
			if err != nil {
				fmt.Println("Failed create shipped email.")
				fmt.Println(err)
			}

			tmpUser, err := userLoader.Load(params.Context, dataloaders.IntKey(transaction.UserID))()
			if err != nil {
				fmt.Println("Failed to find user to email shipped order to.")
				fmt.Println(err)
			} else {
				user := tmpUser.(*db.User)
				err = emailClient.SendMail(user.Email, "Your order has shipped.", toSend)
				if err != nil {
					fmt.Println("Failed to send purchase email.")
					fmt.Println(err)
				}
			}
		}

		return map[string]interface{}{
			"id":       label.ObjectID,
			"labelUrl": label.LabelURL,
		}, nil
	},
}
