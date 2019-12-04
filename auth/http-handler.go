package auth

import (
	"context"
	"net/http"

	core "github.com/jacob-ebey/graphql-core"
)

type HttpHeaderHook struct {
	Source string
	Dest   string
}

func (hook *HttpHeaderHook) PreExecute(ctx context.Context, req core.GraphQLRequest) context.Context {
	request, ok := ctx.Value("request").(*http.Request)

	if !ok || request == nil {
		return ctx
	}

	return context.WithValue(ctx, hook.Dest, request.Header.Get(hook.Source))
}
