package mq

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

type middlewareKey struct{}

func GetMiddlewareFromContext(ctx context.Context) []middleware.Middleware {
	ms := ctx.Value(middlewareKey{})
	if ms != nil {
		return ms.([]middleware.Middleware)
	}
	return nil
}

func MiddlewareWithContext(ctx context.Context, list ...middleware.Middleware) context.Context {
	ms := GetMiddlewareFromContext(ctx)
	return context.WithValue(ctx, middlewareKey{}, append(ms, list...))
}
