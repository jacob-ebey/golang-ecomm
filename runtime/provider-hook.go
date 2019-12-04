package runtime

import (
	"context"

	core "github.com/jacob-ebey/graphql-core"
)

type ProviderHook struct {
	Key   string
	Value interface{}
}

func NewProviderHook(key string, value interface{}) *ProviderHook {
	return &ProviderHook{
		Key:   key,
		Value: value,
	}
}

func (hook *ProviderHook) PreExecute(ctx context.Context, req core.GraphQLRequest) context.Context {
	return context.WithValue(ctx, hook.Key, hook.Value)
}
