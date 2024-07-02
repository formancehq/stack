package webhookcollector

import (
	_ "errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/webhooks/internal/commons"
	"github.com/stretchr/testify/require"

	testutils "github.com/formancehq/webhooks/internal/tests"
)

var ActiveGoodHook *commons.Hook
var ActiveBadHook *commons.Hook
var DeactiveHook *commons.Hook

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

	TestServer = testutils.NewHTTPServer(4567, [2]interface{}{"/good", http.HandlerFunc(GoodHandler)}, [2]interface{}{"/bad", http.HandlerFunc(BadHandler)})
	defer TestServer.Close()

	ActiveGoodHook = commons.NewHook("HookGood", []string{"testevent"}, "http://127.0.0.1:4567/good", "", false)
	ActiveGoodHook.Status = commons.EnableStatus
	ActiveBadHook = commons.NewHook("HookBad", []string{"testevent"}, "http://127.0.0.1:4567/bad", "", false)
	ActiveBadHook.Status = commons.EnableStatus
	DeactiveHook = commons.NewHook("HookDeactive", []string{"testevent"}, "http://127.0.0.1:4567/good", "", false)

	err := Database.SaveHook(*ActiveGoodHook)
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	err = Database.SaveHook(*ActiveBadHook)
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	err = Database.SaveHook(*DeactiveHook)
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
	sAttempt := commons.NewSharedAttempt(ActiveGoodHook.ID,
		ActiveGoodHook.Name, ActiveGoodHook.Endpoint, "testevent", "payload good")
	Database.SaveAttempt(*sAttempt.Val, true)

	wg.Add(1)
	WebhookCollector.AsyncHandleSharedAttempt(sAttempt, &wg)
	wg.Wait()

	attempt, err := Database.GetAttempt(sAttempt.Val.ID)
	require.NoError(t, err)

	require.Equal(t, commons.SuccessStatus, attempt.Status)
}

func HandleBadHook(t *testing.T) {
	var wg sync.WaitGroup
	sAttempt := commons.NewSharedAttempt(ActiveBadHook.ID,
		ActiveBadHook.Name, ActiveBadHook.Endpoint, "testevent", "payload bad")
	Database.SaveAttempt(*sAttempt.Val, true)

	wg.Add(1)
	WebhookCollector.AsyncHandleSharedAttempt(sAttempt, &wg)

	attempt, err := Database.GetAttempt(sAttempt.Val.ID)
	require.NoError(t, err)

	require.Equal(t, commons.WaitingStatus, attempt.Status)

	find := WebhookCollector.State.WaitingAttempts.Find(sAttempt)
	require.NotEqual(t, -1, find)

}

func HandleDeactiveHook(t *testing.T) {
	var wg sync.WaitGroup
	sAttempt := commons.NewSharedAttempt(DeactiveHook.ID,
		DeactiveHook.Name, DeactiveHook.Endpoint, "testevent", "payload bad")
	Database.SaveAttempt(*sAttempt.Val, true)

	wg.Add(1)
	WebhookCollector.AsyncHandleSharedAttempt(sAttempt, &wg)

	attempt, err := Database.GetAttempt(sAttempt.Val.ID)
	require.NoError(t, err)

	require.Equal(t, commons.AbortStatus, attempt.Status)
	require.Equal(t, commons.AbortDisabledHook, attempt.Comment)

}

func HandleMissingHook(t *testing.T) {
	var wg sync.WaitGroup
	sAttempt := commons.NewSharedAttempt("noID",
		"no name", "no endpoint", "testevent", "payload bad")
	Database.SaveAttempt(*sAttempt.Val, true)

	wg.Add(1)
	WebhookCollector.AsyncHandleSharedAttempt(sAttempt, &wg)

	attempt, err := Database.GetAttempt(sAttempt.Val.ID)
	require.NoError(t, err)

	require.Equal(t, commons.AbortStatus, attempt.Status)
	require.Equal(t, commons.AbortMissingHook, attempt.Comment)

}
