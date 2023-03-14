package storage

import (
	"context"
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

func NewCQRSContext(parent context.Context) context.Context {
	return context.WithValue(parent, contextKey, true)
}

func IsCQRSContext(ctx context.Context) bool {
	ok := ctx.Value(contextKey)
	return ok != nil
}
