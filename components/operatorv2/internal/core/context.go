package core

import (
	"context"
)

type Context interface {
	context.Context
	Manager
}

type defaultContext struct {
	context.Context
	Manager
}

func NewContext(mgr Manager, ctx context.Context) *defaultContext {
	return &defaultContext{
		Context: ctx,
		Manager: mgr,
	}
}
