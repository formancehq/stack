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



func TestRunHandleLogs(t *testing.T){

	t.Run("NewHookLog", newHookLog)

	t.Run("ChangeHookStatusLog", changeHookStatusLog)

	t.Run("ChangeHookSecretLog", changeHookSecretLog)

	t.Run("ChangeHookEndpointLog", changeHookEndpointLog)

	t.Run("ChangeHookRetryLog", changeHookRetryLog)

	t.Run("NewWaitingAttemptLog", newWaitingAttemptLog)

	t.Run("FlushWaitingAttemptsLog", flushWaitingAttemptsLog)

	t.Run("FlushWaitingAttemptLog", flushWaitingAttemptLog)

	t.Run("AbortWaitingAttemptLog", abortWaitingAttemptLog)
}



func newHookLog(t *testing.T){
	var start_time time.Time = time.Now()

	newHook := commons.NewHook("Test", []string{"testLog"}, "/localhost/", "", false)
	err := Database.SaveHook(*newHook)
	require.NoError(t, err)

	logs, err := Database.GetFreshLogs([]commons.Channel{commons.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	NewHookLog := (*logs)[0]

	WebhookRunner.HandleFreshLog(NewHookLog)

	saveHook := WebhookRunner.State.HooksById.Get(newHook.ID)
	require.NotNil(t, saveHook)
}

func changeHookStatusLog(t *testing.T){
	var start_time time.Time = time.Now()
	hooks, _, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]

	require.Equal(t, commons.DisableStatus, hook.Status)
	hookIndex := hook.ID


	_, err = Database.ActivateHook(hookIndex)
	require.NoError(t, err)

	logs, err := Database.GetFreshLogs([]commons.Channel{commons.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookStatusLog := (*logs)[0]

	WebhookRunner.HandleFreshLog(ChangeHookStatusLog)
	saveHook := WebhookRunner.State.HooksById.Get(hook.ID)
	require.Equal(t, commons.EnableStatus, saveHook.Val.Status)

	start_time = time.Now()
	_, err = Database.DeactivateHook(hookIndex)
	require.NoError(t, err)

	logs, err = Database.GetFreshLogs([]commons.Channel{commons.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookStatusLog = (*logs)[0]
	WebhookRunner.HandleFreshLog(ChangeHookStatusLog)
	saveHook = WebhookRunner.State.HooksById.Get(hook.ID)
	require.Equal(t, commons.DisableStatus, saveHook.Val.Status)

	start_time = time.Now()
	_, err = Database.ActivateHook(hookIndex)
	require.NoError(t, err)

	logs, err = Database.GetFreshLogs([]commons.Channel{commons.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookStatusLog = (*logs)[0]
	WebhookRunner.HandleFreshLog(ChangeHookStatusLog)
	saveHook = WebhookRunner.State.HooksById.Get(hook.ID)
	require.Equal(t, commons.EnableStatus, saveHook.Val.Status)


	start_time = time.Now()
	_, err = Database.DeleteHook(hookIndex)
	require.NoError(t, err)

	logs, err = Database.GetFreshLogs([]commons.Channel{commons.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookStatusLog = (*logs)[0]
	WebhookRunner.HandleFreshLog(ChangeHookStatusLog)
	saveHook = WebhookRunner.State.HooksById.Get(hook.ID)
	require.Nil(t, saveHook)


}

func changeHookSecretLog(t *testing.T){
	var start_time time.Time = time.Now()
	
	newHook := commons.NewHook("Test", []string{"testLog"}, "/localhost/", "", false)
	err := Database.SaveHook(*newHook)
	require.NoError(t, err)

	logs, err := Database.GetFreshLogs([]commons.Channel{commons.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	NewHookLog := (*logs)[0]

	WebhookRunner.HandleFreshLog(NewHookLog)



	hooks, _, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]
	newSecret := "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh"
	start_time  = time.Now()
	_, err = Database.UpdateHookSecret(hook.ID, newSecret)
	require.NoError(t, err)

	logs, err = Database.GetFreshLogs([]commons.Channel{commons.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookSecretLog := (*logs)[0]

	WebhookRunner.HandleFreshLog(ChangeHookSecretLog)
	saveHook := WebhookRunner.State.HooksById.Get(hook.ID)
	require.Equal(t, newSecret, saveHook.Val.Secret)

}

func changeHookEndpointLog(t *testing.T){
	hooks, _, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]
	newEndpoint := "www.google.fr/top"
	
	var start_time time.Time = time.Now()
	_, err = Database.UpdateHookEndpoint(hook.ID, newEndpoint)
	require.NoError(t, err)

	logs, err := Database.GetFreshLogs([]commons.Channel{commons.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookSecretLog := (*logs)[0]

	WebhookRunner.HandleFreshLog(ChangeHookSecretLog)
	
	saveHook := WebhookRunner.State.HooksById.Get(hook.ID)
	require.Equal(t, newEndpoint, saveHook.Val.Endpoint)
}

func changeHookRetryLog(t *testing.T){
	hooks, _, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]
	require.False(t, hook.Retry)

	
	var start_time time.Time = time.Now()
	_, err = Database.UpdateHookRetry(hook.ID, true)
	require.NoError(t, err)

	logs, err := Database.GetFreshLogs([]commons.Channel{commons.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookRetryLog := (*logs)[0]

	WebhookRunner.HandleFreshLog(ChangeHookRetryLog)

	saveHook := WebhookRunner.State.HooksById.Get(hook.ID)
	require.True(t, saveHook.Val.Retry)
}

func newWaitingAttemptLog(t *testing.T){
	hooks, _, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]

	attempt := commons.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "blabla")

	var start_time time.Time = time.Now()
	err = Database.SaveAttempt(*attempt, true)
	require.NoError(t, err)

	logs, err := Database.GetFreshLogs([]commons.Channel{commons.AttemptChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	NewWaitingAttemptLog := (*logs)[0]

	WebhookRunner.HandleFreshLog(NewWaitingAttemptLog)

	saveAttempt := WebhookRunner.State.WaitingAttempts.FindElement(func(s *commons.SharedAttempt) bool {
		return attempt.ID == s.Val.ID
	})

	require.NotNil(t, saveAttempt)

}

func flushWaitingAttemptLog(t *testing.T){
	now := time.Now()
	attempt := (*WebhookRunner.State.WaitingAttempts.Val)[0]
	attempt.Val.NextTry = now.Add(5*time.Hour)
	var start_time time.Time = time.Now()
	
	ev, err := commons.EventFromType(commons.FlushWaitingAttemptType, attempt.Val, nil)
	require.NoError(t, err)
	log, err := commons.LogFromEvent(ev)
	require.NoError(t, err)

	
	err = Database.WriteLog(log.ID, string(log.Channel), log.Payload, log.CreatedAt)
	require.NoError(t, err)
	
	
	logs, err := Database.GetFreshLogs([]commons.Channel{commons.AttemptChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)
	
	FlushWaitingAttemptLog := (*logs)[0]
	
	WebhookRunner.HandleFreshLog(FlushWaitingAttemptLog)


	attempt = (*WebhookRunner.State.WaitingAttempts.Val)[0]
	require.True(t, attempt.Val.NextTry.Before(now.Add(5*time.Hour)))

}

func flushWaitingAttemptsLog( t *testing.T){
	now := time.Now()
	attempt := (*WebhookRunner.State.WaitingAttempts.Val)[0]
	attempt.Val.NextTry = now.Add(5*time.Hour)
	
	var start_time time.Time = time.Now()
	ev, err := commons.EventFromType(commons.FlushWaitingAttemptsType, attempt.Val, nil)
	require.NoError(t, err)
	log, err := commons.LogFromEvent(ev)
	require.NoError(t, err)

	
	err = Database.WriteLog(log.ID, string(log.Channel), log.Payload, log.CreatedAt)
	require.NoError(t, err)
	
	
	logs, err := Database.GetFreshLogs([]commons.Channel{commons.AttemptChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)
	
	FlushWaitingAttemptsLog := (*logs)[0]
	WebhookRunner.HandleFreshLog(FlushWaitingAttemptsLog)

	attempt = (*WebhookRunner.State.WaitingAttempts.Val)[0]
	require.True(t, attempt.Val.NextTry.Before(now.Add(5*time.Hour)))
}

func abortWaitingAttemptLog( t *testing.T){
	attempts, _, _ := Database.GetWaitingAttempts(0,1)
	require.Len(t, *attempts, 1)

	attempt := (*attempts)[0]

	var start_time time.Time = time.Now()
	_,err := Database.AbortAttempt(attempt.ID, "TEST LOG", true)
	require.NoError(t, err)

	logs, err := Database.GetFreshLogs([]commons.Channel{commons.AttemptChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)
	
	AbortAttemptsLog := (*logs)[0]
	WebhookRunner.HandleFreshLog(AbortAttemptsLog)

	
	require.Len(t, *WebhookRunner.State.WaitingAttempts.Val, 0)

}