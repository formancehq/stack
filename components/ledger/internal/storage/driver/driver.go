package driver

import (
	"context"
	"sync"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"

	"github.com/formancehq/ledger/internal/storage/ledgerstore"

	"github.com/formancehq/ledger/internal/storage/sqlutils"

	"github.com/formancehq/ledger/internal/storage/systemstore"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const defaultBucket = "_default"

var (
	ErrNeedUpgradeBucket = errors.New("need to upgrade bucket before add a new ledger on it")
	ErrAlreadyExists     = errors.New("ledger already exists")
)

type LedgerConfiguration struct {
	Bucket string `json:"bucket"`
}

type Driver struct {
	systemStore       *systemstore.Store
	lock              sync.Mutex
	connectionOptions sqlutils.ConnectionOptions
	buckets           map[string]*ledgerstore.Bucket
	db                *bun.DB
}

func (d *Driver) GetSystemStore() *systemstore.Store {
	return d.systemStore
}

func (d *Driver) OpenBucket(name string) (*ledgerstore.Bucket, error) {

	bucket, ok := d.buckets[name]
	if ok {
		return bucket, nil
	}

	b, err := ledgerstore.ConnectToBucket(d.connectionOptions, name)
	if err != nil {
		return nil, err
	}
	d.buckets[name] = b

	return b, nil
}

func (d *Driver) GetLedgerStore(ctx context.Context, name string) (*ledgerstore.Store, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	ledgerConfiguration, err := d.systemStore.GetLedger(ctx, name)
	if err != nil {
		return nil, err
	}

	bucket, err := d.OpenBucket(ledgerConfiguration.Bucket)
	if err != nil {
		return nil, err
	}

	return bucket.GetLedgerStore(name)
}

func (f *Driver) CreateLedgerStore(ctx context.Context, name string, configuration LedgerConfiguration) (*ledgerstore.Store, error) {
	bucketName := defaultBucket
	if configuration.Bucket != "" {
		bucketName = configuration.Bucket
	}

	bucket, err := f.OpenBucket(bucketName)
	if err != nil {
		return nil, err
	}

	isInitialized, err := bucket.IsInitialized(ctx)
	if err != nil {
		return nil, err
	}

	if isInitialized {
		isUpToDate, err := bucket.IsUpToDate(ctx)
		if err != nil {
			return nil, err
		}
		if !isUpToDate {
			return nil, ErrNeedUpgradeBucket
		}
	} else {
		if err := bucket.Migrate(ctx); err != nil {
			return nil, err
		}
	}

	ledgerExists, err := bucket.HasLedger(ctx, name)
	if err != nil {
		return nil, err
	}
	if ledgerExists {
		return nil, ErrAlreadyExists
	}

	store, err := bucket.CreateLedgerStore(ctx, name)
	if err != nil {
		return nil, err
	}

	_, err = f.systemStore.RegisterLedger(ctx, &systemstore.Ledger{
		Name:    name,
		AddedAt: ledger.Now(),
		Bucket:  bucketName,
	})
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (d *Driver) Initialize(ctx context.Context) error {
	logging.FromContext(ctx).Debugf("Initialize driver")

	var err error
	d.db, err = sqlutils.OpenSQLDB(d.connectionOptions)
	if err != nil {
		return errors.Wrap(err, "connecting to database")
	}

	d.systemStore, err = systemstore.Connect(ctx, d.connectionOptions)
	if err != nil {
		return errors.Wrap(err, "connecting to system store")
	}

	if err := d.systemStore.Migrate(ctx); err != nil {
		return errors.Wrap(err, "migrating data")
	}

	return nil
}

func (d *Driver) UpgradeAllBuckets(ctx context.Context) error {
	systemStore := d.GetSystemStore()
	ledgers, err := systemStore.ListLedgers(ctx)
	if err != nil {
		return err
	}

	buckets := collectionutils.Set[string]{}
	for _, name := range ledgers {
		buckets.Put(name.Name)
	}

	for _, bucket := range collectionutils.Keys(buckets) {
		bucket, err := d.OpenBucket(bucket)
		if err != nil {
			return err
		}

		logging.FromContext(ctx).Infof("Upgrading bucket '%s'", bucket)
		if err := bucket.Migrate(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (d *Driver) Close() error {
	if err := d.systemStore.Close(); err != nil {
		return err
	}
	for _, b := range d.buckets {
		if err := b.Close(); err != nil {
			return err
		}
	}
	if err := d.db.Close(); err != nil {
		return err
	}
	return nil
}

func New(connectionOptions sqlutils.ConnectionOptions) *Driver {
	return &Driver{
		connectionOptions: connectionOptions,
		buckets:           make(map[string]*ledgerstore.Bucket),
	}
}
