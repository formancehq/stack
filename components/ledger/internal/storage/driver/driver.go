package driver

import (
	"context"
	"fmt"
	"sync"

	"github.com/formancehq/ledger/internal/storage/ledgerstore"

	"github.com/formancehq/ledger/internal/storage/sqlutils"

	"github.com/formancehq/ledger/internal/storage/systemstore"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
)

type Driver struct {
	systemStore       *systemstore.Store
	lock              sync.Mutex
	connectionOptions sqlutils.ConnectionOptions
	buckets           map[string]*ledgerstore.Bucket
}

func (d *Driver) GetSystemStore() *systemstore.Store {
	return d.systemStore
}

func (d *Driver) openBucket(name string) (*ledgerstore.Bucket, error) {

	b, err := ledgerstore.ConnectToBucket(d.systemStore, d.connectionOptions, name)
	if err != nil {
		return nil, err
	}
	d.buckets[name] = b

	return b, nil
}

func (d *Driver) createBucket(ctx context.Context, name string) (*ledgerstore.Bucket, error) {
	if name == systemstore.Schema {
		return nil, errors.New("reserved name")
	}

	exists, err := d.systemStore.ExistsBucket(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, sqlutils.ErrBucketAlreadyExists
	}

	_, err = d.systemStore.RegisterBucket(ctx, name)
	if err != nil {
		return nil, err
	}

	bucket, err := d.openBucket(name)
	if err != nil {
		return nil, err
	}

	_, err = bucket.DB().ExecContext(ctx, fmt.Sprintf(`create schema if not exists "%s"`, name))
	if err != nil {
		return nil, sqlutils.PostgresError(err)
	}

	err = bucket.Migrate(ctx)
	if err != nil {
		return nil, err
	}

	return bucket, err
}

func (d *Driver) GetBucket(ctx context.Context, name string) (*ledgerstore.Bucket, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	exists, err := d.systemStore.ExistsBucket(ctx, name)
	if err != nil {
		return nil, err
	}

	var ret *ledgerstore.Bucket
	if !exists {
		ret, err = d.createBucket(ctx, name)
	} else {
		ret, err = d.openBucket(name)
	}
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (d *Driver) GetLedgerStore(ctx context.Context, name string) (*ledgerstore.Store, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	const defaultBucketName = "_default"

	exists, err := d.systemStore.ExistsBucket(ctx, defaultBucketName)
	if err != nil {
		return nil, err
	}

	var ret *ledgerstore.Bucket
	if !exists {
		ret, err = d.createBucket(ctx, name)
	} else {
		ret, err = d.openBucket(name)
	}
	if err != nil {
		return nil, err
	}

	return ret.GetLedgerStore(ctx, name)
}

func (d *Driver) Initialize(ctx context.Context) error {
	logging.FromContext(ctx).Debugf("Initialize driver")

	var err error
	d.systemStore, err = systemstore.Connect(ctx, d.connectionOptions)
	if err != nil {
		return err
	}

	if err := d.systemStore.Migrate(ctx); err != nil {
		return err
	}

	return nil
}

func (d *Driver) UpgradeAllBuckets(ctx context.Context) error {
	systemStore := d.GetSystemStore()
	buckets, err := systemStore.ListBuckets(ctx)
	if err != nil {
		return err
	}

	for _, name := range buckets {
		bucket, err := d.GetBucket(ctx, name)
		if err != nil {
			return err
		}

		logging.FromContext(ctx).Infof("Upgrading storage '%s'", name)
		if err := bucket.Migrate(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (d *Driver) Close() error {
	if d.systemStore != nil {
		if err := d.systemStore.Close(); err != nil {
			return err
		}
	}
	for _, b := range d.buckets {
		if err := b.Close(); err != nil {
			return err
		}
	}
	return nil
}

func New(connectionOptions sqlutils.ConnectionOptions) *Driver {
	return &Driver{
		connectionOptions: connectionOptions,
		buckets:           make(map[string]*ledgerstore.Bucket),
	}
}
