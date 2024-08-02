package service

import (
	"errors"
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"
)

func V1CreateHook(hook utils.V1HookUser) utils.Response[utils.V1Hook] {

	if err := utils.ValidateEndpoint(hook.Endpoint); err != nil {
		return utils.ValidationErrorResp[utils.V1Hook](err)
	}

	if err := utils.ValidateSecret(&hook.Secret); err != nil {
		return utils.ValidationErrorResp[utils.V1Hook](err)
	}

	if len((&hook).EventTypes) == 0 {
		return utils.ValidationErrorResp[utils.V1Hook](errors.New("EventTypes missing"))
	}

	if err := utils.FormatEvents(&hook.EventTypes); err != nil {
		return utils.ValidationErrorResp[utils.V1Hook](err)
	}

	newHook, err := BaseCreateHook("", hook.EventTypes, hook.Endpoint, hook.Secret, true)

	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}

	newHook, err = BaseActivateHook(newHook.ID)

	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}

	return utils.SuccessResp(utils.ToV1Hook(newHook))

}

func V1GetHooks(filterEndpoint, filterId, filterCursor string, pageSize int) utils.Response[bunpaginate.Cursor[utils.V1Hook]] {
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
		hook, err := BaseGetHook(filterId)
		if err != nil {
			return utils.InternalErrorResp[bunpaginate.Cursor[utils.V1Hook]](err)
		}

		if hook.ID != "" {
			v1Hooks = append(v1Hooks, utils.ToV1Hook(hook))
		}

	} else {
		temps, hM, err := BaseGetHooks(filterEndpoint, cursor, pageSize)
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

func V1DeleteHook(id string) utils.Response[utils.V1Hook] {
	hook, err := BaseDeleteHook(id)
	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[utils.V1Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(utils.ToV1Hook(hook))
}

func V1ActiveHook(id string) utils.Response[utils.V1Hook] {
	hook, err := BaseActivateHook(id)
	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[utils.V1Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(utils.ToV1Hook(hook))
}

func V1DeactiveHook(id string) utils.Response[utils.V1Hook] {
	hook, err := BaseDeactivateHook(id)
	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[utils.V1Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(utils.ToV1Hook(hook))
}

func V1ChangeSecret(id, secret string) utils.Response[utils.V1Hook] {
	if err := utils.ValidateSecret(&secret); err != nil {
		return utils.ValidationErrorResp[utils.V1Hook](err)
	}

	hook, err := BaseUpdateSecret(id, secret)

	if err != nil {
		return utils.InternalErrorResp[utils.V1Hook](err)
	}
	if hook.ID == "" {
		return utils.NotFoundErrorResp[utils.V1Hook](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(utils.ToV1Hook(hook))

}

func V1TestHook(id, payload string) utils.Response[utils.V1Attempt] {
	hook, attempt, err := BaseTestHook(id, payload)

	if err != nil {
		return utils.InternalErrorResp[utils.V1Attempt](err)
	}

	if hook.ID == "" {
		return utils.NotFoundErrorResp[utils.V1Attempt](errors.New(fmt.Sprintf("Hook (id : %s) doesn't exist", id)))
	}

	return utils.SuccessResp(utils.ToV1Attempt(*hook, *attempt))
}
