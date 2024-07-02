package component

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/sync"
	"github.com/formancehq/webhooks/internal/commons"
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)

type RunnerParams struct {
	TimeOut   int //ms
	MaxCall   int
	MaxRetry  int
	DelayPull int //second

}

func DefaultRunnerParams() RunnerParams {
	return RunnerParams{
		TimeOut:   2000,
		MaxCall:   20,
		MaxRetry:  60,
		DelayPull: 1,
	}
}

type WebhookRunner struct {
	RunnerParams RunnerParams

	Queue       *sync.Queue
	StopChan    chan struct{}
	LogChannels []commons.Channel

	State *commons.State

	Database storeInterface.IStoreProvider
	Client   clientInterface.IHTTPClient
}

func (wr *WebhookRunner) Stop() {
	wr.StopChan <- struct{}{}
}

func (wr *WebhookRunner) Run() {

}

func (wr *WebhookRunner) StartHandleFreshLogs() {
	go wr.HandleFreshLogs(wr.StopChan)
}

func (wr *WebhookRunner) HandleFreshLogs(stopChan chan struct{}) {

	ticker := time.NewTicker(time.Duration(wr.RunnerParams.DelayPull) * time.Second)
	var last_time time.Time = time.Now()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			freezeTime := time.Now()
			logs, err := wr.Database.GetFreshLogs(wr.LogChannels, last_time)
			last_time = freezeTime
			if err != nil {
				message := fmt.Sprintf("WebhookRunner:HandleFreshLogs() - LogChannels : %s : Error while attempting to reach the database: %x", wr.LogChannels, err)
				logging.Error(message)
				panic(message)
			}

			for _, log := range *logs {
				wr.HandleFreshLog(log)
			}
		}
	}

}

func (wr *WebhookRunner) HandleFreshLog(log *commons.Log) {
	fmt.Println("HandleFreshLog")
	e, err := commons.Event{}.FromPayload(log.Payload)
	if err != nil {
		message := fmt.Sprintf("WebhookRunner:HandleFreshLogs() - LogChannels : %s : Error while Event.FromPayload(log.payload): %x", wr.LogChannels, err)
		logging.Error(message)
		panic(message)
	}
	eventType := commons.TypeFromEvent(e)

	switch eventType {

	case commons.NewHookType:
		hook, err := wr.Database.GetHook(e.ID)
		if err != nil {
			message := fmt.Sprintf("WebhookRunner:HandleFreshLogs() - LogChannels : %s : Case NewHookType Error while attempting to reach the database: %x", wr.LogChannels, err)
			logging.Error(message)
			panic(message)
		} else {
			wr.State.AddNewHook(&hook)
		}

	case commons.ChangeHookStatusType:
		fmt.Println(e.Value)
		strValue := e.Value.(string)
		switch commons.HookStatus(strValue) {
		case commons.EnableStatus:
			fmt.Println("ENABLE CASE")
			wr.State.ActivateHook(e.ID)
		case commons.DisableStatus:
			wr.State.DisableHook(e.ID)
		case commons.DeleteStatus:
			wr.State.DeleteHook(e.ID)
		}
	case commons.ChangeHookEndpointType:
		wr.State.UpdateHookEndpoint(e.ID, e.Value.(string))
	case commons.ChangeHookSecretType:
		wr.State.UpdateHookSecret(e.ID, e.Value.(string))
	case commons.ChangeHookRetryType:
		wr.State.UpdateHookRetry(e.ID, e.Value.(bool))

	case commons.NewWaitingAttemptType:
		attempt, err := wr.Database.GetAttempt(e.ID)
		if err != nil {
			message := fmt.Sprintf("WebhookRunner:HandleFreshLogs() - LogChannels : %s : Error while NewWaitingAttemptType :  wr.Database.GetAttempt(e.ID): %x", wr.LogChannels, err)
			logging.Error(message)
			panic(message)
		}

		wr.State.AddNewAttempt(&attempt)
	case commons.FlushWaitingAttemptType:
		wr.State.FlushAttempt(e.ID)
	case commons.FlushWaitingAttemptsType:
		wr.State.FlushAttempts()
	case commons.AbortWaitingAttemptType:
		wr.State.AbortAttempt(e.ID)
	default:
		message := fmt.Sprintf("WebhookRunner:HandleFreshLogs() - LogChannels : %s : Unknow Log Type: %x", wr.LogChannels, e)
		logging.Error(message)
	}
}

func (wr *WebhookRunner) HandleRequest(ctx context.Context, sAttempt *commons.SharedAttempt, sHook *commons.SharedHook) (int, error) {
	wr.Queue.Lock()
	timeOut := time.Duration(wr.RunnerParams.TimeOut) * time.Millisecond
	requestCtx, cancel := context.WithTimeout(ctx, timeOut)
	statusCode, err := wr.Client.Call(requestCtx, sHook.Val, sAttempt.Val, false)
	cancel()
	wr.Queue.Unlock()
	return statusCode, err
}

func NewWebhookRunner(runnerParams RunnerParams,
	database storeInterface.IStoreProvider,
	client clientInterface.IHTTPClient,
	logChannels ...commons.Channel,
) *WebhookRunner {

	return &WebhookRunner{
		RunnerParams: runnerParams,

		Queue:       sync.NewQueue(runnerParams.MaxCall),
		StopChan:    make(chan struct{}, 0),
		LogChannels: logChannels,
		State:       commons.NewState(),

		Database: database,
		Client:   client,
	}

}
