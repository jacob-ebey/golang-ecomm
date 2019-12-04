package schema

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/jacob-ebey/golang-ecomm/auth"
	"github.com/jacob-ebey/golang-ecomm/dataloaders"
)

var MarkdownScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:         "Markdown",
	Serialize:    graphql.String.Serialize,
	ParseValue:   graphql.String.ParseValue,
	ParseLiteral: graphql.String.ParseLiteral,
})

func PrettyPrint(obj interface{}) {
	b, _ := json.Marshal(obj)
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	fmt.Printf("%s\n", out.Bytes())
}

// Parses an optional string from the graphql args.
func OptionalString(args map[string]interface{}, arg string) *string {
	value, success := args[arg].(string)

	if !success {
		return nil
	}

	return &value
}

func OptionalInt(args map[string]interface{}, arg string) *int {
	value, success := args[arg].(int)

	if !success {
		return nil
	}

	return &value
}

func OptionalFloat(args map[string]interface{}, arg string) *float64 {
	value, success := args[arg].(float64)

	if !success {
		return nil
	}

	return &value
}

// Converts the input object to the output object via json marshaling.
func ConvertObject(input interface{}, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, output); err != nil {
		return err
	}

	return nil
}

// Options used to create a new pagination field.
type PaginationFieldOpts struct {
	Description string
	Type        graphql.Output
	Dataloader  string
	Auth        bool
	AuthRole    string
}

// Creates a new pagination field that wraps a dataloader that uses dataloaders.PaginationKey.
func NewPaginationField(options PaginationFieldOpts) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(options.Type),
		Description: options.Description,
		Args: graphql.FieldConfigArgument{
			"skip": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"limit": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			loader := params.Context.Value(options.Dataloader).(*dataloader.Loader)

			claims := params.Context.Value("claims").(*auth.Claims)

			if options.Auth || options.AuthRole != "" {
				if claims == nil {
					return nil, auth.NotAuthenticatedError
				}

				if options.AuthRole != "" && claims.Role != options.AuthRole {
					return nil, auth.NotAuthorizedError
				}
			}

			skip, _ := params.Args["skip"].(int)
			limit, _ := params.Args["limit"].(int)

			if skip < 0 {
				skip = 0
			}

			if limit <= 0 {
				limit = 20
			}

			thunk := loader.Load(params.Context, dataloaders.PaginationKey{
				Skip:  skip,
				Limit: limit,
			})

			return func() (interface{}, error) {
				return thunk()
			}, nil
		},
	}
}
