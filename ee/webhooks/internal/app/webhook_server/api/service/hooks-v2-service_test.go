package service

import (
	"errors"
	"testing"

	utilsCtrl "github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"
	"github.com/formancehq/webhooks/internal/models"
	"github.com/stretchr/testify/require"
)

func TestRunHookV2(t *testing.T) {
	//Reset Hooks
	resp := V2GetHooks("", "", 15)
	for _, hook := range resp.Data.Data {
		r := V2DeleteHook(hook.ID)
		require.NoError(t, r.Err)
	}
	t.Run("InsertHook", v2InsertHook)

	t.Run("GetHooks", v2GetHooks)

	t.Run("GetHook", v2GetHook)

	t.Run("DeleteHook", v2DeleteHook)

	t.Run("DeactiveHook", v2DeactiveHook)

	t.Run("ActiveHook", v2ActiveHook)

	t.Run("ChangeSecret", v2ChangeSecret)

	t.Run("ChangeEndpoint", v2ChangeEndpoint)
}

func v2InsertHook(t *testing.T) {
	badHook1 := models.HookBodyParams{
		Name:     "Test1",
		Endpoint: "",
		Secret:   "",
		Events:   []string{"event1"}}

	resp := V2CreateHook(badHook1)
	require.Error(t, resp.Err, "Validation error expected for bad endpoint")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad endpoint")

	badHook2 := models.HookBodyParams{
		Name:     "Test1",
		Endpoint: "http://www.exemple-endpoint.com/valide",
		Secret:   "badsecret!",
		Events:   []string{"event1"}}

	resp = V2CreateHook(badHook2)
	require.Error(t, resp.Err, "Validation error expected for bad secret")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad secret")

	badHook3 := models.HookBodyParams{
		Name:     "Test1",
		Endpoint: "http://www.exemple-endpoint.com/valide",
		Secret:   "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh",
		Events:   []string{""}}

	resp = V2CreateHook(badHook3)
	require.Error(t, resp.Err, "Validation error expected for bad events")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad events")

	hook1 := models.HookBodyParams{
		Name:     "Test1",
		Endpoint: "http://www.exemple-endpoint.com/valide",
		Secret:   "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh",
		Events:   []string{"event"}}

	resp = V2CreateHook(hook1)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)
	require.Equal(t, resp.Data.Endpoint, "http://www.exemple-endpoint.com/valide")

	hook2 := models.HookBodyParams{
		Name:     "Test2",
		Endpoint: "http://www.exemple-endpoint.com/valide",
		Secret:   "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh",
		Events:   []string{"event"}}

	resp = V2CreateHook(hook2)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)
	require.Equal(t, resp.Data.Endpoint, "http://www.exemple-endpoint.com/valide")

	hook3 := models.HookBodyParams{
		Name:     "Test3",
		Endpoint: "http://www.exemple-endpoint.com/valide",
		Secret:   "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh",
		Events:   []string{"event"}}

	resp = V2CreateHook(hook3)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)
	require.Equal(t, resp.Data.Endpoint, "http://www.exemple-endpoint.com/valide")

	hook4 := models.HookBodyParams{
		Name:     "Test4",
		Endpoint: "http://www.exemple-endpoint.com/valide2",
		Secret:   "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh",
		Events:   []string{"event"}}

	resp = V2CreateHook(hook4)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)
	require.Equal(t, resp.Data.Endpoint, "http://www.exemple-endpoint.com/valide2")
}

func v2GetHooks(t *testing.T) {
	badCursor := "bad"

	goodEndpoint := "http://www.exemple-endpoint.com/valide"

	resp := V2GetHooks("", badCursor, 15)
	require.Error(t, resp.Err, "Validation error expected for bad cursor")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad cursor")

	resp = V2GetHooks("", "", 15)
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 4)

	resp = V2GetHooks(goodEndpoint, "", 15)
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 3)
}

func v2GetHook(t *testing.T) {
	resp := V2GetHooks("", "", 15)
	require.NoError(t, resp.Err)
	hook := resp.Data.Data[0]

	resp2 := V2GetHook(hook.ID)
	require.NoError(t, resp2.Err)
	require.Equal(t, hook.ID, resp2.Data.ID)
}

func v2DeleteHook(t *testing.T) {
	wrongId := "23"

	resp := V2DeleteHook(wrongId)
	require.Error(t, resp.Err)
	require.Equal(t, utilsCtrl.NotFoundType, resp.T, "NotFound type error expected for bad idea")

	temp := V2GetHooks("", "", 15)
	hook := temp.Data.Data[0]
	resp = V2DeleteHook(hook.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, false, resp.Data.Active)
	temp = V2GetHooks("", "", 15)
	require.NoError(t, temp.Err)
	require.Len(t, temp.Data.Data, 3)
}

func v2DeactiveHook(t *testing.T) {
	wrongId := "23"

	resp := V2DeactiveHook(wrongId)
	require.Error(t, resp.Err)
	require.Equal(t, utilsCtrl.NotFoundType, resp.T, "NotFound type error expected for bad idea")

	temp := V2GetHooks("", "", 15)
	hook := temp.Data.Data[0]
	resp = V2DeactiveHook(hook.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, false, resp.Data.Active)

}

func v2ActiveHook(t *testing.T) {
	wrongId := "23"

	resp := V2ActiveHook(wrongId)
	require.Error(t, resp.Err)
	require.Equal(t, utilsCtrl.NotFoundType, resp.T, "NotFound type error expected for bad idea")

	var inactiveHook models.Hook
	temp := V2GetHooks("", "", 15)
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

	resp = V2ActiveHook(inactiveHook.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, true, resp.Data.Active)

}

func v2ChangeSecret(t *testing.T) {
	badSecret := "badsecret!"
	temp := V2GetHooks("", "", 15)
	hook := temp.Data.Data[0]
	resp := V2ChangeSecret(hook.ID, badSecret)
	require.Error(t, resp.Err, "Validation type error required for bad secret")

	resp = V2ChangeSecret(hook.ID, "")
	require.NoError(t, resp.Err)
	require.NotEqual(t, hook.Secret, resp.Data.Secret)
}

func v2ChangeEndpoint(t *testing.T) {
	badEndpoint := ""
	newEndpoint := "http://www.exemple-endpoint.com/newvalide"
	temp := V2GetHooks("", "", 15)
	hook := temp.Data.Data[0]
	resp := V2ChangeEndpoint(hook.ID, badEndpoint)
	require.Error(t, resp.Err, "Validation type error required for bad endpoint")

	resp = V2ChangeEndpoint(hook.ID, newEndpoint)
	require.NoError(t, resp.Err)
	require.Equal(t, newEndpoint, resp.Data.Endpoint)

}
