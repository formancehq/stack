package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/service"
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"
)

func V1CreateHook(w http.ResponseWriter, r *http.Request) {
	v1HU := utils.V1HookUser{}

	if err := utils.DecodeJSONBody(r, &v1HU); err != nil {
		sharedapi.BadRequest(w, utils.ErrValidation, err)
		return
	}

	resp := service.V1CreateHook(v1HU)

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

	sharedapi.Ok(w, *resp.Data)
	return

}

func V1GetHooks(w http.ResponseWriter, r *http.Request) {

	filterEndpoint := r.URL.Query().Get("endpoint")
	filterId := r.URL.Query().Get("id")
	filterCursor := r.URL.Query().Get("cursor")

	resp := service.V1GetHooks(filterEndpoint, filterId, filterCursor, hookPageSize)

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

func DeleteHook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp := service.V1DeleteHook(id)

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
	return

}

func V1ActivateHook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp := service.V1ActiveHook(id)

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
	return

}

func V1DeactivateHook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp := service.V1DeactiveHook(id)

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
	return
}

func V1ChangeHookSecret(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sec := &utils.Secret{}

	if err := utils.DecodeJSONBody(r, &sec); err != nil {
		sharedapi.BadRequest(w, utils.ErrValidation, err)
		return
	}

	resp := service.V1ChangeSecret(id, sec.Secret)

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
	return

}

func V1TestHook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	payload := PayloadBody{}

	payload.Payload = "{\"data\":\"test\"}"

	resp := service.V1TestHook(id, payload.Payload)

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
	return

}
