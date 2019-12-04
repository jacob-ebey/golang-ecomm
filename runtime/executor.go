package runtime

import (
	"strconv"

	"github.com/braintree-go/braintree-go"
	"github.com/graphql-go/graphql"
	"github.com/jacob-ebey/go-shippo"

	core "github.com/jacob-ebey/graphql-core"
	storage "github.com/jacob-ebey/now-storage-go"

	"github.com/jacob-ebey/golang-ecomm/auth"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	"github.com/jacob-ebey/golang-ecomm/db"
	"github.com/jacob-ebey/golang-ecomm/email"
	"github.com/jacob-ebey/golang-ecomm/schema"
	"github.com/jacob-ebey/golang-ecomm/services"
)

type NewExecutorOpts struct {
	RunBefore []core.PreExecuteHook
	RunAfter  []core.PostExecuteHook
}

func NewExecutor(opts NewExecutorOpts) (*core.GraphQLExecutor, error) {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    schema.QueryType,
		Mutation: schema.MutationType,
	})

	if err != nil {
		return nil, &core.WrappedError{
			Message:       "Could not create schema.",
			InternalError: err,
		}
	}

	baseUrlHook := NewProviderHook("baseUrl", BaseUrl())

	authHook := &auth.JwtAuthHook{
		JwtHeader: "Bearer",
		JwtSecret: JwtSecret(),
	}

	databaseHook, err := db.NewDatabaseHook(GetPgOptions())
	if err != nil {
		return nil, &core.WrappedError{
			Message:       "Failed to create database hook",
			InternalError: err,
		}
	}

	avataxHook := NewProviderHook("avatax", GetAvatax())

	shippoHook := NewProviderHook("shippo", shippo.NewClient(ShippoPrivateToken()))

	braintreeConfig := Braintree()
	braintreeEnvironment := braintree.Production
	if IsDevelopment() {
		braintreeEnvironment = braintree.Sandbox
	}
	braintreeHook := NewProviderHook(
		"braintree",
		braintree.New(braintreeEnvironment, braintreeConfig.MerchantID, braintreeConfig.PublicKey, braintreeConfig.PrivateKey))

	smtpConfig := Smtp()
	smtpPort := strconv.Itoa(smtpConfig.Port)
	emailHook := email.NewSmtpClient(smtpConfig.From, smtpConfig.Host+":"+smtpPort, email.LoginAuth(smtpConfig.Username, smtpConfig.Password))

	nowStorageHook := NewProviderHook("nowStorage", &storage.Client{
		Token:          ZeitToken(),
		DeploymentName: "golang-ecomm",
	})

	return &core.GraphQLExecutor{
		Schema: schema,
		RunBefore: append(opts.RunBefore,
			baseUrlHook,
			authHook,
			databaseHook,
			dataloaders.HooksDataloader,
			avataxHook,
			shippoHook,
			services.ValidateAddressWithShippo,
			services.ResizeImage,
			braintreeHook,
			emailHook,
			nowStorageHook,
		),
		RunAfter: append(opts.RunAfter,
			dataloaders.HooksDataloader,
		),
	}, nil
}
