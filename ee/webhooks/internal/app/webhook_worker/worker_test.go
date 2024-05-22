package webhookworker

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/webhooks/internal/app/cache"
	"github.com/formancehq/webhooks/internal/models"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"
	"github.com/stretchr/testify/require"

	testutils "github.com/formancehq/webhooks/internal/testutils"
)

var Database storage.PostgresStore
var WebhookWorker Worker

func TestMain(m *testing.M) {
	testutils.StartPostgresServer()
	var err error
	Database, err = testutils.GetStoreProvider()
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	WebhookWorker = *NewWorker(cache.DefaultCacheParams(),
		Database, testutils.NewHTTPClient())

	m.Run()
	testutils.StopPostgresServer()
}

var ActiveGoodHook *models.SharedHook
var ActiveBadHook *models.SharedHook
var DeactiveHook *models.SharedHook

var TestServer *http.Server

var GoodHandler func(http.ResponseWriter, *http.Request)
var BadHandler func(http.ResponseWriter, *http.Request)

var HandleHookTrigged func(sHook *models.SharedHook, wg *sync.WaitGroup)
var GlobalError error

func TestRunCollector(t *testing.T) {

	GoodHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

	}

	BadHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)

	}

	TestServer = testutils.NewHTTPServer(45678, [2]interface{}{"/good", http.HandlerFunc(GoodHandler)}, [2]interface{}{"/bad", http.HandlerFunc(BadHandler)})
	defer TestServer.Close()

	ActiveGoodHook = models.NewSharedHook("HookGood", []string{"webhook.testevent"}, "http://127.0.0.1:45678/good", "", false)
	ActiveGoodHook.Val.Status = models.EnableStatus
	ActiveBadHook = models.NewSharedHook("HookBad", []string{"webhook.testevent"}, "http://127.0.0.1:45678/bad", "", false)
	ActiveBadHook.Val.Status = models.EnableStatus
	_, err := Database.SaveHook(*ActiveGoodHook.Val)
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	_, err = Database.SaveHook(*ActiveBadHook.Val)
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	HandleHookTrigged = WebhookWorker.HandlerTriggedHookFactory(context.Background(), "testevent", "payload-general", GlobalError)

	WebhookWorker.State.AddNewHook(ActiveBadHook.Val)
	WebhookWorker.State.AddNewHook(ActiveGoodHook.Val)

	t.Run("HandleGoodHook", HandleGoodHook)

	t.Run("HandleBadHook", HandleBadHook)

	t.Run("HandleMessage", HandleGoodMessage)
}

func HandleGoodHook(t *testing.T) {
	GlobalError = nil
	var wg sync.WaitGroup

	wg.Add(1)
	sAttempt := models.NewSharedAttempt(ActiveGoodHook.Val.ID,
		ActiveGoodHook.Val.Name, ActiveGoodHook.Val.Endpoint, "webhook.goodevent", "payload good")
	statusCode, err := WebhookWorker.HandleRequest(context.Background(), sAttempt, ActiveGoodHook)
	require.NoError(t, err)
	require.Equal(t, 200, statusCode)

	HandleHookTrigged(ActiveGoodHook, &wg)
	require.NoError(t, GlobalError)

}

func HandleBadHook(t *testing.T) {
	GlobalError = nil
	var wg sync.WaitGroup

	sAttempt := models.NewSharedAttempt(ActiveBadHook.Val.ID,
		ActiveBadHook.Val.Name, ActiveBadHook.Val.Endpoint, "webhook.badevent", "payload bad")

	wg.Add(1)
	statusCode, err := WebhookWorker.HandleRequest(context.Background(), sAttempt, ActiveBadHook)
	require.NoError(t, err)

	require.Equal(t, 400, statusCode)

	HandleHookTrigged(ActiveBadHook, &wg)
	require.NoError(t, GlobalError)

}

func HandleGoodMessage(t *testing.T) {
	eventMessage := publish.EventMessage{
		Date:    time.Now(),
		App:     "webhook",
		Version: "1020",
		Type:    "goodevent",
		Payload: "blabla",
	}

	attempts, err := Database.LoadWaitingAttempts()
	require.NoError(t, err)
	nbBefore := len(*attempts)

	data, err := json.Marshal(eventMessage)
	require.NoError(t, err)

	message := message.Message{
		Payload: data,
	}
	gerr := WebhookWorker.HandleMessage(&message)

	require.NoError(t, gerr)

	attempts, err = Database.LoadWaitingAttempts()
	nbAfter := len(*attempts)

	require.Equal(t, nbBefore, nbAfter)

}

func HandleBadMessage(t *testing.T) {
	eventMessage := publish.EventMessage{
		Date:    time.Now(),
		App:     "webhook",
		Version: "1020",
		Type:    "badevent",
		Payload: "blabla",
	}

	attempts, err := Database.LoadWaitingAttempts()
	require.NoError(t, err)
	nbBefore := len(*attempts)

	data, err := json.Marshal(eventMessage)
	require.NoError(t, err)

	message := message.Message{
		Payload: data,
	}
	gerr := WebhookWorker.HandleMessage(&message)

	require.NoError(t, gerr)

	attempts, err = Database.LoadWaitingAttempts()
	nbAfter := len(*attempts)

	require.Equal(t, nbBefore, nbAfter-1)

}
