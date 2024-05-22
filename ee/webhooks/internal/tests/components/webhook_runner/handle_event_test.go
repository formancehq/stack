package webhookrunner

import (
	_ "errors"
	"testing"
	"time"

	"github.com/formancehq/webhooks/internal/commons"
	_ "github.com/formancehq/webhooks/internal/commons"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
)



func TestRunHandleEvent(t *testing.T){

	t.Run("NewHookEvent", newHookEvent)

	t.Run("ChangeHookStatusEvent", changeHookStatusEvent)

	t.Run("ChangeHookSecretEvent", changeHookSecretEvent)

	t.Run("ChangeHookEndpointEvent", changeHookEndpointEvent)

	t.Run("ChangeHookRetryEvent", changeHookRetryEvent)

	t.Run("NewWaitingAttemptEvent", newWaitingAttemptEvent)

	t.Run("FlushWaitingAttemptsEvent", flushWaitingAttemptsEvent)

	t.Run("FlushWaitingAttemptEvent", flushWaitingAttemptEvent)

	t.Run("AbortWaitingAttemptEvent", abortWaitingAttemptEvent)
}



func newHookEvent(t *testing.T){
	newHook := commons.NewHook("Test", []string{"testevent"}, "/localhost/", "", false)
	err := Database.SaveHook(*newHook)
	require.NoError(t, err)
	newHookEvent, err := commons.EventFromType(commons.NewHookType, nil, newHook)
	require.NoError(t, err) 
	WebhookRunner.HandleEvent(newHookEvent)
	saveHook := WebhookRunner.State.HooksById.Get(newHook.ID)
	require.NotNil(t, saveHook)
}

func changeHookStatusEvent(t *testing.T){
	hooks, _, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]
	require.Equal(t, commons.DisableStatus, hook.Status)
	hook.Status = commons.EnableStatus
	newEvent, err := commons.EventFromType(commons.ChangeHookStatusType,nil, hook)
	WebhookRunner.HandleEvent(newEvent)
	saveHook := WebhookRunner.State.HooksById.Get(hook.ID)
	require.Equal(t, commons.EnableStatus, saveHook.Val.Status)

	hook.Status = commons.DisableStatus
	newEvent, err = commons.EventFromType(commons.ChangeHookStatusType,nil, hook)
	WebhookRunner.HandleEvent(newEvent)
	saveHook = WebhookRunner.State.HooksById.Get(hook.ID)
	require.Equal(t, commons.DisableStatus, saveHook.Val.Status)

	hook.Status = commons.DeleteStatus
	newEvent, err = commons.EventFromType(commons.ChangeHookStatusType,nil, hook)
	WebhookRunner.HandleEvent(newEvent)
	saveHook = WebhookRunner.State.HooksById.Get(hook.ID)
	require.Nil(t, saveHook)

	Database.ActivateHook(hook.ID)
	newHookEvent, err := commons.EventFromType(commons.NewHookType, nil, hook)
	require.NoError(t, err) 
	WebhookRunner.HandleEvent(newHookEvent)
}

func changeHookSecretEvent(t *testing.T){
	hooks, _, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]
	newSecret := "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh"
	hook.Secret = newSecret
	newEvent, err := commons.EventFromType(commons.ChangeHookSecretType, nil, hook)
	WebhookRunner.HandleEvent(newEvent)
	saveHook := WebhookRunner.State.HooksById.Get(hook.ID)
	require.Equal(t, newSecret, saveHook.Val.Secret)

}

func changeHookEndpointEvent(t *testing.T){
	hooks, _, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]
	newEndpoint := "www.google.fr/top"
	hook.Endpoint = newEndpoint
	newEvent, err := commons.EventFromType(commons.ChangeHookEndpointType, nil, hook)
	WebhookRunner.HandleEvent(newEvent)
	saveHook := WebhookRunner.State.HooksById.Get(hook.ID)
	require.Equal(t, newEndpoint, saveHook.Val.Endpoint)
}

func changeHookRetryEvent(t *testing.T){
	hooks, _, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]
	require.False(t, hook.Retry)
	hook.Retry=true
	newEvent, err := commons.EventFromType(commons.ChangeHookRetryType, nil, hook)
	WebhookRunner.HandleEvent(newEvent)
	saveHook := WebhookRunner.State.HooksById.Get(hook.ID)
	require.True(t, saveHook.Val.Retry)
}

func newWaitingAttemptEvent(t *testing.T){
	hooks, _, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]

	attempt := commons.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "blabla")
	Database.SaveAttempt(*attempt)

	newEvent, err := commons.EventFromType(commons.NewWaitingAttemptType, attempt, nil)
	WebhookRunner.HandleEvent(newEvent)
	saveAttempt := WebhookRunner.State.WaitingAttempts.FindElement(func(s *commons.SharedAttempt) bool {
		return attempt.ID == s.Val.ID
	})

	require.NotNil(t, saveAttempt)


}

func flushWaitingAttemptEvent(t *testing.T){
	now := time.Now()
	attempt := (*WebhookRunner.State.WaitingAttempts.Val)[0]
	attempt.Val.NextTry = now.Add(5*time.Hour)
	newEvent, err := commons.EventFromType(commons.FlushWaitingAttemptType, attempt.Val, nil)
	require.NoError(t, err)
	WebhookRunner.HandleEvent(newEvent)
	attempt = (*WebhookRunner.State.WaitingAttempts.Val)[0]
	require.True(t, attempt.Val.NextTry.Before(now.Add(5*time.Hour)))


}

func flushWaitingAttemptsEvent( t *testing.T){
	now := time.Now()
	attempt := (*WebhookRunner.State.WaitingAttempts.Val)[0]
	attempt.Val.NextTry = now.Add(5*time.Hour)
	newEvent, err := commons.EventFromType(commons.FlushWaitingAttemptsType, attempt.Val, nil)
	require.NoError(t, err)
	WebhookRunner.HandleEvent(newEvent)

	attempt = (*WebhookRunner.State.WaitingAttempts.Val)[0]
	require.True(t, attempt.Val.NextTry.Before(now.Add(5*time.Hour)))
}

func abortWaitingAttemptEvent( t *testing.T){
	attempts, _, _ := Database.GetWaitingAttempts(0,1)
	require.Len(t, *attempts, 1)
	newEvent, _ := commons.EventFromType(commons.AbortWaitingAttemptType, (*attempts)[0], nil)
	WebhookRunner.HandleEvent(newEvent)

	
	require.Len(t, *WebhookRunner.State.WaitingAttempts.Val, 0)

}