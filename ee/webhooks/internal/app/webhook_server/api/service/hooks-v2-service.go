package service

import (
	"errors"
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"
	"github.com/formancehq/webhooks/internal/models"
)

func V2CreateHook(hookParams models.HookBodyParams) utils.Response[models.Hook] {

	if err := utils.ValidateEndpoint(hookParams.Endpoint); err != nil {
		return utils.ValidationErrorResp[models.Hook](err)
	}

	if err := utils.ValidateSecret(&hookParams.Secret); err != nil {
		return utils.ValidationErrorResp[models.Hook](err)
	}

	if len((&hookParams).Events) == 0 {
		return utils.ValidationErrorResp[models.Hook](errors.New("Events missing"))
	}

	if err := utils.FormatEvents(&hookParams.Events); err != nil {
		return utils.ValidationErrorResp[models.Hook](err)
	}

	hook, err := BaseCreateHook(hookParams.Name, hookParams.Events, hookParams.Endpoint, hookParams.Secret, hookParams.Retry)

	if err != nil {
		return utils.InternalErrorResp[models.Hook](err)
	}

	return utils.SuccessResp(hook)
}

func V2GetHook(id string) utils.Response[models.Hook] {

	hook, err := BaseGetHook(id)
	if err != nil {
		return utils.InternalErrorResp[models.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[models.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}
	return utils.SuccessResp(hook)
}

func V2GetHooks(filterEndpoint, filterCursor string, pageSize int) utils.Response[bunpaginate.Cursor[models.Hook]] {
	hasMore := false
	strPrevious := " "
	strNext := " "

	cursor, err := utils.ReadCursor(filterCursor)

	if err != nil {
		return utils.ValidationErrorResp[bunpaginate.Cursor[models.Hook]](err)
	}

	if filterEndpoint != "" {
		if err := utils.ValidateEndpoint(filterEndpoint); err != nil {
			return utils.ValidationErrorResp[bunpaginate.Cursor[models.Hook]](err)
		}
	}

	hooks, hM, err := BaseGetHooks(filterEndpoint, cursor, pageSize)
	if err != nil {
		return utils.InternalErrorResp[bunpaginate.Cursor[models.Hook]](err)
	}
	hasMore = hM

	if hasMore {
		strPrevious, strNext = utils.PaginationCursor(cursor, hasMore)
	}

	Cursor := bunpaginate.Cursor[models.Hook]{
		PageSize: pageSize,
		HasMore:  hasMore,
		Previous: strPrevious,
		Next:     strNext,
		Data:     utils.ToValues(*hooks),
	}

	return utils.SuccessResp(Cursor)
}

func V2DeleteHook(id string) utils.Response[models.Hook] {
	hook, err := BaseDeleteHook(id)
	if err != nil {
		return utils.InternalErrorResp[models.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[models.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}

func V2ActiveHook(id string) utils.Response[models.Hook] {
	hook, err := BaseActivateHook(id)
	if err != nil {
		return utils.InternalErrorResp[models.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[models.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}

func V2DeactiveHook(id string) utils.Response[models.Hook] {
	hook, err := BaseDeactivateHook(id)
	if err != nil {
		return utils.InternalErrorResp[models.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[models.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}

func V2ChangeSecret(id, secret string) utils.Response[models.Hook] {
	if err := utils.ValidateSecret(&secret); err != nil {
		return utils.ValidationErrorResp[models.Hook](err)
	}

	hook, err := BaseUpdateSecret(id, secret)

	if err != nil {
		return utils.InternalErrorResp[models.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[models.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}

func V2TestHook(id, payload string) utils.Response[models.Attempt] {
	hook, attempt, err := BaseTestHook(id, payload)

	if err != nil {
		return utils.InternalErrorResp[models.Attempt](err)
	}

	if hook.ID == "" {
		return utils.NotFoundErrorResp[models.Attempt](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(*attempt)
}

func V2ChangeEndpoint(id, endpoint string) utils.Response[models.Hook] {
	if err := utils.ValidateEndpoint(endpoint); err != nil {
		return utils.ValidationErrorResp[models.Hook](err)
	}

	hook, err := BaseUpdateEndpoint(id, endpoint)

	if err != nil {
		return utils.InternalErrorResp[models.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[models.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}

func V2ChangeRetry(id string, retry bool) utils.Response[models.Hook] {

	hook, err := BaseUpdateRetry(id, retry)

	if err != nil {
		return utils.InternalErrorResp[models.Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[models.Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(hook)
}
