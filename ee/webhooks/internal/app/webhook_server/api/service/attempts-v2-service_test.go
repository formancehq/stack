package service

import (
	"testing"

	"github.com/formancehq/webhooks/internal/models"

	"github.com/stretchr/testify/require"
)

func TestRunAttemptV2(t *testing.T) {
	t.Run("InsertAttempt", v2InsertAttempt)
	t.Run("GetWaitingAttempts", v2GetWaitingAttempts)
	t.Run("AbortWaitingAttempt", v2AbortWaitingAttempt)
	t.Run("GetAbortedAttempts", v2GetAbortedAttempts)
	t.Run("RetryWaitingAttempt", v2RetryWaitingAttempt)
	t.Run("RetryWaitingAttempts", v2RetryWaitingAttempts)
}

func v2InsertAttempt(t *testing.T) {
	params := models.HookBodyParams{
		Name:     "Test1",
		Endpoint: "http://www.exemple-endpoint.com/valide",
		Secret:   "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh",
		Events:   []string{"event"}}

	resp := V2CreateHook(params)
	require.NoError(t, resp.Err)
	require.NotEmpty(t, resp.Data.ID)

	hook := resp.Data

	attempt := models.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "Attempt1")
	require.NoError(t, getDatabase().SaveAttempt(*attempt, true))
	attempt2 := models.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "Attempt2")
	require.NoError(t, getDatabase().SaveAttempt(*attempt2, true))
	attempt3 := models.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "Attempt3")
	require.NoError(t, getDatabase().SaveAttempt(*attempt3, true))
	attempt4 := models.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "Attempt4")
	require.NoError(t, getDatabase().SaveAttempt(*attempt4, true))
	attempt5 := models.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "AttemptAborted")
	attempt5.Status = models.AbortStatus
	attempt5.Comment = models.AbortUser
	require.NoError(t, getDatabase().SaveAttempt(*attempt5, true))

}

func v2GetWaitingAttempts(t *testing.T) {
	resp := V2GetWaitingAttempts("", 15)
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 4)
}

func v2AbortWaitingAttempt(t *testing.T) {
	temp := V2GetWaitingAttempts("", 15)
	require.NoError(t, temp.Err)
	attempt := temp.Data.Data[0]

	resp := V2AbortWaitingAttempt(attempt.ID)
	require.NoError(t, resp.Err)
	require.Equal(t, attempt.ID, resp.Data.ID)
	require.Equal(t, models.AbortUser, resp.Data.Comment)
}

func v2GetAbortedAttempts(t *testing.T) {
	resp := V2GetAbortedAttempts("", 15)
	require.NoError(t, resp.Err)
	require.Len(t, resp.Data.Data, 2)
}

func v2RetryWaitingAttempt(t *testing.T) {
	temp := V2GetWaitingAttempts("", 15)
	require.NoError(t, temp.Err)
	attempt := temp.Data.Data[0]
	resp := V2RetryWaitingAttempt(attempt.ID)
	require.NoError(t, resp.Err)
}

func v2RetryWaitingAttempts(t *testing.T) {
	resp := V2RetryWaitingAttempts()
	require.NoError(t, resp.Err)
}
