package engine

import (
	"context"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/ledger/internal/engine/command"
	"github.com/formancehq/ledger/internal/opentelemetry/metrics"
	"github.com/formancehq/ledger/internal/storage/driver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/sirupsen/logrus"
)

type option func(r *Resolver)

func WithMessagePublisher(publisher message.Publisher) option {
	return func(r *Resolver) {
		r.publisher = publisher
	}
}

func WithMetricsRegistry(registry metrics.GlobalRegistry) option {
	return func(r *Resolver) {
		r.metricsRegistry = registry
	}
}

func WithCompiler(compiler *command.Compiler) option {
	return func(r *Resolver) {
		r.compiler = compiler
	}
}

func WithLogger(logger logging.Logger) option {
	return func(r *Resolver) {
		r.logger = logger
	}
}

var defaultOptions = []option{
	WithMetricsRegistry(metrics.NewNoOpRegistry()),
	WithCompiler(command.NewCompiler(1024)),
	WithLogger(logging.NewLogrus(logrus.New())),
}

type Resolver struct {
	storageDriver   *driver.Driver
	lock            sync.RWMutex
	metricsRegistry metrics.GlobalRegistry
	//TODO(gfyrag): add a routine to clean old ledger
	buckets   map[string]*Ledger
	compiler  *command.Compiler
	logger    logging.Logger
	publisher message.Publisher
}

func NewResolver(storageDriver *driver.Driver, options ...option) *Resolver {
	r := &Resolver{
		storageDriver: storageDriver,
		buckets:       map[string]*Ledger{},
	}
	for _, opt := range append(defaultOptions, options...) {
		opt(r)
	}

	return r
}

func (r *Resolver) GetLedger(ctx context.Context, name string) (*Ledger, error) {
	r.lock.RLock()
	ledger, ok := r.buckets[name]
	r.lock.RUnlock()

	if !ok {
		r.lock.Lock()
		defer r.lock.Unlock()

		logging.FromContext(ctx).Infof("Initialize new ledger")

		ledger, ok = r.buckets[name]
		if ok {
			return ledger, nil
		}

		bucket, err := r.storageDriver.GetBucket(ctx, name)
		if err != nil {
			return nil, err
		}

		store, err := bucket.GetLedgerStore(ctx, name)
		if err != nil {
			return nil, err
		}

		ledger = New(bucket, store, r.publisher, r.compiler)
		ledger.Start(logging.ContextWithLogger(context.Background(), r.logger))
		r.buckets[name] = ledger
		r.metricsRegistry.ActiveLedgers().Add(ctx, +1)
	}

	return ledger, nil
}

func (r *Resolver) CloseLedgers(ctx context.Context) error {
	r.logger.Info("Close all ledgers")
	defer func() {
		r.logger.Info("All ledgers closed")
	}()
	for name, ledger := range r.buckets {
		r.logger.Infof("Close ledger %s", name)
		ledger.Close(logging.ContextWithLogger(ctx, r.logger.WithField("ledger", name)))
		delete(r.buckets, name)
	}

	return nil
}
