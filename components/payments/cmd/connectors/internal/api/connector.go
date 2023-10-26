package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/integration"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func readConfig[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if connectorNotInstalled(connectorManager, w, r) {
			return
		}

		config, err := connectorManager.ReadConfig(r.Context())
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[Config]{
			Data: config,
		})
		if err != nil {
			panic(err)
		}
	}
}

type listTasksResponseElement struct {
	ID          string            `json:"id"`
	ConnectorID string            `json:"connectorID"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
	Descriptor  json.RawMessage   `json:"descriptor"`
	Status      models.TaskStatus `json:"status"`
	State       json.RawMessage   `json:"state"`
	Error       string            `json:"error"`
}

func listTasks[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if connectorNotInstalled(connectorManager, w, r) {
			return
		}

		pageSize, err := pageSizeQueryParam(r)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		pagination, err := storage.Paginate(pageSize, r.URL.Query().Get("cursor"), nil, nil)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		tasks, paginationDetails, err := connectorManager.ListTasksStates(r.Context(), pagination)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		data := make([]listTasksResponseElement, len(tasks))
		for i, task := range tasks {
			data[i] = listTasksResponseElement{
				ID:          task.ID.String(),
				ConnectorID: task.ConnectorID.String(),
				CreatedAt:   task.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
				Descriptor:  task.Descriptor,
				Status:      task.Status,
				State:       task.State,
				Error:       task.Error,
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[listTasksResponseElement]{
			Cursor: &api.Cursor[listTasksResponseElement]{
				PageSize: paginationDetails.PageSize,
				HasMore:  paginationDetails.HasMore,
				Previous: paginationDetails.PreviousPage,
				Next:     paginationDetails.NextPage,
				Data:     data,
			},
		})
		if err != nil {
			panic(err)
		}
	}
}

func readTask[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if connectorNotInstalled(connectorManager, w, r) {
			return
		}

		taskID, err := uuid.Parse(mux.Vars(r)["taskID"])
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		task, err := connectorManager.ReadTaskState(r.Context(), taskID)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		data := listTasksResponseElement{
			ID:          task.ID.String(),
			ConnectorID: task.ConnectorID.String(),
			CreatedAt:   task.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
			Descriptor:  task.Descriptor,
			Status:      task.Status,
			State:       task.State,
			Error:       task.Error,
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[listTasksResponseElement]{
			Data: &data,
		})
		if err != nil {
			panic(err)
		}
	}
}

func uninstall[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if connectorNotInstalled(connectorManager, w, r) {
			return
		}

		err := connectorManager.Uninstall(r.Context())
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func install[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		installed, err := connectorManager.IsInstalled(r.Context())
		if err != nil {
			handleStorageErrors(w, r, err)
			return
		}

		if installed {
			api.BadRequest(w, ErrValidation, integration.ErrAlreadyInstalled)
			return
		}

		var config Config
		if r.ContentLength > 0 {
			err = json.NewDecoder(r.Body).Decode(&config)
			if err != nil {
				api.BadRequest(w, ErrMissingOrInvalidBody, err)
				return
			}
		}

		err = connectorManager.Install(r.Context(), config)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func reset[Config models.ConnectorConfigObject](
	connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if connectorNotInstalled(connectorManager, w, r) {
			return
		}

		err := connectorManager.Reset(r.Context())
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func connectorNotInstalled[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
	w http.ResponseWriter, r *http.Request,
) bool {
	installed, err := connectorManager.IsInstalled(r.Context())
	if err != nil {
		handleStorageErrors(w, r, err)
		return true
	}

	if !installed {
		api.BadRequest(w, ErrValidation, integration.ErrNotInstalled)
		return true
	}

	return false
}
