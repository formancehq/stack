package webhookworker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/stack/libs/go-libs/contextutil"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/formancehq/webhooks/internal/commons"
	component "github.com/formancehq/webhooks/internal/components/commons"
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)

var Tracer = otel.Tracer("WebhookWorker")

type Worker struct {
	component.WebhookRunner
}

func (w *Worker) Init() {
	w.StartHandleFreshLogs()

	hooks, err := w.Database.LoadHooks()
	if err != nil {
		logging.Error(err.Error())
		os.Exit(1)
	}
	w.State.LoadHooks(hooks)
}

func (w *Worker) HandleMessage(msg *message.Message) error {

	var ev *publish.EventMessage
	span, ev, err := publish.UnmarshalMessage(msg)
	if err != nil {
		logging.FromContext(msg.Context()).Error(err.Error())
		return err
	}

	ctx, span := Tracer.Start(msg.Context(), "WebhookWorker:HandleMessage",
		trace.WithLinks(trace.Link{
			SpanContext: span.SpanContext(),
		}),
		trace.WithAttributes(
			attribute.String("event-id", msg.UUID),
			attribute.Bool("duplicate", false),
			attribute.String("event-type", ev.Type),
			attribute.String("event-payload", string(msg.Payload)),
		),
	)
	defer span.End()
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
	}()
	traceCtx, _ := contextutil.Detached(ctx)

	event := strings.ToLower(ev.Type)
	eventApp := strings.ToLower(ev.App)

	if eventApp != "" {
		event = strings.Join([]string{eventApp, event}, ".")
	}

	triggedSHooks := w.State.ActiveHooksByEvent.Get(event)
	if triggedSHooks == nil || triggedSHooks.Size() == 0 {
		return nil
	}

	payload, err := json.Marshal(ev)
	if err != nil {
		logging.FromContext(traceCtx).Error(err)
		return err
	}

	var globalError error = nil

	triggedSHooks.AsyncApply(w.HandlerTriggedHookFactory(traceCtx, event, string(payload), globalError))

	return globalError

}

func (w *Worker) HandlerTriggedHookFactory(ctx context.Context, event string, payload string, globalError error) func(*commons.SharedHook, *sync.WaitGroup) {

	return func(sHook *commons.SharedHook, wg *sync.WaitGroup) {

		defer wg.Done()

		sAttempt := commons.NewSharedAttempt(sHook.Val.ID, sHook.Val.Name, sHook.Val.Endpoint, event, string(payload))
		hook := sHook.Val
		attempt := sAttempt.Val
		statusCode, err := w.HandleRequest(ctx, sAttempt, sHook)

		if err != nil {
			message := fmt.Sprintf("Worker:triggedSHooks.AsyncApply() - HandleTriggedHookFactory() - func(sHook *commons.SharedHook,wg *sync.WaitGroup) - w.HandleRequest - Something Went wrong while trying to make http request: %x", err)
			logging.Error(message)
			panic(message)

		}

		w.HandleResponse(statusCode, attempt, hook)

	}
}

func (w *Worker) HandleResponse(statusCode int, attempt *commons.Attempt, hook *commons.Hook) error {
	attempt.LastHttpStatusCode = statusCode
	attempt.NbTry += 1
	var err error

	if commons.IsHTTPRequestSuccess(statusCode) {
		commons.SetSuccesStatus(attempt)
		err = w.Database.SaveAttempt(*attempt, true)
	}

	if !hook.Retry && !attempt.IsSuccess() {
		commons.SetAbortNoRetryModeStatus(attempt)
		err = w.Database.SaveAttempt(*attempt, false)
	}

	return err

}

func NewWorker(runnerParams component.RunnerParams, database storeInterface.IStoreProvider, client clientInterface.IHTTPClient) *Worker {

	return &Worker{
		WebhookRunner: *component.NewWebhookRunner(runnerParams, database, client, commons.HookChannel),
	}
}
