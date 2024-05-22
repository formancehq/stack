package webhookcollector

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	cache "github.com/formancehq/webhooks/internal/app/cache"
	"github.com/formancehq/webhooks/internal/models"
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
	utilsHttp "github.com/formancehq/webhooks/internal/utils/http"
)

type Collector struct {
	cache.Cache
}

func (c *Collector) Run() {

	ticker := time.NewTicker(1 * time.Second)
	go c.State.RoutineTime(c.StopChan, ticker, c.HandleWaitingAttempts)
}

func (c *Collector) Init() {

	c.StartHandleFreshLogs()

	hooks, err := c.Database.LoadHooks()
	if err != nil {
		c.Stop()
		logging.Error(err.Error())
		os.Exit(1)
	}
	c.State.LoadHooks(hooks)

	attempts, err := c.Database.LoadWaitingAttempts()
	if err != nil {
		c.Stop()
		logging.Error(err.Error())
		os.Exit(1)
	}
	c.State.LoadWaitingAttempts(attempts)

}

func (c *Collector) HandleWaitingAttempts() {

	if c.State.WaitingAttempts.Size() == 0 {
		return
	}

	now := time.Now()

	wAttempts := c.State.WaitingAttempts.Empty()
	toHandles := models.NewSharedAttempts()

	wAttempts.Apply(func(s *models.SharedAttempt) {

		if s.Val.NextTry.Before(now) {
			toHandles.Add(s)
		} else {
			c.State.WaitingAttempts.Add(s)
		}

	})

	toHandles.AsyncApply(c.AsyncHandleSharedAttempt)

}

func (c *Collector) AsyncHandleSharedAttempt(sAttempt *models.SharedAttempt, wg *sync.WaitGroup) {

	defer wg.Done()
	sHook := c.State.HooksById.Get(sAttempt.Val.HookID)

	if sHook == nil {
		c.handleMissingHook(sAttempt)
		return
	}

	if !sHook.Val.IsActive() {
		c.handleDisabledHook(sAttempt)
		return
	}

	statusCode, err := c.HandleRequest(context.Background(), sAttempt, sHook)

	if err != nil {

		c.State.WaitingAttempts.Add(sAttempt)
		return
	}

	if sAttempt.Val.HookEndpoint != sHook.Val.Endpoint {
		sAttempt.Val.HookEndpoint = sHook.Val.Endpoint
		go func() {
			_, err := c.Database.UpdateAttemptEndpoint(sAttempt.Val.ID, sAttempt.Val.HookEndpoint)
			if err != nil {
				message := fmt.Sprintf("Collector:AsyncHandleSharedAttempt:Database.UpdateAttemptEndpoint() : %s", err)
				logging.Error(message)
				panic(message)
			}
		}()
	}

	sAttempt.Val.LastHttpStatusCode = statusCode

	if utilsHttp.IsHTTPRequestSuccess(statusCode) {

		c.handleSuccess(sAttempt)
	} else {

		c.handleNextRetry(sAttempt)
	}

}

func (c *Collector) handleSuccess(sAttempt *models.SharedAttempt) {
	models.SetSuccesStatus(sAttempt.Val)
	_, err := c.Database.CompleteAttempt(sAttempt.Val.ID)
	if err != nil {
		message := fmt.Sprintf("Collector:handleSuccess:Database.CompleteAttempt() : %s", err)
		logging.Error(message)
		panic(message)
	}
}

func (c *Collector) handleNextRetry(sAttempt *models.SharedAttempt) {
	sAttempt.Val.NbTry += 1

	if c.CacheParams.MaxRetry <= sAttempt.Val.NbTry {
		models.SetAbortMaxRetryStatus(sAttempt.Val)
		_, err := c.Database.AbortAttempt(sAttempt.Val.ID, string(sAttempt.Val.Comment), false)
		if err != nil {
			message := fmt.Sprintf("Collector:handleNextRetry:Database.AbortAttempt: %s", err)
			logging.Error(message)
			panic(message)
		}
	} else {

		models.SetNextRetry(sAttempt.Val)
		c.State.WaitingAttempts.Add(sAttempt)

		go func() {
			_, err := c.Database.UpdateAttemptNextTry(sAttempt.Val.ID, sAttempt.Val.NextTry, sAttempt.Val.LastHttpStatusCode)
			if err != nil {
				message := fmt.Sprintf("Collector:handleNextRetry:Database.UpdateAttemptNextTry: %s", err)
				logging.Error(message)
				panic(message)
			}
		}()
	}
}

func (c *Collector) handleMissingHook(sAttempt *models.SharedAttempt) {
	models.SetAbortMissingHookStatus(sAttempt.Val)
	_, err := c.Database.AbortAttempt(sAttempt.Val.ID, string(sAttempt.Val.Comment), false)
	if err != nil {
		message := fmt.Sprintf("Collector:handleMissingHook:Database.AbortAttempt: %s", err)
		logging.Error(message)
		panic(message)

	}
}

func (c *Collector) handleDisabledHook(sAttempt *models.SharedAttempt) {
	models.SetAbortDisableHook(sAttempt.Val)
	_, err := c.Database.AbortAttempt(sAttempt.Val.ID, string(sAttempt.Val.Comment), false)
	if err != nil {
		message := fmt.Sprintf("Collector:handleDisabledHook:Database.AbortAttempt: %s", err)
		logging.Error(message)
		panic(message)
	}
}

func NewCollector(cacheParams cache.CacheParams, database storeInterface.IStoreProvider,
	client clientInterface.IHTTPClient) *Collector {

	return &Collector{
		Cache: *cache.NewCache(cacheParams, database, client, models.HookChannel, models.AttemptChannel),
	}

}
