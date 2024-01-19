package core

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Context interface {
	context.Context
	GetClient() client.Client
	GetScheme() *runtime.Scheme
	GetPlatform() Platform
}

type defaultContext struct {
	context.Context
	client   client.Client
	scheme   *runtime.Scheme
	platform Platform
}

func (d defaultContext) GetPlatform() Platform {
	return d.platform
}

func (d defaultContext) GetClient() client.Client {
	return d.client
}

func (d defaultContext) GetScheme() *runtime.Scheme {
	return d.scheme
}

var _ Context = &defaultContext{}

func NewContext(client client.Client, scheme *runtime.Scheme, platform Platform, ctx context.Context) *defaultContext {
	return &defaultContext{
		Context:  ctx,
		client:   client,
		scheme:   scheme,
		platform: platform,
	}
}
