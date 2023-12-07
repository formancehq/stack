package backend

import (
	"context"

	"github.com/formancehq/reconciliation/internal/api/service"
)

//go:generate mockgen -source backend.go -destination backend_generated.go -package backend . Service
type Service interface {
	Reconciliation(ctx context.Context, req *service.ReconciliationRequest) (*service.ReconciliationResponse, error)
}

type Backend interface {
	GetService() Service
}

type DefaultBackend struct {
	service Service
}

func (d DefaultBackend) GetService() Service {
	return d.service
}

func NewDefaultBackend(service Service) Backend {
	return &DefaultBackend{
		service: service,
	}
}
