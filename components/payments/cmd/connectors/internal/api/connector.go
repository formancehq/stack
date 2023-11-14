package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/integration"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type APIVersion int

const (
	V0 APIVersion = iota
	V1 APIVersion = iota
)

func getConnectorID[Config models.ConnectorConfigObject](
	connectorManager *integration.ConnectorsManager[Config],
	r *http.Request,
	apiVersion APIVersion,
) (models.ConnectorID, error) {
	switch apiVersion {
	case V0:
		connectors := connectorManager.Connectors()
		if len(connectors) == 0 {
			return models.ConnectorID{}, fmt.Errorf("no connectors installed")
		}

		if len(connectors) > 1 {
			return models.ConnectorID{}, fmt.Errorf("more than one connectors installed")
		}

		for id := range connectors {
			return models.MustConnectorIDFromString(id), nil
		}

	case V1:
		connectorID, err := models.ConnectorIDFromString(mux.Vars(r)["connectorID"])
		if err != nil {
			return models.ConnectorID{}, err
		}

		return connectorID, nil
	}

	return models.ConnectorID{}, fmt.Errorf("unknown API version")
}

func readConfig[Config models.ConnectorConfigObject](
	connectorManager *integration.ConnectorsManager[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connectorID, err := getConnectorID(connectorManager, r, apiVersion)
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if connectorNotInstalled(connectorManager, connectorID, w, r) {
			return
		}

		config, err := connectorManager.ReadConfig(r.Context(), connectorID)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[Config]{
			Data: &config,
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

func listTasks[Config models.ConnectorConfigObject](
	connectorManager *integration.ConnectorsManager[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connectorID, err := getConnectorID(connectorManager, r, apiVersion)
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if connectorNotInstalled(connectorManager, connectorID, w, r) {
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

		tasks, paginationDetails, err := connectorManager.ListTasksStates(r.Context(), connectorID, pagination)
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

func readTask[Config models.ConnectorConfigObject](
	connectorManager *integration.ConnectorsManager[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connectorID, err := getConnectorID(connectorManager, r, apiVersion)
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if connectorNotInstalled(connectorManager, connectorID, w, r) {
			return
		}

		taskID, err := uuid.Parse(mux.Vars(r)["taskID"])
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		task, err := connectorManager.ReadTaskState(r.Context(), connectorID, taskID)
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

func uninstall[Config models.ConnectorConfigObject](
	connectorManager *integration.ConnectorsManager[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connectorID, err := getConnectorID(connectorManager, r, apiVersion)
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if connectorNotInstalled(connectorManager, connectorID, w, r) {
			return
		}

		err = connectorManager.Uninstall(r.Context(), connectorID)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type installResponse struct {
	ConnectorID string `json:"connectorID"`
}

func install[Config models.ConnectorConfigObject](
	connectorManager *integration.ConnectorsManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var config Config
		if r.ContentLength > 0 {
			err := json.NewDecoder(r.Body).Decode(&config)
			if err != nil {
				api.BadRequest(w, ErrMissingOrInvalidBody, err)
				return
			}
		}

		connectorID, err := connectorManager.Install(r.Context(), config.ConnectorName(), config)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(api.BaseResponse[installResponse]{
			Data: &installResponse{
				ConnectorID: connectorID.String(),
			},
		})
		if err != nil {
			panic(err)
		}
	}
}

func reset[Config models.ConnectorConfigObject](
	connectorManager *integration.ConnectorsManager[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connectorID, err := getConnectorID(connectorManager, r, apiVersion)
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if connectorNotInstalled(connectorManager, connectorID, w, r) {
			return
		}

		err = connectorManager.Reset(r.Context(), connectorID)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func connectorNotInstalled[Config models.ConnectorConfigObject](
	connectorManager *integration.ConnectorsManager[Config],
	connectorID models.ConnectorID,
	w http.ResponseWriter, r *http.Request,
) bool {
	installed, err := connectorManager.IsInstalled(r.Context(), connectorID)
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
