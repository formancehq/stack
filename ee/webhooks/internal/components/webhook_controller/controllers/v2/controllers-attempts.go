package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/webhooks/internal/commons"
	"github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/utils"
	r "github.com/formancehq/webhooks/internal/components/webhook_controller/routes"

	serverInterfaces "github.com/formancehq/webhooks/internal/services/httpserver/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)
const (
	attemptPageSize int = 64
)

func RegisterV2AttemptControllers(serverHttp serverInterfaces.IHTTPServer, database storeInterface.IStoreProvider) {

	serverHttp.Register(string(r.V2GetWaitingAttempts.Method), r.V2GetWaitingAttempts.Url, func(w http.ResponseWriter, r *http.Request) {
		filterCursor := r.URL.Query().Get("cursor")

		resp := V2GetWaitingAttemptsController(database, filterCursor)

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

	serverHttp.Register(string(r.V2GetAbortedAttempts.Method), r.V2GetAbortedAttempts.Url, func(w http.ResponseWriter, r *http.Request) {
		filterCursor := r.URL.Query().Get("cursor")

		resp := V2GetAbortedAttemptsController(database, filterCursor)

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

	serverHttp.Register(string(r.V2RetryWaitingAttempts.Method), r.V2RetryWaitingAttempts.Url, func(w http.ResponseWriter, r *http.Request) {

		resp := V2RetryWaitingAttemptsController(database)

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
	})

	serverHttp.Register(string(r.V2RetryWaitingAttempt.Method), r.V2RetryWaitingAttempt.Url, func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		resp := V2RetryWaitingAttemptController(database, id)

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
	})

	serverHttp.Register(string(r.V2AbortWaitingAttempt.Method), r.V2AbortWaitingAttempt.Url, func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		resp := V2AbortWaitingAttemptController(database, id)
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

	})

}

func V2GetWaitingAttemptsController(database storeInterface.IStoreProvider, filterCursor string) utils.Response[bunpaginate.Cursor[commons.Attempt]] {
	hasMore := false

	strPrevious := ""
	strNext := ""

	cursor, err := utils.ReadCursor(filterCursor)

	if err != nil {
		return utils.ValidationErrorResp[bunpaginate.Cursor[commons.Attempt]](err)
	}

	attempts, hM, err := database.GetWaitingAttempts(cursor, attemptPageSize)
	if err != nil {
		return utils.InternalErrorResp[bunpaginate.Cursor[commons.Attempt]](err)
	}

	hasMore = hM

	if hasMore {
		strPrevious, strNext = utils.PaginationCursor(cursor, hasMore)
	}

	Cursor := bunpaginate.Cursor[commons.Attempt]{
		HasMore:  hasMore,
		Previous: strPrevious,
		Next:     strNext,
		Data:     utils.ToValues(*attempts),
	}

	return utils.SuccessResp[bunpaginate.Cursor[commons.Attempt]](Cursor)
}

func V2GetAbortedAttemptsController(database storeInterface.IStoreProvider, filterCursor string) utils.Response[bunpaginate.Cursor[commons.Attempt]] {
	hasMore := false

	strPrevious := ""
	strNext := ""

	cursor, err := utils.ReadCursor(filterCursor)

	if err != nil {
		return utils.ValidationErrorResp[bunpaginate.Cursor[commons.Attempt]](err)
	}

	attempts, hM, err := database.GetAbortedAttempts(cursor, attemptPageSize)
	if err != nil {
		return utils.InternalErrorResp[bunpaginate.Cursor[commons.Attempt]](err)
	}

	hasMore = hM

	if hasMore {
		strPrevious, strNext = utils.PaginationCursor(cursor, hasMore)
	}

	Cursor := bunpaginate.Cursor[commons.Attempt]{
		HasMore:  hasMore,
		Previous: strPrevious,
		Next:     strNext,
		Data:     utils.ToValues(*attempts),
	}

	return utils.SuccessResp[bunpaginate.Cursor[commons.Attempt]](Cursor)
}

func V2RetryWaitingAttemptsController(database storeInterface.IStoreProvider) utils.Response[any] {

	ev, err := commons.EventFromType(commons.FlushWaitingAttemptsType, nil, nil)
	if err != nil {
		return utils.InternalErrorResp[any](err)
	}

	log, err := commons.LogFromEvent(ev)

	if err != nil {
		return utils.InternalErrorResp[any](err)
	}

	err = database.WriteLog(log.ID, log.Payload, string(log.Channel), log.CreatedAt)

	if err != nil {
		return utils.InternalErrorResp[any](err)
	}

	return utils.SuccessResp[any](nil)

}

func V2RetryWaitingAttemptController(database storeInterface.IStoreProvider, id string) utils.Response[any] {

	attempt, err := database.GetAttempt(id)
	if err != nil {
		return utils.InternalErrorResp[any](err)
	}

	if attempt.ID == "" {
		return utils.NotFoundErrorResp[any](errors.New(fmt.Sprintf("Attempt (id : %s) doesn't exist", id)))
	}
	if attempt.Status != commons.WaitingStatus {
		return utils.NotFoundErrorResp[any](errors.New(fmt.Sprintf("Attempt (id : %s) are not waiting anymore", id)))
	}

	ev, err := commons.EventFromType(commons.FlushWaitingAttemptType, &attempt, nil)
	if err != nil {
		return utils.InternalErrorResp[any](err)
	}

	log, err := commons.LogFromEvent(ev)

	if err != nil {
		return utils.InternalErrorResp[any](err)
	}

	err = database.WriteLog(log.ID, log.Payload, string(log.Channel), log.CreatedAt)

	if err != nil {
		return utils.InternalErrorResp[any](err)
	}

	return utils.SuccessResp[any](nil)

}

func V2AbortWaitingAttemptController(database storeInterface.IStoreProvider, id string) utils.Response[commons.Attempt] {
	attempt, err := database.AbortAttempt(id, string(commons.AbortUser), true)

	if err != nil {
		return utils.InternalErrorResp[commons.Attempt](err)
	}

	if attempt.ID == "" {
		return utils.NotFoundErrorResp[commons.Attempt](errors.New(fmt.Sprintf("Attempt (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(attempt)

}
