package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReadConnectors(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	testCases := []testCase{
		{
			name: "nominal",
		},
		{
			name:               "service error duplicate key",
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "service error not found",
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "service error validation",
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "service error invalid ID",
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "service error publish",
			serviceError:       service.ErrPublish,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
		{
			name:               "service error unknown",
			serviceError:       errors.New("unknown"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			listConnectorsResponse := []*models.Connector{
				{
					ID: models.ConnectorID{
						Reference: uuid.New(),
						Provider:  models.ConnectorProviderDummyPay,
					},
					Name:      "c1",
					CreatedAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Provider:  models.ConnectorProviderDummyPay,
				},
			}

			expectedListConnectorsResponse := []readConnectorsResponseElement{
				{
					Provider:    listConnectorsResponse[0].Provider,
					ConnectorID: listConnectorsResponse[0].ID.String(),
					Name:        "c1",
					Enabled:     true,
				},
			}

			backend, mockService := newServiceTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					ListConnectors(gomock.Any()).
					Return(listConnectorsResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ListConnectors(gomock.Any()).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), nil, false)

			req := httptest.NewRequest(http.MethodGet, "/connectors", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[[]readConnectorsResponseElement]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, &expectedListConnectorsResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
