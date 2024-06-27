package webhookcontroller

import (
	
	"testing"

	"github.com/formancehq/webhooks/internal/commons"
	
	v2Ctrl "github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/v2"
	"github.com/stretchr/testify/require"
)


func TestRunAttemptV2(t *testing.T){
	
	t.Run("InsertAttempt", v2InsertAttempt)
	t.Run("GetWaitingAttempts", v2GetWaitingAttempts)
	t.Run("GetAbortedAttempts", v2GetAbortedAttempts)
	t.Run("RetryWaitingAttempt", v2RetryWaitingAttempt)
	t.Run("AbortWaitingAttempt", v2AbortWaitingAttempt)

	t.Run("RetryWaitingAttempts", v2RetryWaitingAttempts)

}

func v2InsertAttempt(t *testing.T){
	params := commons.HookBodyParams{
		Name:"Test1",
		Endpoint:"http://www.exemple-endpoint.com/valide", 
		Secret:"Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh", 
		Events : []string{"event"}}

	resp := v2Ctrl.V2CreateHookController(Database, params)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)

	hook := resp.Data 

	attempt := commons.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "Attempt1")
	require.NoError(t, Database.SaveAttempt(*attempt, true))	
	attempt2 := commons.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "Attempt2")
	require.NoError(t, Database.SaveAttempt(*attempt2, true))	
	attempt3 := commons.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "Attempt3")
	require.NoError(t, Database.SaveAttempt(*attempt3, true))	
	attempt4 := commons.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "Attempt4")
	require.NoError(t, Database.SaveAttempt(*attempt4, true))	


}

func v2GetWaitingAttempts(t *testing.T){
	resp := v2Ctrl.V2GetWaitingAttemptsController(Database, "")
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 4)
}

func v2AbortWaitingAttempt(t *testing.T){
	temp := v2Ctrl.V2GetWaitingAttemptsController(Database, "")
	require.NoError(t, temp.Err)
	attempt := temp.Data.Data[0]

	resp := v2Ctrl.V2AbortWaitingAttemptController(Database, attempt.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, commons.AbortUser, resp.Data.Comment)
}

func v2GetAbortedAttempts(t *testing.T){
	resp := v2Ctrl.V2GetAbortedAttemptsController(Database, "")
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 1)
}

func v2RetryWaitingAttempt(t *testing.T){	
	temp := v2Ctrl.V2GetWaitingAttemptsController(Database, "")
	require.NoError(t, temp.Err)
	attempt := temp.Data.Data[0]
	resp := v2Ctrl.V2RetryWaitingAttemptController(Database, attempt.ID)
	require.NoError(t, resp.Err)

}

func v2RetryWaitingAttempts(t *testing.T){
	resp := v2Ctrl.V2RetryWaitingAttemptsController(Database)
	require.NoError(t, resp.Err)
}