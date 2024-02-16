package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type APIVersion int

const (
	V0 APIVersion = iota
	V1 APIVersion = iota
)

func (a APIVersion) String() string {
	switch a {
	case V0:
		return "v0"
	case V1:
		return "v1"
	default:
		return "unknown"
	}
}

func updateConfig[Config models.ConnectorConfigObject](
	b backend.ManagerBackend[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "updateConfig")
		defer span.End()

		span.SetAttributes(attribute.String("apiVersion", apiVersion.String()))

		connectorID, err := getConnectorID(span, b, r, apiVersion)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if connectorNotInstalled(span, b, connectorID, w, r) {
			return
		}

		var config Config
		if r.ContentLength > 0 {
			err := json.NewDecoder(r.Body).Decode(&config)
			if err != nil {
				otel.RecordError(span, err)
				api.BadRequest(w, ErrMissingOrInvalidBody, err)
				return
			}
		}

		err = b.GetManager().UpdateConfig(ctx, connectorID, config)
		if err != nil {
			otel.RecordError(span, err)
			handleConnectorsManagerErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}

func readConfig[Config models.ConnectorConfigObject](
	b backend.ManagerBackend[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "readConfig")
		defer span.End()

		span.SetAttributes(attribute.String("apiVersion", apiVersion.String()))

		connectorID, err := getConnectorID(span, b, r, apiVersion)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("connectorID", connectorID.String()))

		if connectorNotInstalled(span, b, connectorID, w, r) {
			return
		}

		config, err := b.GetManager().ReadConfig(ctx, connectorID)
		if err != nil {
			otel.RecordError(span, err)
			handleConnectorsManagerErrors(w, r, err)
			return
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[Config]{
			Data: &config,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
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
	b backend.ManagerBackend[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "listTasks")
		defer span.End()

		span.SetAttributes(attribute.String("apiVersion", apiVersion.String()))

		connectorID, err := getConnectorID(span, b, r, apiVersion)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if connectorNotInstalled(span, b, connectorID, w, r) {
			return
		}

		query, err := bunpaginate.Extract[storage.ListTasksQuery](r, func() (*storage.ListTasksQuery, error) {
			pageSize, err := bunpaginate.GetPageSize(r)
			if err != nil {
				return nil, err
			}

			return pointer.For(storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{}).WithPageSize(pageSize))), nil
		})
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		span.SetAttributes(attribute.Int("pageSize", int(query.PageSize)))
		span.SetAttributes(attribute.String("cursor", r.URL.Query().Get("cursor")))

		cursor, err := b.GetManager().ListTasksStates(ctx, connectorID, *query)
		if err != nil {
			otel.RecordError(span, err)
			handleConnectorsManagerErrors(w, r, err)
			return
		}

		tasks := cursor.Data
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
				PageSize: cursor.PageSize,
				HasMore:  cursor.HasMore,
				Previous: cursor.Previous,
				Next:     cursor.Next,
				Data:     data,
			},
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func webhooksMiddleware[Config models.ConnectorConfigObject](
	b backend.ManagerBackend[Config],
	apiVersion APIVersion,
) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := otel.Tracer().Start(r.Context(), "webhooksMiddleware")
			defer span.End()

			connectorID, err := getConnectorID(span, b, r, apiVersion)
			if err != nil {
				otel.RecordError(span, err)
				api.BadRequest(w, ErrInvalidID, err)
				return
			}

			if connectorNotInstalled(span, b, connectorID, w, r) {
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				otel.RecordError(span, err)
				api.BadRequest(w, ErrMissingOrInvalidBody, err)
			}
			defer r.Body.Close()

			webhook := &models.Webhook{
				ID:          uuid.New(),
				ConnectorID: connectorID,
				RequestBody: body,
			}

			span.SetAttributes(attribute.String("webhook.id", webhook.ID.String()))

			ctx, err = b.GetManager().CreateWebhookAndContext(ctx, webhook)
			if err != nil {
				otel.RecordError(span, err)
				handleConnectorsManagerErrors(w, r, err)
				return
			}

			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func readTask[Config models.ConnectorConfigObject](
	b backend.ManagerBackend[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "readTask")
		defer span.End()

		span.SetAttributes(attribute.String("apiVersion", apiVersion.String()))

		connectorID, err := getConnectorID(span, b, r, apiVersion)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if connectorNotInstalled(span, b, connectorID, w, r) {
			return
		}

		taskID, err := uuid.Parse(mux.Vars(r)["taskID"])
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("taskID", taskID.String()))

		task, err := b.GetManager().ReadTaskState(ctx, connectorID, taskID)
		if err != nil {
			otel.RecordError(span, err)
			handleConnectorsManagerErrors(w, r, err)
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
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func uninstall[Config models.ConnectorConfigObject](
	b backend.ManagerBackend[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "uninstall")
		defer span.End()

		span.SetAttributes(attribute.String("apiVersion", apiVersion.String()))

		connectorID, err := getConnectorID(span, b, r, apiVersion)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if connectorNotInstalled(span, b, connectorID, w, r) {
			return
		}

		err = b.GetManager().Uninstall(ctx, connectorID)
		if err != nil {
			otel.RecordError(span, err)
			handleConnectorsManagerErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type installResponse struct {
	ConnectorID string `json:"connectorID"`
}

func install[Config models.ConnectorConfigObject](b backend.ManagerBackend[Config]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "install")
		defer span.End()

		w.Header().Set("Content-Type", "application/json")

		var config Config
		if r.ContentLength > 0 {
			err := json.NewDecoder(r.Body).Decode(&config)
			if err != nil {
				otel.RecordError(span, err)
				api.BadRequest(w, ErrMissingOrInvalidBody, err)
				return
			}
		}

		connectorID, err := b.GetManager().Install(ctx, config.ConnectorName(), config)
		if err != nil {
			otel.RecordError(span, err)
			handleConnectorsManagerErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(api.BaseResponse[installResponse]{
			Data: &installResponse{
				ConnectorID: connectorID.String(),
			},
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func reset[Config models.ConnectorConfigObject](
	b backend.ManagerBackend[Config],
	apiVersion APIVersion,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "reset")
		defer span.End()

		span.SetAttributes(attribute.String("apiVersion", apiVersion.String()))

		connectorID, err := getConnectorID(span, b, r, apiVersion)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if connectorNotInstalled(span, b, connectorID, w, r) {
			return
		}

		err = b.GetManager().Reset(ctx, connectorID)
		if err != nil {
			otel.RecordError(span, err)
			handleConnectorsManagerErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func connectorNotInstalled[Config models.ConnectorConfigObject](
	span trace.Span,
	b backend.ManagerBackend[Config],
	connectorID models.ConnectorID,
	w http.ResponseWriter, r *http.Request,
) bool {
	installed, err := b.GetManager().IsInstalled(r.Context(), connectorID)
	if err != nil {
		otel.RecordError(span, err)
		handleConnectorsManagerErrors(w, r, err)
		return true
	}

	if !installed {
		otel.RecordError(span, fmt.Errorf("connector not installed"))
		api.BadRequest(w, ErrValidation, fmt.Errorf("connector not installed"))
		return true
	}

	return false
}

func getConnectorID[Config models.ConnectorConfigObject](
	span trace.Span,
	b backend.ManagerBackend[Config],
	r *http.Request,
	apiVersion APIVersion,
) (models.ConnectorID, error) {
	switch apiVersion {
	case V0:
		connectors := b.GetManager().Connectors()
		if len(connectors) == 0 {
			return models.ConnectorID{}, fmt.Errorf("no connectors installed")
		}

		span.SetAttributes(attribute.Int("connectors.count", len(connectors)))

		if len(connectors) > 1 {
			return models.ConnectorID{}, fmt.Errorf("more than one connectors installed")
		}

		for id := range connectors {
			return models.MustConnectorIDFromString(id), nil
		}

	case V1:
		c := mux.Vars(r)["connectorID"]

		span.SetAttributes(attribute.String("connectorID", c))

		connectorID, err := models.ConnectorIDFromString(c)
		if err != nil {
			return models.ConnectorID{}, err
		}

		return connectorID, nil
	}

	return models.ConnectorID{}, fmt.Errorf("unknown API version")
}
