package webhookcontroller

import (
	"errors"
	"testing"

	"github.com/formancehq/webhooks/internal/commons"
	utilsCtrl "github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/utils"
	v2Ctrl "github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/v2"
	"github.com/stretchr/testify/require"
)


func TestRunHookV2(t *testing.T){
	//Reset Hooks
	resp := v2Ctrl.V2GetHooksController(Database, "", "")
	for _, hook := range resp.Data.Data {
		r := v2Ctrl.V2DeleteHookController(Database, hook.ID)
		require.NoError(t,r.Err)
	}
	t.Run("InsertHook", v2InsertHook)

	t.Run("GetHooks", v2GetHooks)

	t.Run("DeleteHook", v2DeleteHook)

	t.Run("DeactiveHook", v2DeactiveHook)

	t.Run("ActiveHook", v2ActiveHook)

	t.Run("ChangeSecret", v2ChangeSecret)

	t.Run("ChangeEndpoint", v2ChangeEndpoint)
}


func v2InsertHook(t *testing.T){
	badHook1 := commons.HookBodyParams{
		Name:"Test1",
		Endpoint:"", 
		Secret:"", 
		Events : []string{"event1"}}

	resp := v2Ctrl.V2CreateHookController(Database, badHook1)
	require.Error(t, resp.Err, "Validation error expected for bad endpoint")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad endpoint")
	
	badHook2 := commons.HookBodyParams{
		Name:"Test1",
		Endpoint:"http://www.exemple-endpoint.com/valide", 
		Secret:"badsecret!", 
		Events : []string{"event1"}}

	resp = v2Ctrl.V2CreateHookController(Database, badHook2)
	require.Error(t, resp.Err, "Validation error expected for bad secret")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad secret")

	badHook3 := commons.HookBodyParams{
		Name:"Test1",
		Endpoint:"http://www.exemple-endpoint.com/valide", 
		Secret:"Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh", 
		Events : []string{""}}

	resp = v2Ctrl.V2CreateHookController(Database, badHook3)
	require.Error(t, resp.Err, "Validation error expected for bad events")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad events")


	hook1 := commons.HookBodyParams{
		Name:"Test1",
		Endpoint:"http://www.exemple-endpoint.com/valide", 
		Secret:"Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh", 
		Events : []string{"event"}}

	resp = v2Ctrl.V2CreateHookController(Database, hook1)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)
	require.Equal(t, resp.Data.Endpoint, "http://www.exemple-endpoint.com/valide")

	hook2 := commons.HookBodyParams{
		Name:"Test2",
		Endpoint:"http://www.exemple-endpoint.com/valide", 
		Secret:"Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh", 
		Events : []string{"event"}}

	resp = v2Ctrl.V2CreateHookController(Database, hook2)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)
	require.Equal(t, resp.Data.Endpoint, "http://www.exemple-endpoint.com/valide")

	hook3 := commons.HookBodyParams{
		Name:"Test3",
		Endpoint:"http://www.exemple-endpoint.com/valide", 
		Secret:"Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh", 
		Events : []string{"event"}}

	resp = v2Ctrl.V2CreateHookController(Database, hook3)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)
	require.Equal(t, resp.Data.Endpoint, "http://www.exemple-endpoint.com/valide")
	
	hook4 := commons.HookBodyParams{
		Name:"Test4",
		Endpoint:"http://www.exemple-endpoint.com/valide2", 
		Secret:"Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh", 
		Events : []string{"event"}}

	resp = v2Ctrl.V2CreateHookController(Database, hook4)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)
	require.Equal(t, resp.Data.Endpoint, "http://www.exemple-endpoint.com/valide2")
}

func v2GetHooks(t *testing.T){
	badCursor := "bad"	

	goodEndpoint := "http://www.exemple-endpoint.com/valide"

	resp := v2Ctrl.V2GetHooksController(Database, "", badCursor)
	require.Error(t, resp.Err, "Validation error expected for bad cursor")
	require.Equal(t, resp.T, utilsCtrl.ValidationType, "Validation type error expected for bad cursor")

	resp = v2Ctrl.V2GetHooksController(Database, "", "") 
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 4)

	resp = v2Ctrl.V2GetHooksController(Database, goodEndpoint, "")
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 3)
}

func v2DeleteHook(t *testing.T){
	wrongId := "23"

	resp := v2Ctrl.V2DeleteHookController(Database, wrongId)
	require.Error(t, resp.Err)
	require.Equal(t, utilsCtrl.NotFoundType, resp.T, "NotFound type error expected for bad idea")

	temp := v2Ctrl.V2GetHooksController(Database, "", "")
	hook := temp.Data.Data[0]
	resp = v2Ctrl.V2DeleteHookController(Database, hook.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, false, resp.Data.Active)
	temp = v2Ctrl.V2GetHooksController(Database, "", "")
	require.NoError(t, temp.Err)
	require.Len(t, temp.Data.Data, 3)
}

func v2DeactiveHook(t *testing.T){
	wrongId := "23"

	resp := v2Ctrl.V2DeactiveHookController(Database, wrongId)
	require.Error(t, resp.Err)
	require.Equal(t, utilsCtrl.NotFoundType, resp.T, "NotFound type error expected for bad idea")

	temp := v2Ctrl.V2GetHooksController(Database, "", "")
	hook := temp.Data.Data[0]
	resp = v2Ctrl.V2DeactiveHookController(Database, hook.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, false, resp.Data.Active)
	
}

func v2ActiveHook(t *testing.T){
	wrongId := "23"

	resp := v2Ctrl.V2ActiveHookController(Database, wrongId)
	require.Error(t, resp.Err)
	require.Equal(t, utilsCtrl.NotFoundType, resp.T, "NotFound type error expected for bad idea")

	var inactiveHook commons.Hook
	temp := v2Ctrl.V2GetHooksController(Database, "", "")
	for _,h := range temp.Data.Data {
		if !h.Active {inactiveHook = h; return;}
	}

	if inactiveHook.ID == "" {
		require.NoError(t, errors.New("Inactive hook is missing"))
		return
	}

	resp = v2Ctrl.V2ActiveHookController(Database, inactiveHook.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, true, resp.Data.Active)

}

func v2ChangeSecret(t *testing.T){
	badSecret := "badsecret!"
	temp := v2Ctrl.V2GetHooksController(Database, "", "")
	hook := temp.Data.Data[0]
	resp := v2Ctrl.V2ChangeSecretController(Database, hook.ID, badSecret)
	require.Error(t, resp.Err, "Validation type error required for bad secret")

	resp = v2Ctrl.V2ChangeSecretController(Database, hook.ID, "")
	require.NoError(t, resp.Err)
	require.NotEqual(t, hook.Secret, resp.Data.Secret)
}

func v2ChangeEndpoint(t *testing.T){
	badEndpoint := ""
	newEndpoint := "http://www.exemple-endpoint.com/newvalide"
	temp := v2Ctrl.V2GetHooksController(Database, "", "")
	hook := temp.Data.Data[0]
	resp := v2Ctrl.V2ChangeEndpointController(Database, hook.ID, badEndpoint)
	require.Error(t, resp.Err, "Validation type error required for bad endpoint")

	resp = v2Ctrl.V2ChangeEndpointController(Database, hook.ID, newEndpoint)
	require.NoError(t, resp.Err)
	require.Equal(t, newEndpoint, resp.Data.Endpoint)

}

func V2TestHook(t *testing.T){
	//TODO(Test)
}

