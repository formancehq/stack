package fctl

import "context"

var storeKey string = "_stores"

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
