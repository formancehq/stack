package webhookcollector

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/logging"

	"github.com/formancehq/webhooks/internal/app/cache"
	"github.com/formancehq/webhooks/internal/models"

	"github.com/stretchr/testify/require"

	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"
	testutils "github.com/formancehq/webhooks/internal/testutils"
)

var Database storage.PostgresStore
var WebhookCollector Collector

func TestMain(m *testing.M) {
	testutils.StartPostgresServer()
	var err error
	Database, err = testutils.GetStoreProvider()
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	WebhookCollector = *NewCollector(cache.DefaultCacheParams(),
		Database, testutils.NewHTTPClient())

	m.Run()
	testutils.StopPostgresServer()
}

var ActiveGoodHook *models.Hook
var ActiveBadHook *models.Hook
var DeactiveHook *models.Hook

var TestServer *http.Server

var GoodHandler func(http.ResponseWriter, *http.Request)
var BadHandler func(http.ResponseWriter, *http.Request)

func TestRunCollector(t *testing.T) {

	GoodHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	}

	BadHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "NO OK")
	}

	TestServer = testutils.NewHTTPServer(45679, [2]interface{}{"/good", http.HandlerFunc(GoodHandler)}, [2]interface{}{"/bad", http.HandlerFunc(BadHandler)})
	defer TestServer.Close()

	ActiveGoodHook = models.NewHook("HookGood", []string{"testevent"}, "http://127.0.0.1:45679/good", "", false)
	ActiveGoodHook.Status = models.EnableStatus
	ActiveBadHook = models.NewHook("HookBad", []string{"testevent"}, "http://127.0.0.1:45679/bad", "", false)
	ActiveBadHook.Status = models.EnableStatus
	DeactiveHook = models.NewHook("HookDeactive", []string{"testevent"}, "http://127.0.0.1:45679/good", "", false)

	_, err := Database.SaveHook(*ActiveGoodHook)
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	_, err = Database.SaveHook(*ActiveBadHook)
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	_, err = Database.SaveHook(*DeactiveHook)
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	WebhookCollector.State.AddNewHook(ActiveBadHook)
	WebhookCollector.State.AddNewHook(ActiveGoodHook)
	WebhookCollector.State.AddNewHook(DeactiveHook)

	t.Run("HandleGoodHook", HandleGoodHook)

	t.Run("HandleBadHook", HandleBadHook)

	t.Run("HandleDeactiveHook", HandleDeactiveHook)

	t.Run("HandleMissingHook", HandleMissingHook)
}

func HandleGoodHook(t *testing.T) {
	var wg sync.WaitGroup
	sAttempt := models.NewSharedAttempt(ActiveGoodHook.ID,
		ActiveGoodHook.Name, ActiveGoodHook.Endpoint, "testevent", "payload good")
	Database.SaveAttempt(*sAttempt.Val, true)

	wg.Add(1)
	WebhookCollector.AsyncHandleSharedAttempt(sAttempt, &wg)
	wg.Wait()

	attempt, err := Database.GetAttempt(sAttempt.Val.ID)
	require.NoError(t, err)

	require.Equal(t, models.SuccessStatus, attempt.Status)
}

func HandleBadHook(t *testing.T) {
	var wg sync.WaitGroup
	sAttempt := models.NewSharedAttempt(ActiveBadHook.ID,
		ActiveBadHook.Name, ActiveBadHook.Endpoint, "testevent", "payload bad")
	Database.SaveAttempt(*sAttempt.Val, true)

	wg.Add(1)
	WebhookCollector.AsyncHandleSharedAttempt(sAttempt, &wg)

	attempt, err := Database.GetAttempt(sAttempt.Val.ID)
	require.NoError(t, err)

	require.Equal(t, models.WaitingStatus, attempt.Status)

	find := WebhookCollector.State.WaitingAttempts.Find(sAttempt)
	require.NotEqual(t, -1, find)

}

func HandleDeactiveHook(t *testing.T) {
	var wg sync.WaitGroup
	sAttempt := models.NewSharedAttempt(DeactiveHook.ID,
		DeactiveHook.Name, DeactiveHook.Endpoint, "testevent", "payload bad")
	Database.SaveAttempt(*sAttempt.Val, true)

	wg.Add(1)
	WebhookCollector.AsyncHandleSharedAttempt(sAttempt, &wg)
	wg.Wait()

	attempt, err := Database.GetAttempt(sAttempt.Val.ID)
	require.NoError(t, err)

	require.Equal(t, models.AbortStatus, attempt.Status)
	require.Equal(t, models.AbortDisabledHook, attempt.Comment)

}

func HandleMissingHook(t *testing.T) {
	var wg sync.WaitGroup
	sAttempt := models.NewSharedAttempt("noID",
		"no name", "no endpoint", "testevent", "payload bad")
	Database.SaveAttempt(*sAttempt.Val, true)

	wg.Add(1)
	WebhookCollector.AsyncHandleSharedAttempt(sAttempt, &wg)

	attempt, err := Database.GetAttempt(sAttempt.Val.ID)
	require.NoError(t, err)

	require.Equal(t, models.AbortStatus, attempt.Status)
	require.Equal(t, models.AbortMissingHook, attempt.Comment)

}
