package webhookcontroller

import (
	"errors"
	"testing"

	utilsCtrl "github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/utils"
	v1Ctrl "github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/v1"
	"github.com/stretchr/testify/require"
)

func TestRunV1(t *testing.T) {
	//Reset Hooks
	resp := v1Ctrl.V1GetHooksController(Database, "", "", "")
	for _, hook := range resp.Data.Data {
		r := v1Ctrl.V1DeleteHookController(Database, hook.ID)
		require.NoError(t, r.Err)
	}

	t.Run("InsertHook", v1InsertHook)

	t.Run("GetHooks", v1GetHooks)

	t.Run("DeleteHook", v1DeleteHook)

	t.Run("DeactiveHook", v1DeactiveHook)

	t.Run("ActiveHook", v1ActiveHook)

	t.Run("ChangeSecret", v1ChangeSecret)
}

func v1InsertHook(t *testing.T) {
	badHook1 := utilsCtrl.V1HookUser{
		Endpoint:   "",
		Secret:     "",
		EventTypes: []string{"event1"}}

	resp := v1Ctrl.V1CreateHookController(Database, badHook1)
	require.Error(t, resp.Err, "Validation error expected for bad endpoint")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad endpoint")

	badHook2 := utilsCtrl.V1HookUser{
		Endpoint:   "http://www.exemple-endpoint.com/valide",
		Secret:     "badsecret!",
		EventTypes: []string{"event1"}}

	resp = v1Ctrl.V1CreateHookController(Database, badHook2)
	require.Error(t, resp.Err, "Validation error expected for bad secret")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad secret")

	badHook3 := utilsCtrl.V1HookUser{
		Endpoint:   "http://www.exemple-endpoint.com/valide",
		Secret:     "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh",
		EventTypes: []string{""}}

	resp = v1Ctrl.V1CreateHookController(Database, badHook3)
	require.Error(t, resp.Err, "Validation error expected for bad events")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad events")

	hook1 := utilsCtrl.V1HookUser{
		Endpoint:   "http://www.exemple-endpoint.com/valide",
		Secret:     "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh",
		EventTypes: []string{"event"}}

	resp = v1Ctrl.V1CreateHookController(Database, hook1)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)
	require.Equal(t, resp.Data.Endpoint, "http://www.exemple-endpoint.com/valide")

	hook2 := utilsCtrl.V1HookUser{
		Endpoint:   "http://www.exemple-endpoint.com/valide",
		Secret:     "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh",
		EventTypes: []string{"event"}}

	resp = v1Ctrl.V1CreateHookController(Database, hook2)
	require.NoError(t, resp.Err)

	hook3 := utilsCtrl.V1HookUser{
		Endpoint:   "http://www.exemple-endpoint.com/valide2",
		Secret:     "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh",
		EventTypes: []string{"event"}}

	resp = v1Ctrl.V1CreateHookController(Database, hook3)
	require.NoError(t, resp.Err)

}

func v1GetHooks(t *testing.T) {
	badCursor := "bad"
	wrongId := "23"

	goodEndpoint := "http://www.exemple-endpoint.com/valide"

	resp := v1Ctrl.V1GetHooksController(Database, "", "", badCursor)
	require.Error(t, resp.Err, "Validation error expected for bad cursor")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad cursor")

	resp = v1Ctrl.V1GetHooksController(Database, "", wrongId, "")
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 0)

	resp = v1Ctrl.V1GetHooksController(Database, "", "", "")
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 3)

	resp = v1Ctrl.V1GetHooksController(Database, goodEndpoint, "", "")
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 2)
}

func v1DeleteHook(t *testing.T) {
	wrongId := "23"

	resp := v1Ctrl.V1DeleteHookController(Database, wrongId)
	require.Error(t, resp.Err)
	require.Equal(t, utilsCtrl.NotFoundType, resp.T, "NotFound type error expected for bad idea")

	temp := v1Ctrl.V1GetHooksController(Database, "", "", "")
	hook := temp.Data.Data[0]
	resp = v1Ctrl.V1DeleteHookController(Database, hook.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, false, resp.Data.Active)
	temp = v1Ctrl.V1GetHooksController(Database, "", "", "")
	require.NoError(t, temp.Err)
	require.Len(t, temp.Data.Data, 2)
}

func v1DeactiveHook(t *testing.T) {
	wrongId := "23"

	resp := v1Ctrl.V1DeactiveHookController(Database, wrongId)
	require.Error(t, resp.Err)
	require.Equal(t, utilsCtrl.NotFoundType, resp.T, "NotFound type error expected for bad idea")

	temp := v1Ctrl.V1GetHooksController(Database, "", "", "")
	hook := temp.Data.Data[0]
	resp = v1Ctrl.V1DeactiveHookController(Database, hook.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, false, resp.Data.Active)

}

func v1ActiveHook(t *testing.T) {
	wrongId := "23"

	resp := v1Ctrl.V1ActiveHookController(Database, wrongId)
	require.Error(t, resp.Err)
	require.Equal(t, utilsCtrl.NotFoundType, resp.T, "NotFound type error expected for bad idea")

	var inactiveHook utilsCtrl.V1Hook
	temp := v1Ctrl.V1GetHooksController(Database, "", "", "")
	for _, h := range temp.Data.Data {
		if !h.Active {
			inactiveHook = h
			return
		}
	}

	if inactiveHook.ID == "" {
		require.NoError(t, errors.New("Inactive hook is missing"))
		return
	}

	resp = v1Ctrl.V1ActiveHookController(Database, inactiveHook.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, true, resp.Data.Active)

}

func v1ChangeSecret(t *testing.T) {
	badSecret := "badsecret!"
	temp := v1Ctrl.V1GetHooksController(Database, "", "", "")
	hook := temp.Data.Data[0]
	resp := v1Ctrl.V1ChangeSecretController(Database, hook.ID, badSecret)
	require.Error(t, resp.Err, "Validation type error required for bad secret")

	resp = v1Ctrl.V1ChangeSecretController(Database, hook.ID, "")
	require.NoError(t, resp.Err)
	require.NotEqual(t, hook.Secret, resp.Data.Secret)
}
