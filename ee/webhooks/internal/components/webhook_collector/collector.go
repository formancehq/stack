package webhookcollector

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/webhooks/internal/commons"
	component "github.com/formancehq/webhooks/internal/components/commons"
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)

type Collector struct {
	wg sync.WaitGroup
	component.WebhookRunner
}

func (c *Collector) Run() {

	ticker := time.NewTicker(1*time.Second)


	go c.State.RoutineTime(c.StopChan, ticker, c.HandleWaitingAttempts)
}

func (c *Collector) Init(){

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

func (c *Collector) HandleWaitingAttempts(){

	if(c.State.WaitingAttempts.Size() == 0){return}
				
	now := time.Now()

	wAttempts := c.State.WaitingAttempts.Empty()
	toHandles := commons.NewSharedAttempts()
	
	
	wAttempts.Apply(func(s *commons.SharedAttempt) {
			
		if(s.Val.NextTry.Before(now)){ 
			toHandles.Add(s)
		} else {
			c.State.WaitingAttempts.Add(s)
		}

	})
	
	toHandles.AsyncApply(c.AsyncHandleSharedAttempt)

}


func (c *Collector) AsyncHandleSharedAttempt(sAttempt *commons.SharedAttempt, wg *sync.WaitGroup){

	defer wg.Done()
	sHook := c.State.HooksById.Get(sAttempt.Val.HookID)
	
	if(sHook == nil){
		c.handleMissingHook(sAttempt)
		return
	}
	
	if(!sHook.Val.IsActive()){
		c.handleDisabledHook(sAttempt)
		return
	}

	statusCode, err := c.HandleRequest(context.Background(), sAttempt, sHook)
	
	if (err != nil) {
		logging.Errorf("Collector:AsyncHandleSharedAttempt:HandleRequest : %x", err)
		c.State.WaitingAttempts.Add(sAttempt)
		return
	}

	sAttempt.Val.LastHttpStatusCode = statusCode

	if(commons.IsHTTPRequestSuccess(statusCode)) {
		c.handleSuccess(sAttempt)
	} else {
		c.handleNextRetry(sAttempt)
	}

}

func (c *Collector) handleSuccess(sAttempt *commons.SharedAttempt){
	commons.SetSuccesStatus(sAttempt.Val)
	_, err := c.Database.CompleteAttempt(sAttempt.Val.ID)
	if err != nil {
		message := fmt.Sprintf("Collector:handleSuccess:Database.CompleteAttempt() : %x", err)
		logging.Error(message)
		panic(message)
	}
}

func (c *Collector) handleNextRetry(sAttempt *commons.SharedAttempt){
	sAttempt.Val.NbTry += 1
	commons.SetNextRetry(sAttempt.Val)
	c.State.WaitingAttempts.Add(sAttempt)
	
	go func(){
		_, err := c.Database.UpdateAttemptNextTry(sAttempt.Val.ID, sAttempt.Val.NextTry, sAttempt.Val.LastHttpStatusCode)
		if err != nil {
			message := fmt.Sprintf("Collector:handleNextRetry:Database.UpdateAttemptNextTry: %x", err)
			logging.Error(message)
			panic(message)
		}
	}()
}

func (c *Collector) handleMissingHook(sAttempt *commons.SharedAttempt){
	commons.SetAbortMissingHookStatus(sAttempt.Val)
	_, err := c.Database.AbortAttempt(sAttempt.Val.ID, string(sAttempt.Val.Comment), false)
	if err != nil {
		message := fmt.Sprintf("Collector:handleMissingHook:Database.AbortAttempt: %x", err)
		logging.Error(message)
		panic(message)
		
	}
}

func (c *Collector) handleDisabledHook(sAttempt *commons.SharedAttempt){
	commons.SetAbortDisableHook(sAttempt.Val)
	_, err := c.Database.AbortAttempt(sAttempt.Val.ID, string(sAttempt.Val.Comment), false)
	if err != nil {
		message := fmt.Sprintf("Collector:handleDisabledHook:Database.AbortAttempt: %x", err)
		logging.Error(message)
		panic(message)
	}
}




func NewCollector(runnerParams component.RunnerParams, database storeInterface.IStoreProvider,  
	client clientInterface.IHTTPClient) *Collector {

		return &Collector{
			WebhookRunner: *component.NewWebhookRunner(runnerParams, database,client, commons.HookChannel, commons.AttemptChannel),
		}

}