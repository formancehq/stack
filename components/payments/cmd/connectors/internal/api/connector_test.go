package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/dummypay"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReadConfig(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		connectorID        string
		connectors         map[string]*manager.ConnectorManager
		installed          *bool
		apiVersion         APIVersion
		expectedStatusCode int
		expectedErrorCode  string
		managerError       error
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	testCases := []testCase{
		// V0 tests
		{
			name:       "nominal V0",
			apiVersion: V0,
			connectors: map[string]*manager.ConnectorManager{
				connectorID.String(): nil,
			},
			installed: ptr(true),
		},
		{
			name: "too many connectors for V0",
			connectors: map[string]*manager.ConnectorManager{
				connectorID.String(): nil,
				"1":                  nil,
			},
			apiVersion:         V0,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "no connector for V0",
			connectors:         map[string]*manager.ConnectorManager{},
			apiVersion:         V0,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		// V1 tests
		{
			name:        "nominal V1",
			connectorID: connectorID.String(),
			apiVersion:  V1,
			installed:   ptr(true),
		},
		// Common test for V0 and V1
		{
			name:               "connector not installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(false),
		},
		{
			name:               "manager error duplicate key value storage",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
			installed:          ptr(true),
		},
		{
			name:               "manager error err not found storage",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			installed:          ptr(true),
		},
		{
			name:               "manager error already installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrAlreadyInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error not installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrNotInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error connector not found",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrConnectorNotFound,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error err not found",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			installed:          ptr(true),
		},
		{
			name:               "manager error err validation",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error other errors",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
			installed:          ptr(true),
		},
	}

	for _, tc := range testCases {
		testCase := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			readConfigResponse := dummypay.Config{
				Name:      "test",
				Directory: "test",
				FilePollingPeriod: connectors.Duration{
					Duration: 2 * time.Minute,
				},
			}

			backend, _ := newServiceTestingBackend(t)
			managerBackend, mockManager := newConnectorManagerTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockManager.EXPECT().
					ReadConfig(gomock.Any(), connectorID).
					Return(readConfigResponse, nil)
			}

			if testCase.managerError != nil {
				mockManager.EXPECT().
					ReadConfig(gomock.Any(), connectorID).
					Return(dummypay.Config{}, testCase.managerError)
			}

			if testCase.apiVersion == V0 {
				mockManager.EXPECT().
					Connectors().
					Return(testCase.connectors)
			}

			if testCase.installed != nil {
				mockManager.EXPECT().
					IsInstalled(gomock.Any(), connectorID).
					Return(*testCase.installed, nil)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), []connectorHandler{
				{
					Handler:  connectorRouter[dummypay.Config](models.ConnectorProviderDummyPay, managerBackend),
					Provider: models.ConnectorProviderDummyPay,
					initiatePayment: func(ctx context.Context, transfer *models.TransferInitiation) error {
						return nil
					},
				},
			})

			var endpoint string
			switch testCase.apiVersion {
			case V0:
				endpoint = "/connectors/dummy-pay/config"
			case V1:
				endpoint = fmt.Sprintf("/connectors/dummy-pay/%s/config", testCase.connectorID)
			}
			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[dummypay.Config]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, &readConfigResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestListTasks(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		connectorID        string
		connectors         map[string]*manager.ConnectorManager
		installed          *bool
		apiVersion         APIVersion
		queryParams        url.Values
		pageSize           int
		expectedQuery      storage.ListTasksQuery
		expectedStatusCode int
		expectedErrorCode  string
		managerError       error
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	testCases := []testCase{
		// V0 tests
		{
			name:       "nominal V0",
			apiVersion: V0,
			connectors: map[string]*manager.ConnectorManager{
				connectorID.String(): nil,
			},
			expectedQuery: storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
			connectorID:   connectorID.String(),
			installed:     ptr(true),
			queryParams:   map[string][]string{},
			pageSize:      15,
		},
		{
			name: "too many connectors for V0",
			connectors: map[string]*manager.ConnectorManager{
				connectorID.String(): nil,
				"1":                  nil,
			},
			apiVersion:         V0,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "no connector for V0",
			connectors:         map[string]*manager.ConnectorManager{},
			apiVersion:         V0,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		// V1 tests
		{
			name:          "nominal V1",
			expectedQuery: storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
			connectorID:   connectorID.String(),
			installed:     ptr(true),
			apiVersion:    V1,
			pageSize:      15,
		},
		// Common test for V0 and V1
		{
			name: "page size too low",
			queryParams: url.Values{
				"pageSize": {"0"},
			},
			expectedQuery: storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
			pageSize:      15,
			connectorID:   connectorID.String(),
			installed:     ptr(true),
			apiVersion:    V1,
		},
		{
			name: "page size too high",
			queryParams: url.Values{
				"pageSize": {"100000"},
			},
			expectedQuery: storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(100)),
			pageSize:      100,
			connectorID:   connectorID.String(),
			installed:     ptr(true),
			apiVersion:    V1,
		},
		{
			name: "with invalid page size",
			queryParams: url.Values{
				"pageSize": {"nan"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			connectorID:        connectorID.String(),
			installed:          ptr(true),
			apiVersion:         V1,
		},
		{
			name:               "connector not installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(false),
		},
		{
			name:               "manager error duplicate key value storage",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
			installed:          ptr(true),
			expectedQuery:      storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
		},
		{
			name:               "manager error err not found storage",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			installed:          ptr(true),
			expectedQuery:      storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
		},
		{
			name:               "manager error already installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrAlreadyInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
			expectedQuery:      storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
		},
		{
			name:               "manager error not installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrNotInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
			expectedQuery:      storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
		},
		{
			name:               "manager error connector not found",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrConnectorNotFound,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
			expectedQuery:      storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
		},
		{
			name:               "manager error err not found",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			installed:          ptr(true),
			expectedQuery:      storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
		},
		{
			name:               "manager error err validation",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
			expectedQuery:      storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
		},
		{
			name:               "manager error other errors",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
			installed:          ptr(true),
			expectedQuery:      storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(15)),
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			tasks := []models.Task{
				{
					ID:          uuid.New(),
					ConnectorID: connectorID,
					CreatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Name:        "test",
					Descriptor:  []byte("{}"),
					Status:      models.TaskStatusActive,
					State:       json.RawMessage("{}"),
				},
			}
			listTasksResponse := &api.Cursor[models.Task]{
				PageSize: testCase.pageSize,
				HasMore:  false,
				Previous: "",
				Next:     "",
				Data:     tasks,
			}

			expectedListTasksResponse := []listTasksResponseElement{
				{
					ID:          tasks[0].ID.String(),
					ConnectorID: tasks[0].ConnectorID.String(),
					CreatedAt:   tasks[0].CreatedAt.Format(time.RFC3339),
					UpdatedAt:   tasks[0].UpdatedAt.Format(time.RFC3339),
					Descriptor:  tasks[0].Descriptor,
					Status:      tasks[0].Status,
					State:       tasks[0].State,
					Error:       tasks[0].Error,
				},
			}

			backend, _ := newServiceTestingBackend(t)
			managerBackend, mockManager := newConnectorManagerTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockManager.EXPECT().
					ListTasksStates(gomock.Any(), connectorID, testCase.expectedQuery).
					Return(listTasksResponse, nil)
			}

			if testCase.managerError != nil {
				mockManager.EXPECT().
					ListTasksStates(gomock.Any(), connectorID, testCase.expectedQuery).
					Return(nil, testCase.managerError)
			}

			if testCase.apiVersion == V0 {
				mockManager.EXPECT().
					Connectors().
					Return(testCase.connectors)
			}

			if testCase.installed != nil {
				mockManager.EXPECT().
					IsInstalled(gomock.Any(), connectorID).
					Return(*testCase.installed, nil)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), []connectorHandler{
				{
					Handler:  connectorRouter[dummypay.Config](models.ConnectorProviderDummyPay, managerBackend),
					Provider: models.ConnectorProviderDummyPay,
					initiatePayment: func(ctx context.Context, transfer *models.TransferInitiation) error {
						return nil
					},
				},
			})

			var endpoint string
			switch testCase.apiVersion {
			case V0:
				endpoint = "/connectors/dummy-pay/tasks"
			case V1:
				endpoint = fmt.Sprintf("/connectors/dummy-pay/%s/tasks", testCase.connectorID)
			}
			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[listTasksResponseElement]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedListTasksResponse, resp.Cursor.Data)
				require.Equal(t, listTasksResponse.PageSize, resp.Cursor.PageSize)
				require.Equal(t, listTasksResponse.HasMore, resp.Cursor.HasMore)
				require.Equal(t, listTasksResponse.Next, resp.Cursor.Next)
				require.Equal(t, listTasksResponse.Previous, resp.Cursor.Previous)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestReadTask(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		connectorID        string
		taskID             string
		connectors         map[string]*manager.ConnectorManager
		installed          *bool
		apiVersion         APIVersion
		expectedStatusCode int
		expectedErrorCode  string
		managerError       error
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	taskID := uuid.New()

	testCases := []testCase{
		// V0 tests
		{
			name:       "nominal V0",
			apiVersion: V0,
			connectors: map[string]*manager.ConnectorManager{
				connectorID.String(): nil,
			},
			taskID:    taskID.String(),
			installed: ptr(true),
		},
		{
			name: "too many connectors for V0",
			connectors: map[string]*manager.ConnectorManager{
				connectorID.String(): nil,
				"1":                  nil,
			},
			apiVersion:         V0,
			taskID:             taskID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "no connector for V0",
			connectors:         map[string]*manager.ConnectorManager{},
			apiVersion:         V0,
			taskID:             taskID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		// V1 tests
		{
			name:        "nominal V1",
			connectorID: connectorID.String(),
			taskID:      taskID.String(),
			apiVersion:  V1,
			installed:   ptr(true),
		},
		// Common test for V0 and V1
		{
			name:               "connector not installed",
			connectorID:        connectorID.String(),
			taskID:             taskID.String(),
			apiVersion:         V1,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(false),
		},
		{
			name:               "manager error duplicate key value storage",
			connectorID:        connectorID.String(),
			taskID:             taskID.String(),
			apiVersion:         V1,
			managerError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
			installed:          ptr(true),
		},
		{
			name:               "manager error err not found storage",
			connectorID:        connectorID.String(),
			taskID:             taskID.String(),
			apiVersion:         V1,
			managerError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			installed:          ptr(true),
		},
		{
			name:               "manager error already installed",
			connectorID:        connectorID.String(),
			taskID:             taskID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrAlreadyInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error not installed",
			connectorID:        connectorID.String(),
			taskID:             taskID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrNotInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error connector not found",
			connectorID:        connectorID.String(),
			taskID:             taskID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrConnectorNotFound,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error err not found",
			connectorID:        connectorID.String(),
			taskID:             taskID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			installed:          ptr(true),
		},
		{
			name:               "manager error err validation",
			connectorID:        connectorID.String(),
			taskID:             taskID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error other errors",
			connectorID:        connectorID.String(),
			taskID:             taskID.String(),
			apiVersion:         V1,
			managerError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
			installed:          ptr(true),
		},
		{
			name:               "invalid task ID",
			apiVersion:         V1,
			connectorID:        connectorID.String(),
			taskID:             "invalid",
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
			installed:          ptr(true),
		},
	}

	for _, tc := range testCases {
		testCase := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			readTaskResponse := &models.Task{
				ID:          uuid.New(),
				ConnectorID: connectorID,
				CreatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Name:        "test",
				Descriptor:  []byte("{}"),
				Status:      models.TaskStatusActive,
				State:       json.RawMessage("{}"),
			}

			expectedReadTasksResponse := listTasksResponseElement{
				ID:          readTaskResponse.ID.String(),
				ConnectorID: readTaskResponse.ConnectorID.String(),
				CreatedAt:   readTaskResponse.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   readTaskResponse.UpdatedAt.Format(time.RFC3339),
				Descriptor:  readTaskResponse.Descriptor,
				Status:      readTaskResponse.Status,
				State:       readTaskResponse.State,
				Error:       readTaskResponse.Error,
			}

			backend, _ := newServiceTestingBackend(t)
			managerBackend, mockManager := newConnectorManagerTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockManager.EXPECT().
					ReadTaskState(gomock.Any(), connectorID, taskID).
					Return(readTaskResponse, nil)
			}

			if testCase.managerError != nil {
				mockManager.EXPECT().
					ReadTaskState(gomock.Any(), connectorID, taskID).
					Return(nil, testCase.managerError)
			}

			if testCase.apiVersion == V0 {
				mockManager.EXPECT().
					Connectors().
					Return(testCase.connectors)
			}

			if testCase.installed != nil {
				mockManager.EXPECT().
					IsInstalled(gomock.Any(), connectorID).
					Return(*testCase.installed, nil)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), []connectorHandler{
				{
					Handler:  connectorRouter[dummypay.Config](models.ConnectorProviderDummyPay, managerBackend),
					Provider: models.ConnectorProviderDummyPay,
					initiatePayment: func(ctx context.Context, transfer *models.TransferInitiation) error {
						return nil
					},
				},
			})

			var endpoint string
			switch testCase.apiVersion {
			case V0:
				endpoint = fmt.Sprintf("/connectors/dummy-pay/tasks/%s", testCase.taskID)
			case V1:
				endpoint = fmt.Sprintf("/connectors/dummy-pay/%s/tasks/%s", testCase.connectorID, testCase.taskID)
			}
			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[listTasksResponseElement]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, &expectedReadTasksResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestUninstall(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		connectorID        string
		connectors         map[string]*manager.ConnectorManager
		installed          *bool
		apiVersion         APIVersion
		expectedStatusCode int
		expectedErrorCode  string
		managerError       error
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	testCases := []testCase{
		// V0 tests
		{
			name:       "nominal V0",
			apiVersion: V0,
			connectors: map[string]*manager.ConnectorManager{
				connectorID.String(): nil,
			},
			installed: ptr(true),
		},
		{
			name: "too many connectors for V0",
			connectors: map[string]*manager.ConnectorManager{
				connectorID.String(): nil,
				"1":                  nil,
			},
			apiVersion:         V0,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "no connector for V0",
			connectors:         map[string]*manager.ConnectorManager{},
			apiVersion:         V0,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		// V1 tests
		{
			name:        "nominal V1",
			connectorID: connectorID.String(),
			apiVersion:  V1,
			installed:   ptr(true),
		},
		// Common test for V0 and V1
		{
			name:               "connector not installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(false),
		},
		{
			name:               "manager error duplicate key value storage",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
			installed:          ptr(true),
		},
		{
			name:               "manager error err not found storage",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			installed:          ptr(true),
		},
		{
			name:               "manager error already installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrAlreadyInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error not installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrNotInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error connector not found",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrConnectorNotFound,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error err not found",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			installed:          ptr(true),
		},
		{
			name:               "manager error err validation",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error other errors",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
			installed:          ptr(true),
		},
	}

	for _, tc := range testCases {
		testCase := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusNoContent
			}

			backend, _ := newServiceTestingBackend(t)
			managerBackend, mockManager := newConnectorManagerTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockManager.EXPECT().
					Uninstall(gomock.Any(), connectorID).
					Return(nil)
			}

			if testCase.managerError != nil {
				mockManager.EXPECT().
					Uninstall(gomock.Any(), connectorID).
					Return(testCase.managerError)
			}

			if testCase.apiVersion == V0 {
				mockManager.EXPECT().
					Connectors().
					Return(testCase.connectors)
			}

			if testCase.installed != nil {
				mockManager.EXPECT().
					IsInstalled(gomock.Any(), connectorID).
					Return(*testCase.installed, nil)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), []connectorHandler{
				{
					Handler:  connectorRouter[dummypay.Config](models.ConnectorProviderDummyPay, managerBackend),
					Provider: models.ConnectorProviderDummyPay,
					initiatePayment: func(ctx context.Context, transfer *models.TransferInitiation) error {
						return nil
					},
				},
			})

			var endpoint string
			switch testCase.apiVersion {
			case V0:
				endpoint = "/connectors/dummy-pay"
			case V1:
				endpoint = fmt.Sprintf("/connectors/dummy-pay/%s", testCase.connectorID)
			}
			req := httptest.NewRequest(http.MethodDelete, endpoint, nil)
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

func TestInstall(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		body               []byte
		expectedStatusCode int
		expectedErrorCode  string
		managerError       error
	}

	dummypayConfig := dummypay.Config{
		Name:      "test",
		Directory: "test",
		FilePollingPeriod: connectors.Duration{
			Duration: 2 * time.Minute,
		},
	}

	body, err := json.Marshal(dummypayConfig)
	require.NoError(t, err)

	testCases := []testCase{
		{
			name: "nominal",
			body: body,
		},
		{
			name:               "invalid body",
			body:               []byte("invalid"),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name:               "manager error duplicate key value storage",
			body:               body,
			managerError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "manager error err not found storage",
			body:               body,
			managerError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "manager error already installed",
			body:               body,
			managerError:       manager.ErrAlreadyInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "manager error not installed",
			body:               body,
			managerError:       manager.ErrNotInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "manager error connector not found",
			body:               body,
			managerError:       manager.ErrConnectorNotFound,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "manager error err not found",
			body:               body,
			managerError:       manager.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "manager error err validation",
			body:               body,
			managerError:       manager.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "manager error other errors",
			body:               body,
			managerError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	for _, tc := range testCases {
		testCase := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusCreated
			}

			connectorID := models.ConnectorID{
				Reference: uuid.New(),
				Provider:  models.ConnectorProviderDummyPay,
			}

			expectedResponse := installResponse{
				ConnectorID: connectorID.String(),
			}

			backend, _ := newServiceTestingBackend(t)
			managerBackend, mockManager := newConnectorManagerTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockManager.EXPECT().
					Install(gomock.Any(), dummypayConfig.Name, dummypayConfig).
					Return(connectorID, nil)
			}

			if testCase.managerError != nil {
				mockManager.EXPECT().
					Install(gomock.Any(), dummypayConfig.Name, dummypayConfig).
					Return(models.ConnectorID{}, testCase.managerError)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), []connectorHandler{
				{
					Handler:  connectorRouter[dummypay.Config](models.ConnectorProviderDummyPay, managerBackend),
					Provider: models.ConnectorProviderDummyPay,
					initiatePayment: func(ctx context.Context, transfer *models.TransferInitiation) error {
						return nil
					},
				},
			})

			req := httptest.NewRequest(http.MethodPost, "/connectors/dummy-pay", bytes.NewReader(testCase.body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode >= 300 || testCase.expectedStatusCode < 200 {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			} else {
				var resp sharedapi.BaseResponse[installResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, &expectedResponse, resp.Data)
			}

		})
	}
}

func TestReset(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		connectorID        string
		connectors         map[string]*manager.ConnectorManager
		installed          *bool
		apiVersion         APIVersion
		expectedStatusCode int
		expectedErrorCode  string
		managerError       error
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	testCases := []testCase{
		// V0 tests
		{
			name:       "nominal V0",
			apiVersion: V0,
			connectors: map[string]*manager.ConnectorManager{
				connectorID.String(): nil,
			},
			installed: ptr(true),
		},
		{
			name: "too many connectors for V0",
			connectors: map[string]*manager.ConnectorManager{
				connectorID.String(): nil,
				"1":                  nil,
			},
			apiVersion:         V0,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "no connector for V0",
			connectors:         map[string]*manager.ConnectorManager{},
			apiVersion:         V0,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		// V1 tests
		{
			name:        "nominal V1",
			connectorID: connectorID.String(),
			apiVersion:  V1,
			installed:   ptr(true),
		},
		// Common test for V0 and V1
		{
			name:               "connector not installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(false),
		},
		{
			name:               "manager error duplicate key value storage",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
			installed:          ptr(true),
		},
		{
			name:               "manager error err not found storage",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			installed:          ptr(true),
		},
		{
			name:               "manager error already installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrAlreadyInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error not installed",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrNotInstalled,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error connector not found",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrConnectorNotFound,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error err not found",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			installed:          ptr(true),
		},
		{
			name:               "manager error err validation",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       manager.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			installed:          ptr(true),
		},
		{
			name:               "manager error other errors",
			connectorID:        connectorID.String(),
			apiVersion:         V1,
			managerError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
			installed:          ptr(true),
		},
	}

	for _, tc := range testCases {
		testCase := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusNoContent
			}

			backend, _ := newServiceTestingBackend(t)
			managerBackend, mockManager := newConnectorManagerTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockManager.EXPECT().
					Reset(gomock.Any(), connectorID).
					Return(nil)
			}

			if testCase.managerError != nil {
				mockManager.EXPECT().
					Reset(gomock.Any(), connectorID).
					Return(testCase.managerError)
			}

			if testCase.apiVersion == V0 {
				mockManager.EXPECT().
					Connectors().
					Return(testCase.connectors)
			}

			if testCase.installed != nil {
				mockManager.EXPECT().
					IsInstalled(gomock.Any(), connectorID).
					Return(*testCase.installed, nil)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), []connectorHandler{
				{
					Handler:  connectorRouter[dummypay.Config](models.ConnectorProviderDummyPay, managerBackend),
					Provider: models.ConnectorProviderDummyPay,
					initiatePayment: func(ctx context.Context, transfer *models.TransferInitiation) error {
						return nil
					},
				},
			})

			var endpoint string
			switch testCase.apiVersion {
			case V0:
				endpoint = "/connectors/dummy-pay/reset"
			case V1:
				endpoint = fmt.Sprintf("/connectors/dummy-pay/%s/reset", testCase.connectorID)
			}
			req := httptest.NewRequest(http.MethodPost, endpoint, nil)
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
