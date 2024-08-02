package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"

	"github.com/formancehq/webhooks/internal/models"

	"github.com/formancehq/webhooks/internal/app/webhook_server/api/service"
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"
)

func V2CreateHook(w http.ResponseWriter, r *http.Request) {

	hookParams := models.HookBodyParams{}
	hookParams.Retry = true

	if err := utils.DecodeJSONBody(r, &hookParams); err != nil {
		sharedapi.BadRequest(w, utils.ErrValidation, err)
		return
	}

	resp := service.V2CreateHook(hookParams)

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

	sharedapi.Created(w, *resp.Data)
	return
}

func V2GetHook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	resp := service.V2GetHook(id)

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

func V2GetHooks(w http.ResponseWriter, r *http.Request) {

	filterEndpoint := r.URL.Query().Get("endpoint")
	filterCursor := r.URL.Query().Get("cursor")

	resp := service.V2GetHooks(filterEndpoint, filterCursor, hookPageSize)

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

func V2DeleteHook(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	resp := service.V2DeleteHook(id)

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

func V2ActivateHook(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	resp := service.V2ActiveHook(id)

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

func V2DeactivateHook(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	resp := service.V2DeactiveHook(id)

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

func V2ChangeHookSecret(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	sec := &utils.Secret{}

	if err := utils.DecodeJSONBody(r, &sec); err != nil {
		sharedapi.BadRequest(w, utils.ErrValidation, err)
		return
	}

	resp := service.V2ChangeSecret(id, sec.Secret)

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

func V2TestHook(w http.ResponseWriter, r *http.Request) {
	logging.Infof("V2TEStHook")
	id := chi.URLParam(r, "id")
	payload := PayloadBody{}

	if err := utils.DecodeJSONBody(r, &payload); err != nil {
		sharedapi.BadRequest(w, utils.ErrValidation, err)
		return
	}

	resp := service.V2TestHook(id, payload.Payload)

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

func V2ChangeHookEndpoint(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	ep := &utils.Endpoint{}

	if err := utils.DecodeJSONBody(r, &ep); err != nil {
		sharedapi.BadRequest(w, utils.ErrValidation, err)
		return
	}

	if err := utils.ValidateEndpoint(ep.Endpoint); err != nil {
		sharedapi.BadRequest(w, utils.ErrValidation, err)
	}

	resp := service.V2ChangeEndpoint(id, ep.Endpoint)

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

func V2ChangeHookRetry(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	retry := &utils.Retry{}
	if err := utils.DecodeJSONBody(r, &retry); err != nil {
		sharedapi.BadRequest(w, utils.ErrValidation, err)
		return
	}

	resp := service.V2ChangeRetry(id, retry.Retry)

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
