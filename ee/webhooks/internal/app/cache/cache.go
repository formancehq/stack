package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/sync"
	"github.com/formancehq/webhooks/internal/models"
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)

type CacheParams struct {
	TimeOut   int //ms
	MaxCall   int
	MaxRetry  int
	DelayPull int //second

}

func DefaultCacheParams() CacheParams {
	return CacheParams{
		TimeOut:   2000,
		MaxCall:   20,
		MaxRetry:  60,
		DelayPull: 1,
	}
}

type Cache struct {
	CacheParams CacheParams

	Queue       *sync.Queue
	StopChan    chan struct{}
	LogChannels []models.Channel

	State *State

	Database storeInterface.IStoreProvider
	Client   clientInterface.IHTTPClient
}

func (wr *Cache) Stop() {
	wr.StopChan <- struct{}{}
}

func (wr *Cache) StartHandleFreshLogs() {
	go wr.HandleFreshLogs(wr.StopChan)
}

func (wr *Cache) HandleFreshLogs(stopChan chan struct{}) {
	delay := time.Duration(wr.CacheParams.DelayPull) * time.Second
	ticker := time.NewTicker(time.Duration(delay))
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
				message := fmt.Sprintf("Cache:HandleFreshLogs() - LogChannels : %s : Error while attempting to reach the database: %s", wr.LogChannels, err)
				logging.Error(message)
				panic(message)
			}

			for _, log := range *logs {
				wr.HandleFreshLog(log)
			}
		}
	}

}

func (wr *Cache) HandleFreshLog(log *models.Log) {
	e, err := models.Event{}.FromPayload(log.Payload)
	if err != nil {
		message := fmt.Sprintf("Cache:HandleFreshLogs() - LogChannels : %s : Error while Event.FromPayload(log.payload): %s", wr.LogChannels, err)
		logging.Error(message)
		panic(message)
	}
	eventType := models.TypeFromEvent(e)
	switch eventType {

	case models.NewHookType:
		hook, err := wr.Database.GetHook(e.ID)

		if err != nil {
			message := fmt.Sprintf("Cache:HandleFreshLogs() - LogChannels : %s : Case NewHookType Error while attempting to reach the database: %s", wr.LogChannels, err)
			logging.Error(message)
			panic(message)
		} else {
			wr.State.AddNewHook(&hook)
		}

	case models.ChangeHookStatusType:
		strValue := e.Value.(string)
		switch models.HookStatus(strValue) {
		case models.EnableStatus:
			wr.State.ActivateHook(e.ID)
		case models.DisableStatus:
			wr.State.DisableHook(e.ID)
		case models.DeleteStatus:
			wr.State.DeleteHook(e.ID)
		}
	case models.ChangeHookEndpointType:
		wr.State.UpdateHookEndpoint(e.ID, e.Value.(string))
	case models.ChangeHookSecretType:
		wr.State.UpdateHookSecret(e.ID, e.Value.(string))
	case models.ChangeHookRetryType:
		wr.State.UpdateHookRetry(e.ID, e.Value.(bool))

	case models.NewWaitingAttemptType:
		attempt, err := wr.Database.GetAttempt(e.ID)
		if err != nil {
			message := fmt.Sprintf("Cache:HandleFreshLogs() - LogChannels : %s : Error while NewWaitingAttemptType :  wr.Database.GetAttempt(e.ID): %s", wr.LogChannels, err)
			logging.Error(message)
			panic(message)
		}
		wr.State.AddNewAttempt(&attempt)
	case models.FlushWaitingAttemptType:
		wr.State.FlushAttempt(e.ID)
	case models.FlushWaitingAttemptsType:
		wr.State.FlushAttempts()
	case models.AbortWaitingAttemptType:
		wr.State.AbortAttempt(e.ID)
	default:
		message := fmt.Sprintf("Cache:HandleFreshLogs() - LogChannels : %s : Unknow Log Type: %s", wr.LogChannels, e)
		logging.Error(message)
	}
}

func (wr *Cache) HandleRequest(ctx context.Context, sAttempt *models.SharedAttempt, sHook *models.SharedHook) (int, error) {
	wr.Queue.Lock()
	timeOut := time.Duration(wr.CacheParams.TimeOut) * time.Millisecond
	requestCtx, cancel := context.WithTimeout(ctx, timeOut)
	statusCode, err := wr.Client.Call(requestCtx, sHook.Val, sAttempt.Val, false)
	cancel()
	wr.Queue.Unlock()
	return statusCode, err
}

func NewCache(CacheParams CacheParams,
	database storeInterface.IStoreProvider,
	client clientInterface.IHTTPClient,
	logChannels ...models.Channel,
) *Cache {

	return &Cache{
		CacheParams: CacheParams,

		Queue:       sync.NewQueue(CacheParams.MaxCall),
		StopChan:    make(chan struct{}, 0),
		LogChannels: logChannels,
		State:       NewState(),

		Database: database,
		Client:   client,
	}

}
