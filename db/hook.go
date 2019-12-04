package db

import (
	"context"
	"reflect"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"

	core "github.com/jacob-ebey/graphql-core"
)

type DatabaseHook struct {
	Database *pg.DB
}

func NewDatabaseHook(options *pg.Options) (*DatabaseHook, error) {
	database := pg.Connect(options)

	types := []interface{}{
		(*User)(nil),
		(*Image)(nil),
		(*Address)(nil),
		(*Product)(nil),
		(*ProductImage)(nil),
		(*ProductOption)(nil),
		(*ProductOptionValue)(nil),
		(*ProductVariant)(nil),
		(*ProductVariantOption)(nil),
		(*ProductVariantImage)(nil),
		(*Transaction)(nil),
		(*TransactionAddressInfo)(nil),
		(*TransactionLineItem)(nil),
		(*TransactionStatus)(nil),
	}

	for _, model := range types {
		err := database.CreateTable(model, &orm.CreateTableOptions{
			// Temp:          dev,
			IfNotExists:   true,
			FKConstraints: true,
		})

		if err != nil {
			return nil, &core.WrappedError{
				Message:       "Failed to create table for type \"" + reflect.TypeOf(model).String() + "\".",
				InternalError: err,
			}
		}
	}

	return &DatabaseHook{
		Database: database,
	}, nil
}

func (hook *DatabaseHook) PreExecute(ctx context.Context, req core.GraphQLRequest) context.Context {
	return context.WithValue(ctx, "database", hook.Database)
}
