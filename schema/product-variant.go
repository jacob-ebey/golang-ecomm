package schema

import (
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	core "github.com/jacob-ebey/graphql-core"

	"github.com/jacob-ebey/golang-ecomm/auth"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	"github.com/jacob-ebey/golang-ecomm/db"
	"github.com/jacob-ebey/golang-ecomm/services"
	"github.com/jacob-ebey/golang-ecomm/utilities"
)

var ProductVariantType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ProductVariant",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					product := params.Context.Value("product").(*dataloader.Loader)

					variant := params.Source.(*db.ProductVariant)

					if variant.Name == "" {
						thunk := product.Load(params.Context, dataloaders.IntKey(variant.ProductID))

						return func() (interface{}, error) {
							result, err := thunk()

							if err != nil || result.(*db.Product) == nil {
								return nil, err
							}

							return result.(*db.Product).Name, nil
						}, nil
					}

					return variant.Name, nil
				},
			},
			"price": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"length": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"width": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"height": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"weight": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"selectedOptions": &graphql.Field{
				Type: graphql.NewList(ProductOptionValueType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					productVariantOptions := params.Context.Value("productVariantOptions").(*dataloader.Loader)

					variant := params.Source.(*db.ProductVariant)

					thunk := productVariantOptions.Load(params.Context, dataloaders.IntKey(variant.ID))

					return func() (interface{}, error) {
						return thunk()
					}, nil
				},
			},
			"images": &graphql.Field{
				Type:        graphql.NewList(ImageType),
				Description: "The images for the product variant.",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					productVariantImagesLoader := params.Context.Value("productVariantImages").(*dataloader.Loader)

					variant := params.Source.(*db.ProductVariant)

					thunk := productVariantImagesLoader.Load(params.Context, dataloaders.IntKey(variant.ID))

					return func() (interface{}, error) {
						return thunk()
					}, nil
				},
			},
		},
	},
)

var ProductVariantsByIdsField = &graphql.Field{
	Type:        graphql.NewList(ProductVariantType),
	Description: "Get product variants by IDs.",
	Args: graphql.FieldConfigArgument{
		"variantIds": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int))),
			Description: "The product variants to retrieve.",
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		productVariant := params.Context.Value("productVariant").(*dataloader.Loader)

		idParams := params.Args["variantIds"].([]interface{})
		ids := make(dataloader.Keys, len(idParams))
		for index, param := range idParams {
			ids[index] = dataloaders.IntKey(param.(int))
		}

		thunk := productVariant.LoadMany(params.Context, ids)

		return func() (interface{}, error) {
			variants, errs := thunk()

			if len(errs) > 0 {
				return variants, dataloaders.HandleErrors(errs)
			}

			return variants, nil
		}, nil
	},
}

var ProductVariantBySelectedOptionsField = &graphql.Field{
	Type:        ProductVariantType,
	Description: "Get a product variant for a product by selected options.",
	Args: graphql.FieldConfigArgument{
		"productId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"selectedOptions": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int))),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		productVariantBySelectedOptions := params.Context.Value("productVariantBySelectedOptions").(*dataloader.Loader)

		id := dataloaders.SelectedOptionsKey{}
		if err := ConvertObject(params.Args, &id); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not convert arguments.",
				InternalError: err,
			}
		}

		thunk := productVariantBySelectedOptions.Load(params.Context, id)

		return func() (interface{}, error) {
			return thunk()
		}, nil
	},
}

var RemoveProductVariantField = &graphql.Field{
	Type:        ProductVariantType,
	Description: "Remove a product variant.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
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

		id := params.Args["id"].(int)

		toDelete := db.ProductVariant{}
		if err := database.Model(&toDelete).Where("id = ?", id).Select(); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not retrieve product variant to remove.",
				InternalError: err,
			}
		}

		if err := database.Delete(&toDelete); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not remove product variant.",
				InternalError: err,
			}
		}

		return &toDelete, nil
	},
}

var CreateProductVariantInputSchema = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateProductVariantInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"price": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "The price in cents (¢).",
		},
		"length": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.Float),
			Description: "The length in inches.",
		},
		"width": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.Float),
			Description: "The width in inches.",
		},
		"height": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.Float),
			Description: "The height in inches.",
		},
		"weight": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.Float),
			Description: "The weight in ounces.",
		},
	},
})

var CreateProductVariantField = &graphql.Field{
	Type:        graphql.NewList(ProductVariantType),
	Description: "Create a variant for a product.",
	Args: graphql.FieldConfigArgument{
		"productId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(CreateProductVariantInputSchema),
		},
		"selectedProductOptionValues": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int))),
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

		input := db.ProductVariant{}
		if err := ConvertObject(params.Args["input"], &input); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not decode input.",
				InternalError: err,
			}
		}
		input.ProductID = params.Args["productId"].(int)

		selectedProductOptionValues := []int{}
		if err := ConvertObject(params.Args["selectedProductOptionValues"], &selectedProductOptionValues); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not decode input.",
				InternalError: err,
			}
		}

		if err := database.Insert(&input); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create product variant.",
				InternalError: err,
			}
		}

		inserted := []db.ProductVariantOption{}
		var err error = nil
		for _, selectedOptionValueId := range selectedProductOptionValues {
			selected := db.ProductVariantOption{
				ProductOptionValueID: selectedOptionValueId,
				ProductVariantID:     input.ID,
				ProductID:            input.ProductID,
			}

			err = database.Insert(&selected)
			if err != nil {
				break
			}

			inserted = append(inserted, selected)
		}

		if err != nil {
			for _, toDelete := range inserted {
				database.Delete(&toDelete)      // TODO: Log error
				database.ForceDelete(&toDelete) // TODO: Log error
			}

			return nil, &core.WrappedError{
				Message:       "Could not create product variant.",
				InternalError: err,
			}
		}

		return &input, nil
	},
}

var UpdateProductVariantInputSchema = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateProductVariantInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "An optional name. If none is provided, the name of the product is used.",
		},
		"price": &graphql.InputObjectFieldConfig{
			Type:        graphql.Int,
			Description: "The price in cents (¢).",
		},
		"length": &graphql.InputObjectFieldConfig{
			Type:        graphql.Float,
			Description: "The length in inches.",
		},
		"width": &graphql.InputObjectFieldConfig{
			Type:        graphql.Float,
			Description: "The width in inches.",
		},
		"height": &graphql.InputObjectFieldConfig{
			Type:        graphql.Float,
			Description: "The height in inches.",
		},
		"weight": &graphql.InputObjectFieldConfig{
			Type:        graphql.Float,
			Description: "The weight in ounces.",
		},
	},
})

var UpdateProductVariantField = &graphql.Field{
	Type:        ProductVariantType,
	Description: "Update a variant for a product.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(UpdateProductVariantInputSchema),
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

		id := params.Args["id"].(int)
		input := params.Args["input"].(map[string]interface{})
		name := OptionalString(input, "name")
		price := OptionalInt(input, "price")
		length := OptionalFloat(input, "length")
		width := OptionalFloat(input, "width")
		height := OptionalFloat(input, "height")
		weight := OptionalFloat(input, "weight")

		result := db.ProductVariant{ID: id}
		if err := database.Select(&result); err != nil {
			if err != nil {
				return nil, &core.WrappedError{
					Message:       "Could not find product variant to update.",
					InternalError: err,
				}
			}
		}

		if name != nil {
			result.Name = strings.TrimSpace(*name)
		}
		if price != nil {
			result.Price = *price
		}
		if length != nil {
			result.Length = *length
		}
		if width != nil {
			result.Width = *width
		}
		if height != nil {
			result.Height = *height
		}
		if weight != nil {
			result.Weight = *weight
		}

		if err := database.Update(&result); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not update product variant.",
				InternalError: err,
			}
		}

		return &result, nil
	},
}

var CreateProductPermutationsInputSchema = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateProductPermutationsInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"price": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "The price in cents (¢).",
		},
		"length": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.Float),
			Description: "The length in inches.",
		},
		"width": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.Float),
			Description: "The width in inches.",
		},
		"height": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.Float),
			Description: "The height in inches.",
		},
		"weight": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.Float),
			Description: "The weight in ounces.",
		},
	},
})

var CreateProductPermutationsField = &graphql.Field{
	Type:        graphql.NewList(ProductVariantType),
	Description: "Create all the variant permutations for a product based on it's options. There must not be any existing product variants for the product, otheriwse this will fail.",
	Args: graphql.FieldConfigArgument{
		"productId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(CreateProductPermutationsInputSchema),
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

		input := db.ProductVariant{}
		if err := ConvertObject(params.Args["input"], &input); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not decode input.",
				InternalError: err,
			}
		}
		input.ProductID = params.Args["productId"].(int)

		if exists, err := database.
			Model(&db.ProductVariant{}).
			Where("product_id = ?", input.ProductID).
			Exists(); exists || err != nil {
			return nil, &core.WrappedError{
				Message:       "Product already contains product vairants.",
				InternalError: err,
			}
		}

		options := []db.ProductOption{}
		if err := database.
			Model(&options).
			Column("product_option.*").
			Where("product_option.product_id = ?", input.ProductID).
			Relation("Values").
			Select(); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not get product options.",
				InternalError: err,
			}
		}

		// Get option value ids for permutation
		optionIDs := make([][]int, len(options))
		for optionIndex, option := range options {
			optionIDs[optionIndex] = make([]int, len(option.Values))

			for valueIndex, value := range option.Values {
				optionIDs[optionIndex][valueIndex] = value.ID
			}
		}

		permutations := utilities.Permutate(optionIDs)

		createdVariants := []*db.ProductVariant{}
		var err error
		for permutaitonIndex, permutation := range permutations {
			variant := db.ProductVariant{
				Price:     input.Price,
				Length:    input.Length,
				Width:     input.Width,
				Height:    input.Height,
				Weight:    input.Weight,
				ProductID: input.ProductID,
			}
			err = database.Insert(&variant)
			if err != nil {
				break
			}

			createdVariants = append(createdVariants, &variant)

			options := make([]*db.ProductVariantOption, len(permutation))
			for index, optionID := range permutation {
				options[index] = &db.ProductVariantOption{
					ProductOptionValueID: optionID,
					ProductVariantID:     variant.ID,
					ProductID:            input.ProductID,
				}
			}

			err = database.Insert(&options)
			if err != nil {
				break
			}

			createdVariants[permutaitonIndex].SelectedOptions = options
		}

		// Rollback
		if err != nil {
			for _, variant := range createdVariants {
				if variant.SelectedOptions != nil {
					for _, selectedOption := range variant.SelectedOptions {
						database.Delete(&selectedOption)
						database.ForceDelete(&selectedOption)
					}
				}

				database.Delete(variant)
				database.ForceDelete(variant)
			}

			return nil, &core.WrappedError{
				Message:       "Failed to create permutations.",
				InternalError: err,
			}
		}

		return createdVariants, nil
	},
}

var AddProductVariantImageField = &graphql.Field{
	Type:        ImageType,
	Description: "Add a product variant image.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"image": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(UploadScalar),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		database := params.Context.Value("database").(*pg.DB)
		productVariantLoader := params.Context.Value("productVariant").(*dataloader.Loader)
		imageResizer := params.Context.Value("imageResizer").(services.ImageResizer)

		claims := params.Context.Value("claims").(*auth.Claims)
		if claims == nil {
			return nil, auth.NotAuthenticatedError
		}
		if claims.Role != "ADMIN" {
			return nil, auth.NotAuthorizedError
		}

		id := params.Args["id"].(int)

		file := params.Args["image"].(*core.MultipartFile)
		defer file.File.Close()

		_, err := productVariantLoader.Load(params.Context, dataloaders.IntKey(id))()
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not retrieve product variant to add image.",
				InternalError: err,
			}
		}

		image, err := imageResizer.ResizeImage(params.Context, file)
		if err != nil {
			return nil, err
		}

		if err := database.Insert(image); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not save image.",
				InternalError: err,
			}
		}

		productVariantImage := db.ProductVariantImage{
			ProductVariantID: id,
			ImageID:          image.ID,
		}
		if err := database.Insert(&productVariantImage); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not save product variant image.",
				InternalError: err,
			}
		}

		return image, nil
	},
}

var RemoveProductVariantImageField = &graphql.Field{
	Type:        graphql.Boolean,
	Description: "Remove a product variant image.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"imageId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		database := params.Context.Value("database").(*pg.DB)

		id := params.Args["id"].(int)
		image := params.Args["imageId"].(int)

		if _, err := database.
			Model(&db.ProductVariantImage{}).
			Where("product_variant_id = ?", id).
			Where("image_id = ?", image).
			Delete(); err != nil {
			return false, &core.WrappedError{
				Message:       "Could not remove image from product variant.",
				InternalError: err,
			}
		}

		return true, nil
	},
}
