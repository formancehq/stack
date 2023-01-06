package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/app/storage"

	"github.com/google/uuid"

	"github.com/formancehq/payments/internal/app/models"

	"github.com/gorilla/mux"

	"github.com/formancehq/go-libs/sharedapi"
	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/formancehq/payments/internal/app/integration"
)

func handleErrorBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)

	sharedlogging.GetLogger(r.Context()).Error(err)
	// TODO: Opentracing
	err = json.NewEncoder(w).Encode(sharedapi.ErrorResponse{
		ErrorCode:    http.StatusText(http.StatusBadRequest),
		ErrorMessage: err.Error(),
	})
	if err != nil {
		panic(err)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	sharedlogging.GetLogger(r.Context()).Error(err)
	// TODO: Opentracing
	err = json.NewEncoder(w).Encode(sharedapi.ErrorResponse{
		ErrorCode:    "INTERNAL",
		ErrorMessage: err.Error(),
	})
	if err != nil {
		panic(err)
	}
}

func readConfig[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config, err := connectorManager.ReadConfig(r.Context())
		if err != nil {
			handleError(w, r, err)

			return
		}

		err = json.NewEncoder(w).Encode(config)
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
		skip, err := integerWithDefault(r, "skip", 0)
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		limit, err := integerWithDefault(r, "limit", maxPerPage)
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		if limit > maxPerPage {
			limit = maxPerPage
		}

		tasks, err := connectorManager.ListTasksStates(r.Context(), storage.Paginate(skip, limit))
		if err != nil {
			handleError(w, r, err)

			return
		}

		response := make([]listTasksResponseElement, len(tasks))
		for i, task := range tasks {
			response[i] = listTasksResponseElement{
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

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			panic(err)
		}
	}
}

func readTask[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID, err := uuid.Parse(mux.Vars(r)["taskID"])
		if err != nil {
			handleErrorBadRequest(w, r, err)

			return
		}

		task, err := connectorManager.ReadTaskState(r.Context(), taskID)
		if err != nil {
			handleError(w, r, err)

			return
		}

		response := listTasksResponseElement{
			ID:          task.ID.String(),
			ConnectorID: task.ConnectorID.String(),
			CreatedAt:   task.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
			Descriptor:  task.Descriptor,
			Status:      task.Status,
			State:       task.State,
			Error:       task.Error,
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			panic(err)
		}
	}
}

func uninstall[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := connectorManager.Uninstall(r.Context())
		if err != nil {
			handleError(w, r, err)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func install[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		installed, err := connectorManager.IsInstalled(context.Background())
		if err != nil {
			handleError(w, r, err)

			return
		}

		if installed {
			handleError(w, r, integration.ErrAlreadyInstalled)

			return
		}

		var config Config
		if r.ContentLength > 0 {
			err = json.NewDecoder(r.Body).Decode(&config)
			if err != nil {
				handleError(w, r, err)

				return
			}
		}

		err = connectorManager.Install(r.Context(), config)
		if err != nil {
			handleError(w, r, err)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func reset[Config models.ConnectorConfigObject](connectorManager *integration.ConnectorManager[Config],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		installed, err := connectorManager.IsInstalled(context.Background())
		if err != nil {
			handleError(w, r, err)

			return
		}

		if !installed {
			handleError(w, r, errors.New("connector not installed"))

			return
		}

		err = connectorManager.Reset(r.Context())
		if err != nil {
			handleError(w, r, err)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
