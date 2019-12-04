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
	"golang.org/x/crypto/bcrypt"
)

var AuthResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthResponse",
	Fields: graphql.Fields{
		"refreshToken": &graphql.Field{
			Type:        graphql.String,
			Description: "The token to include in the 'Authorization' header with the 'refreshToken' mutation. Example: 'Bearer yourtokenhere'",
		},
		"token": &graphql.Field{
			Type:        graphql.String,
			Description: "The auth token to include in the 'Authorization' header. Example: 'Bearer yourtokenhere'",
		},
	},
})

var LoginError = fmt.Errorf("Email or password is invalid.")
var UserExistsError = fmt.Errorf("User with the provided email already exists.")
var PasswordsDoNotMatchError = fmt.Errorf("The provided passwords do not match.")

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var SignUpField = &graphql.Field{
	Type:        AuthResponseType,
	Description: "Sign up with email and password.",
	Args: graphql.FieldConfigArgument{
		"email": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"password": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"confirmPassword": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		database := params.Context.Value("database").(*pg.DB)
		tokenGenerator := params.Context.Value("auth").(auth.TokenGenerator)

		email := params.Args["email"].(string)
		password := params.Args["password"].(string)
		confirmPassword := params.Args["confirmPassword"].(string)

		if password != confirmPassword {
			return nil, PasswordsDoNotMatchError
		}

		user := db.User{}
		if exists, err := database.
			Model(&user).
			Where("email = ?", email).
			Exists(); exists || err != nil {
			return nil, &core.WrappedError{
				Message:       "Email already in use.",
				InternalError: err,
			}
		}

		hashedPassword, err := hashPassword(password)
		if err != nil {
			return nil, err
		}

		user = db.User{
			Email:    email,
			Password: hashedPassword,
		}

		if err := database.Insert(&user); err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create user.",
				InternalError: err,
			}
		}

		token, err := tokenGenerator.GenerateToken(params.Context, auth.Claims{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		}, 60) // 1 hour expiration
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create token.",
				InternalError: err,
			}
		}

		refreshToken, err := tokenGenerator.GenerateToken(params.Context, auth.Claims{
			ID: user.ID,
		}, 1440) // 1 day expiration
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create refresh token.",
				InternalError: err,
			}
		}

		return map[string]string{
			"token":        token,
			"refreshToken": refreshToken,
		}, nil
	},
}

var SignInField = &graphql.Field{
	Type:        AuthResponseType,
	Description: "Sign in with email and password.",
	Args: graphql.FieldConfigArgument{
		"email": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"password": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		database := params.Context.Value("database").(*pg.DB)
		tokenGenerator := params.Context.Value("auth").(auth.TokenGenerator)

		email := params.Args["email"].(string)
		password := params.Args["password"].(string)

		user := db.User{}
		if err := database.
			Model(&user).
			Where("email = ?", email).
			Select(); err != nil {
			return nil, LoginError
		}

		if !checkPasswordHash(password, user.Password) {
			return nil, LoginError
		}

		token, err := tokenGenerator.GenerateToken(params.Context, auth.Claims{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		}, 60) // 1 hour expiration
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create token.",
				InternalError: err,
			}
		}

		refreshToken, err := tokenGenerator.GenerateToken(params.Context, auth.Claims{
			ID: user.ID,
		}, 1440) // 1 day expiration
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create refresh token.",
				InternalError: err,
			}
		}

		return map[string]string{
			"token":        token,
			"refreshToken": refreshToken,
		}, nil
	},
}

var RefreshTokenField = &graphql.Field{
	Type:        AuthResponseType,
	Description: "Refrsh tokens for a user.",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		userLoader := params.Context.Value("user").(*dataloader.Loader)
		tokenGenerator := params.Context.Value("auth").(auth.TokenGenerator)
		claims := params.Context.Value("claims").(*auth.Claims)

		if claims == nil {
			return nil, auth.NotAuthenticatedError
		}

		user, err := userLoader.Load(params.Context, dataloaders.IntKey(claims.ID))()
		if err != nil {
			return nil, err
		}

		fullUser, ok := user.(*db.User)
		if !ok {
			return nil, &core.WrappedError{
				Message: "Failed to load user.",
			}
		}

		token, err := tokenGenerator.GenerateToken(params.Context, auth.Claims{
			ID:    fullUser.ID,
			Email: fullUser.Email,
			Role:  fullUser.Role,
		}, 60) // 1 hour expiration
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create token.",
				InternalError: err,
			}
		}

		refreshToken, err := tokenGenerator.GenerateToken(params.Context, auth.Claims{
			ID: fullUser.ID,
		}, 1440) // 1 day expiration
		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Could not create refresh token.",
				InternalError: err,
			}
		}

		return map[string]string{
			"token":        token,
			"refreshToken": refreshToken,
		}, nil
	},
}
