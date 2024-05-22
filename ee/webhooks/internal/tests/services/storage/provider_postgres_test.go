package storage

import (
	"os"
	"testing"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/webhooks/internal/commons"
	"github.com/formancehq/webhooks/internal/services/storage/postgres"
	testutils "github.com/formancehq/webhooks/internal/tests"
	"github.com/stretchr/testify/require"
)

var Database storage.PostgresStore



func TestMain(m *testing.M) {
	testutils.StartPostgresServer()
	m.Run()
	testutils.StopPostgresServer()
}



func TestRun(t *testing.T) {
	var err error 
	Database, err  = testutils.GetStoreProvider()

	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	t.Run("InsertHook", insertHook)
	t.Run("GetHook", getHook)
	t.Run("ActivateHook", activateHook)
	t.Run("DeactivateHook", deactivateHook)

	t.Run("UpdateHookEndpoint", updateHookEndpoint)
	t.Run("UpdateHookSecret", updateHookSecret)
	t.Run("LoadHooks", loadHooks)
	t.Run("DeleteHook", deleteHook)

	t.Run("InsertHook", insertHook)
	t.Run("InsertAttempt", insertAttempt)

	t.Run("GetAttempt", getAttempt)

	t.Run("CompleteAttempt", completeAttempt)

	t.Run("AbortAttempt", abortAttempt) 

	t.Run("GetAbortedAttempts", getAbortedAttempts)

	t.Run("UpdateAttemptNextTry", updateAttemptNextTry)

	t.Run("LoadWaitingAttempts", loadWaitingAttempts)

}
func insertHook(t *testing.T){
	hook := commons.NewHook("TestHook", []string{"test","test2"}, "www.google.com", "xxx-foo-bar", true)
	require.NoError(t, Database.SaveHook(*hook))
}

func getHook(t *testing.T){
	hooks,hasMore, err := Database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Equal(t, hasMore, false)
	require.Len(t, *hooks, 1)
	hook, err := Database.GetHook((*hooks)[0].ID)
	require.NoError(t, err)
	require.NotNil(t, hook.ID)

}

func activateHook(t *testing.T){
	hooks,_, _ := Database.GetHooks(0, 1, "")
	hook := *(*hooks)[0]
	require.Equal(t, commons.DisableStatus, hook.Status)
	hook, err := Database.ActivateHook(hook.ID)
	require.NoError(t, err)
	require.Equal(t, commons.EnableStatus, hook.Status)

}

func deactivateHook(t *testing.T){
	hooks,_, _ := Database.GetHooks(0, 1, "")
	hook := *(*hooks)[0]
	require.Equal(t, commons.EnableStatus, hook.Status)
	hook, err := Database.DeactivateHook(hook.ID)
	require.NoError(t, err)
	require.Equal(t, commons.DisableStatus, hook.Status)
	
}

func updateHookEndpoint(t *testing.T){
	oldEndpoint := "www.google.com"
	newEndpoint := "www.newendpoint.com"
	hooks,_,_ := Database.GetHooks(0, 1, "")
	hook := *(*hooks)[0]
	require.Equal(t, oldEndpoint, hook.Endpoint)
	hook, err := Database.UpdateHookEndpoint(hook.ID, newEndpoint)
	require.NoError(t, err)
	require.Equal(t, newEndpoint, hook.Endpoint)
}

func updateHookSecret(t *testing.T){
	oldSecret := "xxx-foo-bar"
	newSecret := "new-xxx-foo-bar"
	hooks, _, _ := Database.GetHooks(0, 1, "")
	hook := *(*hooks)[0]
	require.Equal(t, oldSecret, hook.Secret)
	hook, err := Database.UpdateHookSecret(hook.ID, newSecret)
	require.NoError(t, err)
	require.Equal(t, newSecret, hook.Secret)
}

func loadHooks(t *testing.T){
	hooks,_, _ := Database.GetHooks(0, 1, "")
	hook := *(*hooks)[0]
	hook, _ = Database.ActivateHook(hook.ID)
	hooks, _ = Database.LoadHooks()
	require.Len(t, *hooks, 1)
}

func deleteHook(t *testing.T){
	hooks,_, _ := Database.GetHooks(0, 1, "")
	require.Len(t, *hooks, 1)
	hook := *(*hooks)[0]
	hook, err := Database.DeleteHook(hook.ID)
	require.NoError(t, err)
	hooks, _ = Database.LoadHooks()
	require.Len(t, *hooks, 0)
}

func insertAttempt(t *testing.T){
	hooks,_, _ := Database.GetHooks(0, 1, "")
	hook := *(*hooks)[0]
	attempt := commons.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "TEST")
	require.NoError(t, Database.SaveAttempt(*attempt))	
}

func getAttempt(t *testing.T){
	attempts,hasMore, err := Database.GetWaitingAttempts(0, 1)
	require.NoError(t, err)
	require.Equal(t, hasMore, false)
	require.Len(t, *attempts, 1)

	attempt, err := Database.GetAttempt((*attempts)[0].ID)
	require.NoError(t, err)
	require.NotNil(t, attempt.ID)
}

func completeAttempt(t *testing.T){
	attempts,_, _ := Database.GetWaitingAttempts(0, 1)
	attempt := *(*attempts)[0]
	attempt, err := Database.CompleteAttempt(attempt.ID)
	require.NoError(t, err)
	require.Equal(t, commons.SuccessStatus, attempt.Status)
	attempts,_, _ = Database.GetWaitingAttempts(0, 1)
	require.Len(t, *attempts, 0)
}

func abortAttempt(t *testing.T){
	hooks,_, _ := Database.GetHooks(0, 1, "")
	hook := *(*hooks)[0]
	attempt := commons.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "TEST")
	require.NoError(t, Database.SaveAttempt(*attempt))	
	attempts,_, _ := Database.GetWaitingAttempts(0, 1)
	attempt1 := *(*attempts)[0]
	attempt1, err := Database.AbortAttempt(attempt1.ID, "TESTABORT")
	require.NoError(t, err)
	require.Equal(t, commons.AbortStatus, attempt1.Status)
	require.Equal(t, "TESTABORT", string(attempt1.Comment))
	attempts,_, _ = Database.GetWaitingAttempts(0, 1)
	require.Len(t, *attempts, 0)
}

func getAbortedAttempts(t *testing.T){
	attempts,hasMore, err := Database.GetAbortedAttempts(0, 1)
	require.NoError(t, err)
	require.Equal(t, hasMore, false)
	require.Len(t, *attempts, 1)
}

func updateAttemptNextTry(t *testing.T){
	hooks,_, _ := Database.GetHooks(0, 1, "")
	hook := *(*hooks)[0]
	attempt := commons.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "TEST")
	require.NoError(t, Database.SaveAttempt(*attempt))	
	attempts,_, _ := Database.GetWaitingAttempts(0, 1)
	attempt1 := *(*attempts)[0]
	now := time.Now()
	attempt1, err := Database.UpdateAttemptNextTry(attempt1.ID, now.Add(25*time.Minute), attempt.LastHttpStatusCode)
	require.NoError(t, err)
	require.Equal(t, now.Add(25*time.Minute).UTC().Format(time.RFC3339), attempt1.NextTry.UTC().Format(time.RFC3339))
}

func loadWaitingAttempts(t *testing.T){
	attempts, err := Database.LoadWaitingAttempts()
	require.NoError(t, err)
	require.Len(t, *attempts, 1)
}