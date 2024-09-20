package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMetadata(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		paymentID          string
		body               string
		expectedPaymentID  models.PaymentID
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	paymentID := models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: "p1",
			Type:      models.PaymentTypePayIn,
		},
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	testCases := []testCase{
		{
			name:              "nominal",
			paymentID:         paymentID.String(),
			body:              "{\"foo\":\"bar\"}",
			expectedPaymentID: paymentID,
		},
		{
			name:               "missing body",
			paymentID:          paymentID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name:               "invalid body",
			paymentID:          paymentID.String(),
			body:               "invalid",
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name:               "invalid paymentID",
			paymentID:          "invalid",
			body:               "{\"foo\":\"bar\"}",
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "nominal",
			paymentID:          paymentID.String(),
			body:               "{\"foo\":\"bar\"}",
			expectedPaymentID:  paymentID,
			serviceError:       service.ErrValidation,
			expectedErrorCode:  ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "nominal",
			paymentID:          paymentID.String(),
			body:               "{\"foo\":\"bar\"}",
			expectedPaymentID:  paymentID,
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "nominal",
			paymentID:          paymentID.String(),
			body:               "{\"foo\":\"bar\"}",
			expectedPaymentID:  paymentID,
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "nominal",
			paymentID:          paymentID.String(),
			body:               "{\"foo\":\"bar\"}",
			expectedPaymentID:  paymentID,
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusNoContent
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					UpdatePaymentMetadata(gomock.Any(), testCase.expectedPaymentID, map[string]string{"foo": "bar"}).
					Return(nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					UpdatePaymentMetadata(gomock.Any(), testCase.expectedPaymentID, map[string]string{"foo": "bar"}).
					Return(testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{
				Debug: testing.Verbose(),
			}, auth.NewNoAuth(), false)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/payments/%s/metadata", testCase.paymentID), strings.NewReader(testCase.body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode >= 300 || testCase.expectedStatusCode < 200 {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
