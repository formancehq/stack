package api

import (
	"testing"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/golang/mock/gomock"
)

func newTestingBackend(t *testing.T) (*backend.MockBackend, *backend.MockService) {
	ctrl := gomock.NewController(t)
	mockService := backend.NewMockService(ctrl)
	backend := backend.NewMockBackend(ctrl)
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
