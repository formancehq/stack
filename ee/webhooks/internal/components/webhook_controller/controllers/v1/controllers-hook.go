package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"

	controllersCommons "github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/commons"
	"github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/utils"
	r "github.com/formancehq/webhooks/internal/components/webhook_controller/routes"

	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	serverInterfaces "github.com/formancehq/webhooks/internal/services/httpserver/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)

const (
	pageSize int = 20
)

type PayloadBody struct {
	Payload string `json:"payload"`
}

func RegisterV1HookControllers(server serverInterfaces.IHTTPServer, database storeInterface.IStoreProvider, client clientInterface.IHTTPClient) {

	server.Register(string(r.V1CreateHook.Method), r.V1CreateHook.Url, func(w http.ResponseWriter, r *http.Request) {
		v1HU := utils.V1HookUser{}

		if err := utils.DecodeJSONBody(r, &v1HU); err != nil {
			sharedapi.BadRequest(w, utils.ErrValidation, err)
			return
		}

		resp := V1CreateHookController(database, v1HU)

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

	server.Register(string(r.V1GetHooks.Method), r.V1GetHooks.Url, func(w http.ResponseWriter, r *http.Request) {

		filterEndpoint := r.URL.Query().Get("endpoint")
		filterId := r.URL.Query().Get("id")
		filterCursor := r.URL.Query().Get("cursor")

		resp := V1GetHooksController(database, filterEndpoint, filterId, filterCursor)

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

	server.Register(string(r.V1DeleteHook.Method), r.V1DeleteHook.Url, func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		resp := V1DeleteHookController(database, id)

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

	server.Register(string(r.V1ActiveHook.Method), r.V1ActiveHook.Url, func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		resp := V1ActiveHookController(database, id)

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

	server.Register(string(r.V1DeactiveHook.Method), r.V1DeactiveHook.Url, func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		resp := V1DeactiveHookController(database, id)

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

	server.Register(string(r.V1ChangeSecret.Method), r.V1ChangeSecret.Url, func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		sec := &utils.Secret{}

		if err := utils.DecodeJSONBody(r, &sec); err != nil {
			sharedapi.BadRequest(w, utils.ErrValidation, err)
			return
		}

		resp := V1ChangeSecretController(database, id, sec.Secret)

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

	server.Register(string(r.V1TestHook.Method), string(r.V1TestHook.Method), func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		payload := PayloadBody{}

		if err := utils.DecodeJSONBody(r, &payload); err != nil {
			sharedapi.BadRequest(w, utils.ErrValidation, err)
			return
		}

		resp := V1TestHookController(database, client, id, payload.Payload)

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

}

func V1CreateHookController(database storeInterface.IStoreProvider, hook utils.V1HookUser) utils.Response[utils.V1Hook] {

	if err := utils.ValidateEndpoint(hook.Endpoint); err != nil {
		return utils.ValidationErrorResp[utils.V1Hook](err)
	}

	if err := utils.ValidateSecret(&hook.Secret); err != nil {
		return utils.ValidationErrorResp[utils.V1Hook](err)
	}

	if err := utils.FormatEvents(&hook.EventTypes); err != nil {
		return utils.ValidationErrorResp[utils.V1Hook](err)
	}

	newHook, err := controllersCommons.CreateHook(database, "", hook.EventTypes, hook.Endpoint, hook.Secret, true)

	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}

	return utils.SuccessResp(utils.ToV1Hook(newHook))

}

func V1GetHooksController(database storeInterface.IStoreProvider, filterEndpoint, filterId, filterCursor string) utils.Response[bunpaginate.Cursor[utils.V1Hook]] {
	v1Hooks := make([]utils.V1Hook, 0)
	hasMore := false

	strPrevious := ""
	strNext := ""

	cursor, err := utils.ReadCursor(filterCursor)

	if err != nil {
		return utils.ValidationErrorResp[bunpaginate.Cursor[utils.V1Hook]](err)
	}

	if filterEndpoint != "" {
		if err := utils.ValidateEndpoint(filterEndpoint); err != nil {
			return utils.ValidationErrorResp[bunpaginate.Cursor[utils.V1Hook]](err)
		}
	}

	if filterId != "" {
		hook, err := controllersCommons.GetHook(database, filterId)
		if err != nil {
			return utils.InternalErrorResp[bunpaginate.Cursor[utils.V1Hook]](err)
		}

		if hook.ID != "" {
			v1Hooks = append(v1Hooks, utils.ToV1Hook(hook))
		}

	} else {
		temps, hM, err := controllersCommons.GetHooks(database, filterEndpoint, cursor, pageSize)
		if err != nil {
			return utils.InternalErrorResp[bunpaginate.Cursor[utils.V1Hook]](err)
		}
		hasMore = hM
		v1Hooks = append(v1Hooks, utils.ToV1Hooks(temps)...)
	}

	if hasMore {
		strPrevious, strNext = utils.PaginationCursor(cursor, hasMore)
	}

	Cursor := bunpaginate.Cursor[utils.V1Hook]{
		HasMore:  hasMore,
		Previous: strPrevious,
		Next:     strNext,
		Data:     v1Hooks,
	}

	return utils.SuccessResp(Cursor)
}

func V1DeleteHookController(database storeInterface.IStoreProvider, id string) utils.Response[utils.V1Hook] {
	hook, err := controllersCommons.DeleteHook(database, id)
	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[utils.V1Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(utils.ToV1Hook(hook))
}

func V1ActiveHookController(database storeInterface.IStoreProvider, id string) utils.Response[utils.V1Hook] {
	hook, err := controllersCommons.ActivateHook(database, id)
	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[utils.V1Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(utils.ToV1Hook(hook))
}

func V1DeactiveHookController(database storeInterface.IStoreProvider, id string) utils.Response[utils.V1Hook] {
	hook, err := controllersCommons.DeactivateHook(database, id)
	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[utils.V1Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(utils.ToV1Hook(hook))
}

func V1ChangeSecretController(database storeInterface.IStoreProvider, id, secret string) utils.Response[utils.V1Hook] {
	if err := utils.ValidateSecret(&secret); err != nil {
		return utils.ValidationErrorResp[utils.V1Hook](err)
	}

	hook, err := controllersCommons.UpdateSecret(database, id, secret)

	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[utils.V1Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(utils.ToV1Hook(hook))

}

func V1TestHookController(database storeInterface.IStoreProvider, client clientInterface.IHTTPClient, id, payload string) utils.Response[utils.V1Attempt] {
	hook, attempt, err := controllersCommons.TestHook(database, client, id, payload)

	if err != nil {
		return utils.InternalErrorResp[utils.V1Attempt](err)
	}

	if hook.ID == "" {
		return utils.NotFoundErrorResp[utils.V1Attempt](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(utils.ToV1Attempt(*hook, *attempt))
}
