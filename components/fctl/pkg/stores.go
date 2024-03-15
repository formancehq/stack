package fctl

import (
	"context"

	"github.com/formancehq/fctl/membershipclient"
	v2 "github.com/formancehq/formance-sdk-go/v2"
)

var storeKey string = "_stores"
var stackKey string = "_stack"

func ContextWithStore(ctx context.Context, key string, store interface{}) context.Context {
	var stores map[string]interface{}
	stores, ok := ctx.Value(storeKey).(map[string]interface{})
	if !ok {
		stores = map[string]interface{}{}
	}
	stores[key] = store

	return context.WithValue(ctx, storeKey, stores)
}

func GetStore(ctx context.Context, key string) any {
	stores, ok := ctx.Value(storeKey).(map[string]interface{})
	if !ok {
		return nil
	}
	store, ok := stores[key]
	if !ok {
		return nil
	}
	return store
}

type StackStore struct {
	Config         *Config
	stack          *membershipclient.Stack
	stackClient    *v2.Formance
	organizationId string
}

func (cns StackStore) Client() *v2.Formance {
	return cns.stackClient
}

func (cns StackStore) Stack() *membershipclient.Stack {
	return cns.stack
}

func (cns StackStore) OrganizationId() string {
	return cns.organizationId
}

func StackNode(config *Config, stack *membershipclient.Stack, organization string, stackClient *v2.Formance) *StackStore {
	return &StackStore{
		Config:         config,
		stack:          stack,
		organizationId: organization,
		stackClient:    stackClient,
	}
}

func GetStackStore(ctx context.Context) *StackStore {
	return GetStore(ctx, stackKey).(*StackStore)
}

func ContextWithStackStore(ctx context.Context, store *StackStore) context.Context {
	return ContextWithStore(ctx, stackKey, store)
}
