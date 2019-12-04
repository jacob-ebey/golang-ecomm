package schema

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/jacob-ebey/golang-ecomm/auth"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	"github.com/jacob-ebey/golang-ecomm/db"
	core "github.com/jacob-ebey/graphql-core"
	"strings"
)

var ProductOptionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ProductOption",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"label": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"values": &graphql.Field{
				Type: graphql.NewList(ProductOptionValueType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					productOptionValues := params.Context.Value("productOptionValues").(*dataloader.Loader)

					option := params.Source.(*db.ProductOption)

					thunk := productOptionValues.Load(params.Context, dataloaders.IntKey(option.ID))

					return func() (interface{}, error) {
						return thunk()
					}, nil
				},
			},
		},
	},
)

var ProductOptionValueType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ProductOptionValue",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"value": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"productOptionId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
	},
)

var CreateProductOptionValueInputSchema = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateProductOptionValueInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"value": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "A value for the option. Make these human readable.",
		},
	},
})

var CreateProductOptionInputSchema = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateProductOptionInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"label": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The label of the option.",
		},
		"values": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(CreateProductOptionValueInputSchema))),
			Description: "The values for the option.",
		},
	},
})

var CreateProductOptionField = &graphql.Field{
	Type:        ProductOptionType,
	Description: "Create a new product draft.",
	Args: graphql.FieldConfigArgument{
		"productId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"option": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(CreateProductOptionInputSchema),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		database := params.Context.Value("database").(*pg.DB)

		claims := params.Context.Value("claims").(*auth.Claims)

		if claims == nil {
			return nil, auth.NotAuthenticatedError
		}

		if claims.Role != "ADMIN" {
			return nil, auth.NotAuthorizedError
		}

		productID := params.Args["productId"].(int)

		label := params.Args["option"].(map[string]interface{})["label"].(string)

		inputValues := make([]db.ProductOptionValue, len(params.Args["option"].(map[string]interface{})["values"].([]interface{})))
		if err := ConvertObject(params.Args["option"].(map[string]interface{})["values"], &inputValues); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not convert arguments.",
				InternalError: err,
			}
		}

		option := db.ProductOption{
			Label:     label,
			ProductID: productID,
		}

		if err := database.Insert(&option); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create product draft.",
				InternalError: err,
			}
		}

		var err error

		for index, inputValue := range inputValues {
			inputValues[index].Value = strings.TrimSpace(inputValue.Value)
			if inputValues[index].Value == "" {
				err = fmt.Errorf("Values must not be empty.")
				break
			}

			inputValues[index].ProductOptionID = option.ID

			if err := database.Insert(&inputValues[index]); err != nil {
				err = &core.WrappedError{
					Message:       "Could not create product draft.",
					InternalError: err,
				}
				break
			}
		}

		if err != nil {
			database.Delete(&option)
			database.ForceDelete(&option)

			for _, value := range inputValues {
				database.Delete(&value)
				database.ForceDelete(&value)
			}

			return nil, err
		}

		return &option, nil
	},
}

var RemoveProductOptionField = &graphql.Field{
	Type:        ProductOptionType,
	Description: "Create a new product draft.",
	Args: graphql.FieldConfigArgument{
		"productId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"productOptionId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		database := params.Context.Value("database").(*pg.DB)

		claims := params.Context.Value("claims").(*auth.Claims)

		if claims == nil {
			return nil, auth.NotAuthenticatedError
		}

		if claims.Role != "ADMIN" {
			return nil, auth.NotAuthorizedError
		}

		productID := params.Args["productId"].(int)
		productOptionID := params.Args["productOptionId"].(int)

		option := db.ProductOption{}
		if err := database.
			Model(&option).
			Where("product_id = ?", productID).
			Where("id = ?", productOptionID).
			Select(); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not retrieve product option to remove.",
				InternalError: err,
			}
		}

		if err := database.Delete(&option); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not remove product option.",
				InternalError: err,
			}
		}

		return &option, nil
	},
}
