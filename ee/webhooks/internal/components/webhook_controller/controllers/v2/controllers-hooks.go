package controllers

import (
	"errors"
	"fmt"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/go-chi/chi/v5"
	"net/http"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"

	"github.com/formancehq/webhooks/internal/commons"
	controllersCommons "github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/commons"
	"github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/utils"
	r "github.com/formancehq/webhooks/internal/components/webhook_controller/routes"

	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	serverInterfaces "github.com/formancehq/webhooks/internal/services/httpserver/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)

const (
	pageSize int = 16
)

type PayloadBody struct {
	Payload string `json:"payload"`
}

func RegisterV2HookControllers(server serverInterfaces.IHTTPServer, database storeInterface.IStoreProvider, client clientInterface.IHTTPClient) {

	server.Register(string(r.V2CreateHook.Method), r.V2CreateHook.Url, func(w http.ResponseWriter, r *http.Request) {

		hookParams := commons.HookBodyParams{}

		if err := utils.DecodeJSONBody(r, &hookParams); err != nil {
			sharedapi.BadRequest(w, utils.ErrValidation, err)
			return
		}

		resp := V2CreateHookController(database, hookParams)

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

	})

	server.Register(string(r.V2GetHooks.Method), r.V2GetHooks.Url, func(w http.ResponseWriter, r *http.Request) {

		filterEndpoint := r.URL.Query().Get("endpoint")
		filterCursor := r.URL.Query().Get("cursor")

		resp := V2GetHooksController(database, filterEndpoint, filterCursor)

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

	})

	server.Register(string(r.V2DeleteHook.Method), r.V2DeleteHook.Url, func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		resp := V2DeleteHookController(database, id)

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

	})

	server.Register(string(r.V2ActiveHook.Method), r.V2ActiveHook.Url, func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		resp := V2ActiveHookController(database, id)

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
	})

	server.Register(string(r.V2DeactiveHook.Method), r.V2DeactiveHook.Url, func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		resp := V2DeactiveHookController(database, id)

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
	})

	server.Register(string(r.V2ChangeHookSecret.Method), r.V2ChangeHookSecret.Url, func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		sec := &utils.Secret{}

		if err := utils.DecodeJSONBody(r, &sec); err != nil {
			sharedapi.BadRequest(w, utils.ErrValidation, err)
			return
		}

		resp := V2ChangeSecretController(database, id, sec.Secret)

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

	})

	server.Register(string(r.V2TestHook.Method), r.V2TestHook.Url, func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		payload := PayloadBody{}

		if err := utils.DecodeJSONBody(r, &payload); err != nil {
			sharedapi.BadRequest(w, utils.ErrValidation, err)
			return
		}

		resp := V2TestHookController(database, client, id, payload.Payload)

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

	})

	server.Register(string(r.V2ChangeHookEndpoint.Method), r.V2ChangeHookEndpoint.Url, func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		ep := &utils.Endpoint{}

		if err := utils.DecodeJSONBody(r, &ep); err != nil {
			sharedapi.BadRequest(w, utils.ErrValidation, err)
			return
		}

		if err := utils.ValidateEndpoint(ep.Endpoint); err != nil {
			sharedapi.BadRequest(w, utils.ErrValidation, err)
		}

		hook, err := controllersCommons.UpdateEndpoint(database, id, ep.Endpoint)

		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		if hook.ID == "" {
			sharedapi.NotFound(w, errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
			return
		}

		sharedapi.Ok(w, hook)
	})

	server.Register(string(r.V2ChangeHookRetry.Method), r.V2ChangeHookRetry.Url, func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		retry := &utils.Retry{}
		if err := utils.DecodeJSONBody(r, &retry); err != nil {
			sharedapi.BadRequest(w, utils.ErrValidation, err)
			return
		}

		hook, err := controllersCommons.UpdateRetry(database, id, retry.Retry)

		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		if hook.ID == "" {
			sharedapi.NotFound(w, errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
			return
		}

		sharedapi.Ok(w, hook)

	})
}

func V2CreateHookController(database storeInterface.IStoreProvider, hookParams commons.HookBodyParams) utils.Response[commons.Hook] {

	if err := utils.ValidateEndpoint(hookParams.Endpoint); err != nil {
		return utils.ValidationErrorResp[commons.Hook](err)
	}

	if err := utils.ValidateSecret(&hookParams.Secret); err != nil {
		return utils.ValidationErrorResp[commons.Hook](err)
	}

	if err := utils.FormatEvents(&hookParams.Events); err != nil {
		return utils.ValidationErrorResp[commons.Hook](err)
	}

	hook, err := controllersCommons.CreateHook(database, hookParams.Name, hookParams.Events, hookParams.Endpoint, hookParams.Secret, hookParams.Retry)

	if err != nil {
		return utils.InternalErrorResp[commons.Hook](err)
	}

	return utils.SuccessResp(hook)
}

func V2GetHooksController(database storeInterface.IStoreProvider, filterEndpoint, filterCursor string) utils.Response[bunpaginate.Cursor[commons.Hook]] {
	hasMore := false
	strPrevious := ""
	strNext := ""

	cursor, err := utils.ReadCursor(filterCursor)

	if err != nil {
		return utils.ValidationErrorResp[bunpaginate.Cursor[commons.Hook]](err)
	}

	if filterEndpoint != "" {
		if err := utils.ValidateEndpoint(filterEndpoint); err != nil {
			return utils.ValidationErrorResp[bunpaginate.Cursor[commons.Hook]](err)
		}
	}

	hooks, hM, err := controllersCommons.GetHooks(database, filterEndpoint, cursor, pageSize)
	if err != nil {
		return utils.InternalErrorResp[bunpaginate.Cursor[commons.Hook]](err)
	}
	hasMore = hM

	if hasMore {
		strPrevious, strNext = utils.PaginationCursor(cursor, hasMore)
	}

	Cursor := bunpaginate.Cursor[commons.Hook]{
		PageSize: pageSize,
		HasMore:  hasMore,
		Previous: strPrevious,
		Next:     strNext,
		Data:     utils.ToValues(*hooks),
	}

	return utils.SuccessResp(Cursor)
}

func V2DeleteHookController(database storeInterface.IStoreProvider, id string) utils.Response[commons.Hook] {
	hook, err := controllersCommons.DeleteHook(database, id)
	if err != nil {
		return utils.InternalErrorResp[commons.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[commons.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}

func V2ActiveHookController(database storeInterface.IStoreProvider, id string) utils.Response[commons.Hook] {
	hook, err := controllersCommons.ActivateHook(database, id)
	if err != nil {
		return utils.InternalErrorResp[commons.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[commons.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}

func V2DeactiveHookController(database storeInterface.IStoreProvider, id string) utils.Response[commons.Hook] {
	hook, err := controllersCommons.DeactivateHook(database, id)
	if err != nil {
		return utils.InternalErrorResp[commons.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[commons.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}

func V2ChangeSecretController(database storeInterface.IStoreProvider, id, secret string) utils.Response[commons.Hook] {
	if err := utils.ValidateSecret(&secret); err != nil {
		return utils.ValidationErrorResp[commons.Hook](err)
	}

	hook, err := controllersCommons.UpdateSecret(database, id, secret)

	if err != nil {
		return utils.InternalErrorResp[commons.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[commons.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}

func V2TestHookController(database storeInterface.IStoreProvider, client clientInterface.IHTTPClient, id, payload string) utils.Response[commons.Attempt] {
	hook, attempt, err := controllersCommons.TestHook(database, client, id, payload)

	if err != nil {
		return utils.InternalErrorResp[commons.Attempt](err)
	}

	if hook.ID == "" {
		return utils.NotFoundErrorResp[commons.Attempt](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(*attempt)
}

func V2ChangeEndpointController(database storeInterface.IStoreProvider, id, endpoint string) utils.Response[commons.Hook] {
	if err := utils.ValidateEndpoint(endpoint); err != nil {
		return utils.ValidationErrorResp[commons.Hook](err)
	}

	hook, err := controllersCommons.UpdateEndpoint(database, id, endpoint)

	if err != nil {
		return utils.InternalErrorResp[commons.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[commons.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}

func V2ChangeRetryController(database storeInterface.IStoreProvider, id string, retry bool) utils.Response[commons.Hook] {

	hook, err := controllersCommons.UpdateRetry(database, id, retry)

	if err != nil {
		return utils.InternalErrorResp[commons.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[commons.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}
