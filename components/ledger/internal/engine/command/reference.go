package command

import (
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/pkg/errors"
)

var (
	ErrAlreadyTaken = errors.New("already taken")
)

type Reference int

const (
	referenceReverts = iota
	referenceIks
	referenceTxReference
)

type Referencer struct {
	references map[Reference]*sync.Map

	idempotencyCache *expirable.LRU[string, struct{}]
}

func (r *Referencer) take(ref Reference, key any) error {
	_, loaded := r.references[ref].LoadOrStore(fmt.Sprintf("%d/%s", ref, key), struct{}{})
	if loaded {
		return ErrAlreadyTaken
	}

	if ref == referenceIks {
		_, ok := r.idempotencyCache.Get(fmt.Sprintf("%d/%s", ref, key))
		if ok {
			return ErrAlreadyTaken
		}

		r.idempotencyCache.Add(fmt.Sprintf("%d/%s", ref, key), struct{}{})
	}
	return nil
}

func (r *Referencer) release(ref Reference, key any) {
	r.references[ref].Delete(fmt.Sprintf("%d/%s", ref, key))
}

func NewReferencer() *Referencer {
	return &Referencer{
		references: map[Reference]*sync.Map{
			referenceReverts:     {},
			referenceIks:         {},
			referenceTxReference: {},
		},
		idempotencyCache: expirable.NewLRU[string, struct{}](0, nil, 5*time.Minute),
	}
}
