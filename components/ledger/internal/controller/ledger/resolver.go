package ledger

import (
	"context"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/pkg/errors"
)

type Option func(r *Resolver)

func WithCompiler(compiler Compiler) Option {
	return func(r *Resolver) {
		r.compiler = compiler
	}
}

var defaultOptions = []Option{
	WithCompiler(NewDefaultCompiler()),
}

type Resolver struct {
	storageDriver StorageDriver
	compiler      Compiler
	listener      Listener
	registry      *StateRegistry
}

func NewResolver(storageDriver StorageDriver, listener Listener, options ...Option) *Resolver {
	r := &Resolver{
		storageDriver: storageDriver,
		listener:      listener,
		registry:      NewStateRegistry(),
	}
	for _, opt := range append(defaultOptions, options...) {
		opt(r)
	}

	return r
}

func (r *Resolver) GetLedger(ctx context.Context, name string) (Controller, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}

	store, l, err := r.storageDriver.OpenLedger(ctx, name)
	if err != nil {
		return nil, err
	}

	return NewControllerWithCache(
		*l,
		NewDefaultController(*l, store, r.listener, NewDefaultMachineFactory(r.compiler)),
		r.registry,
	), nil
}

func (r *Resolver) CreateLedger(ctx context.Context, name string, configuration ledger.Configuration) error {
	if name == "" {
		return errors.New("empty name")
	}

	return r.storageDriver.CreateLedger(ctx, name, configuration)
}
