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
	TimeOut int //ms
	MaxCall int 
	MaxRetry int 
	DelayPull int //second

}

func DefaultRunnerParams() RunnerParams{
	return RunnerParams{
		TimeOut:2000,
		MaxCall:20, 
		MaxRetry: 60,
		DelayPull: 1,
	}
}


type WebhookRunner struct{
	RunnerParams RunnerParams
	
	Queue *sync.Queue
	StopChan chan struct{}
	EventChan chan commons.Event

	State *commons.State
	
	Database storeInterface.IStoreProvider
	Client clientInterface.IHTTPClient
} 

func (wr *WebhookRunner) Stop(){
	wr.StopChan <- struct{}{}
}

func (wr *WebhookRunner) Run(){

}

func (wr *WebhookRunner) StartHandleEventFromDatabase(){
	go wr.State.RoutineEvent(wr.StopChan, wr.EventChan, wr.HandleEvent)
}

func (wr *WebhookRunner) StartRetrySaveToDatabase(){
	ticker := time.NewTicker(10*time.Second)
	go wr.State.RoutineTime(wr.StopChan, ticker, wr.ToSaveProcess)
}
 
func (wr *WebhookRunner) HandleEvent(e commons.Event){
	eventType := commons.TypeFromEvent(e) 
	switch eventType {

	case commons.NewHookType :
		hook, err := wr.Database.GetHook(e.ID)
		if (err!=nil){
			//TODO(CriticPolitic)
			logging.Errorf("WebhookRunner:NewHookEvent:Database.GetHook() : %x", err)
		}else {
			wr.State.AddNewHook(&hook)
		}
		
	case commons.ChangeHookStatusType:
		switch e.Value {
		case commons.EnableStatus:
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
		if (err!=nil){
			//TODO(CriticPolitic)
			logging.Errorf("WebhookRunner:NewWaitingAttemptEvent:Database.GetAttempt() : %x", err)
		}
		wr.State.AddNewAttempt(&attempt)
	case commons.FlushWaitingAttemptType:
		wr.State.FlushAttempt(e.ID)
	case commons.FlushWaitingAttemptsType:
		wr.State.FlushAttempts()
	case commons.AbortWaitingAttemptType:
		wr.State.AbortAttempt(e.ID)
	default:
		//TODO(Translate)
		logging.Error(fmt.Sprintf("Un évenement non traité est survenue : %x", e))
	}
}

func (wr *WebhookRunner) ToSaveProcess(){
	if(wr.State.ToSaveAttempts.Size()==0){return}
	
	toSave := wr.State.ToSaveAttempts.Empty()
	for _,sAttempt := range *toSave.Val {
		_, err := wr.Database.ChangeAttemptStatus(sAttempt.Val.ID, sAttempt.Val.Status, string(sAttempt.Val.Comment))
		if err != nil {
			wr.State.ToSaveAttempts.Add(sAttempt)
			//TODO(Translate)
			logging.Error(fmt.Sprintf("Une erreur survient pendant une requête vers la DB: %x", err))
		}
	}
}

func (wr *WebhookRunner) HandleRequest(ctx context.Context, sAttempt *commons.SharedAttempt, sHook *commons.SharedHook) (int, error){
	wr.Queue.Lock()
	timeOut := time.Duration(wr.RunnerParams.TimeOut)*time.Millisecond
	requestCtx, cancel := context.WithTimeout(ctx, timeOut)
	statusCode, err := wr.Client.Call(requestCtx, sHook.Val, sAttempt.Val, false)
	cancel()
	wr.Queue.Unlock()
	return statusCode, err
}

func (wr *WebhookRunner) SendEvent(t commons.EventType, attempt *commons.Attempt, hook *commons.Hook) error{
	event, err := commons.EventFromType(t, attempt, hook)
	if err != nil {return err}
	
	return wr.Database.NotifyUpdate(event) 
}


func NewWebhookRunner(runnerParams RunnerParams, 
					eventChan chan commons.Event,
					database storeInterface.IStoreProvider, 
					client clientInterface.IHTTPClient,
					) *WebhookRunner{

	return &WebhookRunner{
		RunnerParams: runnerParams,

		Queue: sync.NewQueue(runnerParams.MaxCall),
		StopChan: make(chan struct{},0),
		EventChan: eventChan,

		State: commons.NewState(),

		Database: database,
		Client: client,
		
	}

}