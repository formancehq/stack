package service

import (
	"errors"
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"
	"github.com/formancehq/webhooks/internal/models"
)

func V2GetWaitingAttempts(filterCursor string, pageSize int) utils.Response[bunpaginate.Cursor[models.Attempt]] {
	hasMore := false

	strPrevious := " "
	strNext := " "

	cursor, err := utils.ReadCursor(filterCursor)

	if err != nil {
		return utils.ValidationErrorResp[bunpaginate.Cursor[models.Attempt]](err)
	}

	attempts, hM, err := getDatabase().GetWaitingAttempts(cursor, pageSize)
	if err != nil {
		return utils.InternalErrorResp[bunpaginate.Cursor[models.Attempt]](err)
	}

	hasMore = hM

	if hasMore {
		strPrevious, strNext = utils.PaginationCursor(cursor, hasMore)
	}

	Cursor := bunpaginate.Cursor[models.Attempt]{
		PageSize: pageSize,
		HasMore:  hasMore,
		Previous: strPrevious,
		Next:     strNext,
		Data:     utils.ToValues(*attempts),
	}

	return utils.SuccessResp[bunpaginate.Cursor[models.Attempt]](Cursor)
}

func V2GetAbortedAttempts(filterCursor string, pageSize int) utils.Response[bunpaginate.Cursor[models.Attempt]] {
	hasMore := false

	strPrevious := " "
	strNext := " "

	cursor, err := utils.ReadCursor(filterCursor)

	if err != nil {
		return utils.ValidationErrorResp[bunpaginate.Cursor[models.Attempt]](err)
	}

	attempts, hM, err := getDatabase().GetAbortedAttempts(cursor, pageSize)
	if err != nil {
		return utils.InternalErrorResp[bunpaginate.Cursor[models.Attempt]](err)
	}

	hasMore = hM

	if hasMore {
		strPrevious, strNext = utils.PaginationCursor(cursor, hasMore)
	}

	Cursor := bunpaginate.Cursor[models.Attempt]{
		PageSize: pageSize,
		HasMore:  hasMore,
		Previous: strPrevious,
		Next:     strNext,
		Data:     utils.ToValues(*attempts),
	}

	return utils.SuccessResp[bunpaginate.Cursor[models.Attempt]](Cursor)
}

func V2RetryWaitingAttempts() utils.Response[any] {

	err := getDatabase().FlushAttempts("")
	
	if err != nil {
		return utils.InternalErrorResp[any](err)
	}

	return utils.SuccessResp[any](nil)

}

func V2RetryWaitingAttempt(id string) utils.Response[any] {

	attempt, err := getDatabase().GetAttempt(id)
	if err != nil {
		return utils.InternalErrorResp[any](err)
	}

	if attempt.ID == "" {
		return utils.NotFoundErrorResp[any](errors.New(fmt.Sprintf("Attempt (id : %s) doesn't exist", id)))
	}
	if attempt.Status != models.WaitingStatus {
		return utils.NotFoundErrorResp[any](errors.New(fmt.Sprintf("Attempt (id : %s) are not waiting anymore", id)))
	}

	err = getDatabase().FlushAttempts(attempt.ID)

	if err != nil {
		return utils.InternalErrorResp[any](err)
	}

	return utils.SuccessResp[any](nil)

}

func V2AbortWaitingAttempt(id string) utils.Response[models.Attempt] {

	attempt, err := getDatabase().AbortAttempt(id, string(models.AbortUser), true)

	if err != nil {
		return utils.InternalErrorResp[models.Attempt](err)
	}

	if attempt.ID == "" {
		return utils.NotFoundErrorResp[models.Attempt](errors.New(fmt.Sprintf("Attempt (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(attempt)

}
