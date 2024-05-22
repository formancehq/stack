package models

import (
	"fmt"
	"time"

	"github.com/formancehq/stack/libs/go-libs/sync/shared"
	"github.com/google/uuid"
)

// #############################################################################
// #############################################################################

type HookStatus string

const (
	EnableStatus  HookStatus = "ENABLED"
	DisableStatus HookStatus = "DISABLED"
	DeleteStatus  HookStatus = "DELETED"
)

type Hook struct {
	ID           string     `json:"id" bun:",pk"`
	Name         string     `json:"name" bun:"name"`
	Status       HookStatus `json:"status" bun:"status"`
	Events       []string   `json:"events" bun:"event_types,array"`
	Endpoint     string     `json:"endpoint" bun:"endpoint"`
	Secret       string     `json:"secret" bun:"secret"`
	Retry        bool       `json:"retry" bun:"retry"`
	DateCreation time.Time  `json:"dateCreation" bun:"created_at"`
	DateStatus   time.Time  `json:"dateStatus" bun:"date_status"`
	Active       bool       `json:"active"`                                                                //v1
	UpdatedAt    time.Time  `json:"updatedAt" bun:"updated_at,nullzero,notnull,default:current_timestamp"` //v1
}

func (h *Hook) IsActive() bool {
	return h.Status == EnableStatus
}

func (h *Hook) Delete() {
	h.Status = DeleteStatus
	h.DateStatus = time.Now()
}

func (h *Hook) Enable() {
	h.Status = EnableStatus
	h.DateStatus = time.Now()
}

func (h *Hook) Disable() {
	h.Status = DisableStatus
	h.DateStatus = time.Now()
}

func NewHook(name string, events []string, endpoint string, secret string, retry bool) *Hook {

	pId := uuid.NewString()
	pName := name
	if pName == "" {
		pName = fmt.Sprintf("Hook-%s", pId[:5])
	}

	return &Hook{
		ID:           pId,
		Name:         pName,
		Status:       DisableStatus,
		Events:       events,
		Endpoint:     endpoint,
		Secret:       secret,
		Retry:        retry,
		DateCreation: time.Now(),
		DateStatus:   time.Now(),
	}

}

// #############################################################################
// #############################################################################
type SharedHook = shared.Shared[Hook]

func NewSharedHook(name string, events []string, endpoint string, secret string, retry bool) *SharedHook {
	s := shared.NewShared(NewHook(name, events, endpoint, secret, retry))
	return &s
}

// #############################################################################
// #############################################################################

type SharedHooks = shared.SharedArr[Hook]

func NewSharedHooks() *SharedHooks {
	s := shared.NewSharedArr[Hook]()
	return &s
}

// #############################################################################
// #############################################################################

type MapSharedHooks = shared.SharedMapArr[Hook]

func NewMapSharedHooks() *MapSharedHooks {
	m := shared.NewSharedMapArr[Hook]()
	return &m
}

// #############################################################################
// #############################################################################

type MapSharedHook = shared.SharedMap[Hook]

func NewMapSharedHook() *MapSharedHook {
	m := shared.NewSharedMap[Hook]()
	return &m
}

// #############################################################################
// #############################################################################
// For Controller

type HookBodyParams struct {
	Name     string   `json:"name"`
	Endpoint string   `json:"endpoint"`
	Secret   string   `json:"secret"`
	Events   []string `json:"events"`
	Retry    bool     `json:"retry"`
}
