package schema

import (
	"fmt"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"

	"github.com/jacob-ebey/golang-ecomm/auth"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
	"github.com/jacob-ebey/golang-ecomm/db"
	"github.com/jacob-ebey/golang-ecomm/services"
	core "github.com/jacob-ebey/graphql-core"
)

const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64
const (
	MaxInt = 1<<(UintSize-1) - 1 // 1<<31 - 1 or 1<<63 - 1
	MinInt = -MaxInt - 1         // -1 << 31 or -1 << 63
)

type ProductPriceRange struct {
	Min int
	Max int
}

var ProductType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "An ID unique to the Product type. May conflict with other types.",
			},
			"slug": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "A unique identifier for the product used in places like URL's.",
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The name of the product. Captialize things here.",
			},
			"description": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "A short description of the product. Try to keep this under 200 characters for the best UX.",
			},
			"details": &graphql.Field{
				Type:        MarkdownScalar,
				Description: "More in-depth details about the product in Markdown format.",
			},
			"published": &graphql.Field{
				Type: graphql.Boolean,
			},
			"priceRange": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "ProductPriceRange",
					Fields: graphql.Fields{
						"min": &graphql.Field{Type: graphql.Int},
						"max": &graphql.Field{Type: graphql.Int},
					},
				}),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					productVariants := params.Context.Value("productVariants").(*dataloader.Loader)

					product := params.Source.(*db.Product)

					thunk := productVariants.Load(params.Context, dataloaders.IntKey(product.ID))

					return func() (interface{}, error) {
						results, err := thunk()

						if err != nil {
							return nil, &core.WrappedError{
								Message:       "Could not get price range for product.",
								InternalError: err,
							}
						}

						variants := results.([]*db.ProductVariant)

						if variants == nil || len(variants) == 0 {
							return nil, nil
						}

						minValue := MaxInt
						maxValue := MinInt
						for _, variant := range variants {
							if variant.Price < minValue {
								minValue = variant.Price
							}

							if variant.Price > maxValue {
								maxValue = variant.Price
							}
						}

						return &ProductPriceRange{
							Max: maxValue,
							Min: minValue,
						}, nil
					}, nil
				},
			},
			"options": &graphql.Field{
				Type:        graphql.NewList(ProductOptionType),
				Description: "The configurable options for the product.",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					productOptions := params.Context.Value("productOptions").(*dataloader.Loader)

					product := params.Source.(*db.Product)

					thunk := productOptions.Load(params.Context, dataloaders.IntKey(product.ID))

					return func() (interface{}, error) {
						return thunk()
					}, nil
				},
			},
			"variants": &graphql.Field{
				Type:        graphql.NewList(ProductVariantType),
				Description: "The purchasable variants of the product.",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					productVariants := params.Context.Value("productVariants").(*dataloader.Loader)

					product := params.Source.(*db.Product)

					thunk := productVariants.Load(params.Context, dataloaders.IntKey(product.ID))

					return func() (interface{}, error) {
						return thunk()
					}, nil
				},
			},
			"images": &graphql.Field{
				Type:        graphql.NewList(ImageType),
				Description: "The images for the product.",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					productImagesLoader := params.Context.Value("productImages").(*dataloader.Loader)

					product := params.Source.(*db.Product)

					thunk := productImagesLoader.Load(params.Context, dataloaders.IntKey(product.ID))

					return func() (interface{}, error) {
						return thunk()
					}, nil
				},
			},
		},
	},
)

var ProductField = &graphql.Field{
	Type:        ProductType,
	Description: "Get a product by ID.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		product := params.Context.Value("product").(*dataloader.Loader)

		id := params.Args["id"].(int)

		thunk := product.Load(params.Context, dataloaders.IntKey(id))

		return func() (interface{}, error) {
			return thunk()
		}, nil
	},
}

var ProductBySlugField = &graphql.Field{
	Type:        ProductType,
	Description: "Get a product from the catalog by slug.",
	Args: graphql.FieldConfigArgument{
		"slug": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		productBySlug := params.Context.Value("productBySlug").(*dataloader.Loader)

		slug := params.Args["slug"].(string)

		thunk := productBySlug.Load(params.Context, dataloader.StringKey(slug))

		return func() (interface{}, error) {
			return thunk()
		}, nil
	},
}

var CreateProductInputSchema = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateProductInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"slug": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "A unique identifier for the product used in places like URL's.",
		},
		"name": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The name of the product. Captialize things here.",
		},
		"description": &graphql.InputObjectFieldConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "A short description of the product. Try to keep this under 200 characters for the best UX.",
		},
		"details": &graphql.InputObjectFieldConfig{
			Type:        MarkdownScalar,
			Description: "More in-depth details about the product in Markdown format.",
		},
	},
})

var UpdateProductInputSchema = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateProductInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "The name of the product. Captialize things here.",
		},
		"description": &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "A short description of the product. Try to keep this under 200 characters for the best UX.",
		},
		"details": &graphql.InputObjectFieldConfig{
			Type:        MarkdownScalar,
			Description: "More in-depth details about the product in Markdown format.",
		},
	},
})

var CreateProductDraftField = &graphql.Field{
	Type:        ProductType,
	Description: "Create a new product draft.",
	Args: graphql.FieldConfigArgument{
		"product": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(CreateProductInputSchema),
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

		product := db.Product{}
		if err := ConvertObject(params.Args["product"], &product); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not convert arguments.",
				InternalError: err,
			}
		}

		product.Slug = strings.TrimSpace(product.Slug)
		product.Name = strings.TrimSpace(product.Name)
		product.Description = strings.TrimSpace(product.Description)

		if product.Slug == "" {
			return nil, fmt.Errorf("Slug is required.")
		}

		if product.Name == "" {
			return nil, fmt.Errorf("Name is required.")
		}

		if product.Description == "" {
			return nil, fmt.Errorf("Description is required.")
		}

		if err := database.Insert(&product); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create product draft.",
				InternalError: err,
			}
		}

		return &product, nil
	},
}

var PublishProductField = &graphql.Field{
	Type:        ProductType,
	Description: "Publish or un-publish a product from the catalog.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"published": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Boolean),
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
		published := params.Args["published"].(bool)

		result := db.Product{ID: id}

		if err := database.Select(&result); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not find product draft to publish.",
				InternalError: err,
			}
		}

		result.Published = published

		if err := database.Update(&result); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not publish product.",
				InternalError: err,
			}
		}

		return &result, nil
	},
}

var UpdateProductField = &graphql.Field{
	Type:        ProductType,
	Description: "Update a product.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"product": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(UpdateProductInputSchema),
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
		productInput := params.Args["product"].(map[string]interface{})
		name := OptionalString(productInput, "name")
		description := OptionalString(productInput, "description")
		details := OptionalString(productInput, "details")

		result := db.Product{ID: id}
		if err := database.Select(&result); err != nil {
			if err != nil {
				return nil, &core.WrappedError{
					Message:       "Could not find product to update.",
					InternalError: err,
				}
			}
		}

		if name != nil {
			result.Name = strings.TrimSpace(*name)
		}
		if description != nil {
			result.Description = strings.TrimSpace(*description)
		}
		if details != nil {
			result.Details = strings.TrimSpace(*details)
		}

		if err := database.Update(&result); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not update product.",
				InternalError: err,
			}
		}

		return &result, nil
	},
}

var AddProductImageField = &graphql.Field{
	Type:        ImageType,
	Description: "Add a product image.",
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
		productLoader := params.Context.Value("product").(*dataloader.Loader)
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

		_, err := productLoader.Load(params.Context, dataloaders.IntKey(id))()
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not retrieve product to add image.",
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

		productImage := db.ProductImage{
			ProductID: id,
			ImageID:   image.ID,
		}
		if err := database.Insert(&productImage); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not save product image.",
				InternalError: err,
			}
		}

		return image, nil
	},
}

var RemoveProductImageField = &graphql.Field{
	Type:        graphql.Boolean,
	Description: "Remove a product image.",
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
			Model(&db.ProductImage{}).
			Where("product_id = ?", id).
			Where("image_id = ?", image).
			Delete(); err != nil {
			return false, &core.WrappedError{
				Message:       "Could not remove image from product.",
				InternalError: err,
			}
		}

		return true, nil
	},
}
