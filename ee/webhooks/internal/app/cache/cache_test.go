package cache

import (
	"os"
	"testing"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/webhooks/internal/models"

	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"
	testutils "github.com/formancehq/webhooks/internal/testutils"

	"github.com/stretchr/testify/require"
)

var (
	database storage.PostgresStore
	cache    *Cache
)

func TestMain(m *testing.M) {

	testutils.StartPostgresServer()
	var err error
	database, err = testutils.GetStoreProvider()

	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	cache = NewCache(DefaultCacheParams(), database, testutils.NewHTTPClient())

	m.Run()
	testutils.StopPostgresServer()
}

func TestRunHandleLogs(t *testing.T) {

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

func newHookLog(t *testing.T) {
	var start_time time.Time = time.Now()

	newHook := models.NewHook("Test", []string{"testLog"}, "/localhost/", "", false)
	_, err := database.SaveHook(*newHook)
	require.NoError(t, err)

	logs, err := database.GetFreshLogs([]models.Channel{models.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	NewHookLog := (*logs)[0]

	cache.HandleFreshLog(NewHookLog)

	saveHook := cache.State.HooksById.Get(newHook.ID)
	require.NotNil(t, saveHook)
}

func changeHookStatusLog(t *testing.T) {
	var start_time time.Time = time.Now()
	hooks, _, err := database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]

	require.Equal(t, models.DisableStatus, hook.Status)
	hookIndex := hook.ID

	_, err = database.ActivateHook(hookIndex)
	require.NoError(t, err)

	logs, err := database.GetFreshLogs([]models.Channel{models.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookStatusLog := (*logs)[0]

	cache.HandleFreshLog(ChangeHookStatusLog)
	saveHook := cache.State.HooksById.Get(hook.ID)
	require.Equal(t, models.EnableStatus, saveHook.Val.Status)

	start_time = time.Now()
	_, err = database.DeactivateHook(hookIndex)
	require.NoError(t, err)

	logs, err = database.GetFreshLogs([]models.Channel{models.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookStatusLog = (*logs)[0]
	cache.HandleFreshLog(ChangeHookStatusLog)
	saveHook = cache.State.HooksById.Get(hook.ID)
	require.Equal(t, models.DisableStatus, saveHook.Val.Status)

	start_time = time.Now()
	_, err = database.ActivateHook(hookIndex)
	require.NoError(t, err)

	logs, err = database.GetFreshLogs([]models.Channel{models.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookStatusLog = (*logs)[0]
	cache.HandleFreshLog(ChangeHookStatusLog)
	saveHook = cache.State.HooksById.Get(hook.ID)
	require.Equal(t, models.EnableStatus, saveHook.Val.Status)

	start_time = time.Now()
	_, err = database.DeleteHook(hookIndex)
	require.NoError(t, err)

	logs, err = database.GetFreshLogs([]models.Channel{models.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookStatusLog = (*logs)[0]
	cache.HandleFreshLog(ChangeHookStatusLog)
	saveHook = cache.State.HooksById.Get(hook.ID)
	require.Nil(t, saveHook)

}

func changeHookSecretLog(t *testing.T) {
	var start_time time.Time = time.Now()

	newHook := models.NewHook("Test", []string{"testLog"}, "/localhost/", "", false)
	_, err := database.SaveHook(*newHook)
	require.NoError(t, err)

	logs, err := database.GetFreshLogs([]models.Channel{models.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	NewHookLog := (*logs)[0]

	cache.HandleFreshLog(NewHookLog)

	hooks, _, err := database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]
	newSecret := "Y2VjaWVzdHVuc2VjcmV0dmFsaWRlcyEh"
	start_time = time.Now()
	_, err = database.UpdateHookSecret(hook.ID, newSecret)
	require.NoError(t, err)

	logs, err = database.GetFreshLogs([]models.Channel{models.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookSecretLog := (*logs)[0]

	cache.HandleFreshLog(ChangeHookSecretLog)
	saveHook := cache.State.HooksById.Get(hook.ID)
	require.Equal(t, newSecret, saveHook.Val.Secret)

}

func changeHookEndpointLog(t *testing.T) {
	hooks, _, err := database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]
	newEndpoint := "www.google.fr/top"

	var start_time time.Time = time.Now()
	_, err = database.UpdateHookEndpoint(hook.ID, newEndpoint)
	require.NoError(t, err)

	logs, err := database.GetFreshLogs([]models.Channel{models.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookSecretLog := (*logs)[0]

	cache.HandleFreshLog(ChangeHookSecretLog)

	saveHook := cache.State.HooksById.Get(hook.ID)
	require.Equal(t, newEndpoint, saveHook.Val.Endpoint)
}

func changeHookRetryLog(t *testing.T) {
	hooks, _, err := database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]
	require.False(t, hook.Retry)

	var start_time time.Time = time.Now()
	_, err = database.UpdateHookRetry(hook.ID, true)
	require.NoError(t, err)

	logs, err := database.GetFreshLogs([]models.Channel{models.HookChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	ChangeHookRetryLog := (*logs)[0]

	cache.HandleFreshLog(ChangeHookRetryLog)

	saveHook := cache.State.HooksById.Get(hook.ID)
	require.True(t, saveHook.Val.Retry)
}

func newWaitingAttemptLog(t *testing.T) {
	hooks, _, err := database.GetHooks(0, 1, "")
	require.NoError(t, err)
	require.Len(t, *hooks, 1)
	hook := (*hooks)[0]

	attempt := models.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], "blabla")

	var start_time time.Time = time.Now()
	err = database.SaveAttempt(*attempt, true)
	require.NoError(t, err)

	logs, err := database.GetFreshLogs([]models.Channel{models.AttemptChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	NewWaitingAttemptLog := (*logs)[0]

	cache.HandleFreshLog(NewWaitingAttemptLog)

	saveAttempt := cache.State.WaitingAttempts.FindElement(func(s *models.SharedAttempt) bool {
		return attempt.ID == s.Val.ID
	})

	require.NotNil(t, saveAttempt)

}

func flushWaitingAttemptLog(t *testing.T) {
	now := time.Now()
	attempt := (*cache.State.WaitingAttempts.Val)[0]
	attempt.Val.NextTry = now.Add(5 * time.Hour)
	var start_time time.Time = time.Now()

	ev, err := models.EventFromType(models.FlushWaitingAttemptType, attempt.Val, nil)
	require.NoError(t, err)
	log, err := models.LogFromEvent(ev)
	require.NoError(t, err)

	err = database.WriteLog(log.ID, string(log.Channel), log.Payload, log.CreatedAt)
	require.NoError(t, err)

	logs, err := database.GetFreshLogs([]models.Channel{models.AttemptChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	FlushWaitingAttemptLog := (*logs)[0]

	cache.HandleFreshLog(FlushWaitingAttemptLog)

	attempt = (*cache.State.WaitingAttempts.Val)[0]
	require.True(t, attempt.Val.NextTry.Before(now.Add(5*time.Hour)))

}

func flushWaitingAttemptsLog(t *testing.T) {
	now := time.Now()
	attempt := (*cache.State.WaitingAttempts.Val)[0]
	attempt.Val.NextTry = now.Add(5 * time.Hour)

	var start_time time.Time = time.Now()
	ev, err := models.EventFromType(models.FlushWaitingAttemptsType, attempt.Val, nil)
	require.NoError(t, err)
	log, err := models.LogFromEvent(ev)
	require.NoError(t, err)

	err = database.WriteLog(log.ID, string(log.Channel), log.Payload, log.CreatedAt)
	require.NoError(t, err)

	logs, err := database.GetFreshLogs([]models.Channel{models.AttemptChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	FlushWaitingAttemptsLog := (*logs)[0]
	cache.HandleFreshLog(FlushWaitingAttemptsLog)

	attempt = (*cache.State.WaitingAttempts.Val)[0]
	require.True(t, attempt.Val.NextTry.Before(now.Add(5*time.Hour)))
}

func abortWaitingAttemptLog(t *testing.T) {
	attempts, _, _ := database.GetWaitingAttempts(0, 1)
	require.Len(t, *attempts, 1)

	attempt := (*attempts)[0]

	var start_time time.Time = time.Now()
	_, err := database.AbortAttempt(attempt.ID, "TEST LOG", true)
	require.NoError(t, err)

	logs, err := database.GetFreshLogs([]models.Channel{models.AttemptChannel}, start_time)
	require.NoError(t, err)
	require.Len(t, *logs, 1)

	AbortAttemptsLog := (*logs)[0]
	cache.HandleFreshLog(AbortAttemptsLog)

	require.Len(t, *cache.State.WaitingAttempts.Val, 0)

}
