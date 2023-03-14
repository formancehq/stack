package cache

import (
	"github.com/bluele/gcache"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/storage"
)

type TxInfo struct {
	Date core.Time
	ID   uint64
}

// TODO: Add a mutex for concurrent ledger access
type Cache struct {
	cache                   gcache.Cache
	lastInsertedTransaction map[string]*TxInfo
	driver                  storage.Driver
}

func (c *Cache) ForLedger(ledger string) *Ledger {
	return newLedgerCache(c, ledger)
}

func NewCache(driver storage.Driver) *Cache {
	return &Cache{
		driver:                  driver,
		cache:                   gcache.New(1000).LRU().Build(), // TODO: Make configurable
		lastInsertedTransaction: map[string]*TxInfo{},
	}
}
