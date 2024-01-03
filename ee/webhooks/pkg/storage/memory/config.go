package memory

import (
	"context"
	"reflect"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	//lint:ignore ST1001 shared definitions
	. "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage"
)

type configStore struct {
	mtx     sync.RWMutex
	configs map[string]Config
}

func NewConfigStore() ConfigStore {
	return &configStore{
		configs: make(map[string]Config),
	}
}

func (c *configStore) FindManyConfigs(ctx context.Context, filter map[string]any) ([]Config, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	result := make([]Config, 0, len(c.configs))
	for _, cfg := range c.configs {
		if configMatchesFilter(cfg, filter) {
			result = append(result, cfg)
		}
	}
	slices.SortFunc(result, func(a, b Config) int {
		return b.UpdatedAt.Compare(a.UpdatedAt)
	})
	return result, nil
}

func configMatchesFilter(config Config, filter map[string]any) bool {
FilterLoop:
	for key, val := range filter {
		switch key {
		case "id":
			if config.ID != val.(string) {
				return false
			}
		case "endpoint":
			if config.Endpoint != val.(string) {
				return false
			}
		case "active":
			if config.Active != val.(bool) {
				return false
			}
		case "event_types":
			for _, eventType := range config.EventTypes {
				if eventType == val.(string) {
					continue FilterLoop
				}
			}
			return false
		}
	}
	return true
}

func (c *configStore) InsertOneConfig(ctx context.Context, cfgUser ConfigUser) (Config, error) {
	cfg := NewConfig(cfgUser)
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.configs[cfg.ID] = cfg
	return cfg, nil
}

func (c *configStore) DeleteOneConfig(ctx context.Context, id string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if _, ok := c.configs[id]; !ok {
		return storage.ErrConfigNotFound
	}
	delete(c.configs, id)
	return nil
}

func (c *configStore) UpdateOneConfigActivation(ctx context.Context, id string, active bool) (Config, error) {
	return c.updateOne(id, func(cfg *Config) {
		cfg.Active = active
	})
}

func (c *configStore) UpdateOneConfigSecret(ctx context.Context, id string, secret string) (Config, error) {
	return c.updateOne(id, func(cfg *Config) {
		cfg.Secret = secret
	})
}

func (c *configStore) updateOne(id string, update func(*Config)) (Config, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	cfg, ok := c.configs[id]
	if !ok {
		return Config{}, storage.ErrConfigNotFound
	}
	newCfg := cfg
	update(&newCfg)
	if reflect.DeepEqual(newCfg, cfg) {
		return Config{}, storage.ErrConfigNotModified
	}
	newCfg.UpdatedAt = time.Now().UTC()
	c.configs[id] = newCfg
	return newCfg, nil
}
