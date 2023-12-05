package api

import (
	"testing"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/dummypay"
	gomock "github.com/golang/mock/gomock"
)

func newServiceTestingBackend(t *testing.T) (*backend.MockServiceBackend, *backend.MockService) {
	ctrl := gomock.NewController(t)
	mockService := backend.NewMockService(ctrl)
	backend := backend.NewMockServiceBackend(ctrl)
	backend.
		EXPECT().
		GetService().
		MinTimes(0).
		Return(mockService)
	t.Cleanup(func() {
		ctrl.Finish()
	})
	return backend, mockService
}

func newConnectorManagerTestingBackend(t *testing.T) (*backend.MockManagerBackend[dummypay.Config], *backend.MockManager[dummypay.Config]) {
	ctrl := gomock.NewController(t)
	mockManager := backend.NewMockManager[dummypay.Config](ctrl)
	backend := backend.NewMockManagerBackend[dummypay.Config](ctrl)
	backend.
		EXPECT().
		GetManager().
		MinTimes(0).
		Return(mockManager)
	t.Cleanup(func() {
		ctrl.Finish()
	})
	return backend, mockManager
}

func ptr[T any](v T) *T {
	return &v
}
