package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/service"
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"
)

func V2GetWaitingAttempts(w http.ResponseWriter, r *http.Request) {
	filterCursor := r.URL.Query().Get("cursor")

	resp := service.V2GetWaitingAttempts(filterCursor, attemptPageSize)

	if resp.Err != nil {
		if resp.T == utils.ValidationType {
			sharedapi.BadRequest(w, string(resp.T), resp.Err)
			return
		}
		if resp.T == utils.InternalType {
			sharedapi.InternalServerError(w, r, resp.Err)
			return
		}

		sharedapi.InternalServerError(w, r, resp.Err)
		return
	}

	sharedapi.RenderCursor(w, *resp.Data)

}
func V2GetAbortedAttempts(w http.ResponseWriter, r *http.Request) {
	filterCursor := r.URL.Query().Get("cursor")

	resp := service.V2GetAbortedAttempts(filterCursor, attemptPageSize)

	if resp.Err != nil {
		if resp.T == utils.ValidationType {
			sharedapi.BadRequest(w, string(resp.T), resp.Err)
			return
		}
		if resp.T == utils.InternalType {
			sharedapi.InternalServerError(w, r, resp.Err)
			return
		}

		sharedapi.InternalServerError(w, r, resp.Err)
		return
	}

	sharedapi.RenderCursor(w, *resp.Data)

	sharedapi.RenderCursor(w, *resp.Data)

}

func V2RetryWaitingAttempts(w http.ResponseWriter, r *http.Request) {

	resp := service.V2RetryWaitingAttempts()

	if resp.Err != nil {
		if resp.T == utils.ValidationType {
			sharedapi.BadRequest(w, string(resp.T), resp.Err)
			return
		}
		if resp.T == utils.InternalType {
			sharedapi.InternalServerError(w, r, resp.Err)
			return
		}

		sharedapi.InternalServerError(w, r, resp.Err)
		return
	}

	sharedapi.Ok(w, nil)
}

func V2RetryWaitingAttempt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp := service.V2RetryWaitingAttempt(id)

	if resp.Err != nil {
		if resp.T == utils.ValidationType {
			sharedapi.BadRequest(w, string(resp.T), resp.Err)
			return
		}
		if resp.T == utils.InternalType {
			sharedapi.InternalServerError(w, r, resp.Err)
			return
		}

		sharedapi.InternalServerError(w, r, resp.Err)
		return
	}

	sharedapi.Ok(w, nil)
}

func V2AbortWaitingAttempt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp := service.V2AbortWaitingAttempt(id)
	if resp.Err != nil {
		if resp.T == utils.ValidationType {
			sharedapi.BadRequest(w, string(resp.T), resp.Err)
			return
		}
		if resp.T == utils.InternalType {
			sharedapi.InternalServerError(w, r, resp.Err)
			return
		}
		if resp.T == utils.NotFoundType {
			sharedapi.NotFound(w, resp.Err)
			return
		}

		sharedapi.InternalServerError(w, r, resp.Err)
		return
	}

	sharedapi.Ok(w, *resp.Data)

}
