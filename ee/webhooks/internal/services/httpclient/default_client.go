package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"time"

	"github.com/formancehq/webhooks/internal/models"
	"github.com/formancehq/webhooks/internal/utils/security"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DefaultHttpClient struct {
	httpClient *http.Client
}

var Tracer = otel.Tracer("webhook")

func (dc DefaultHttpClient) Call(context context.Context, hook *models.Hook, attempt *models.Attempt, isTest bool) (int, error) {

	parentSpan := trace.SpanFromContext(context)
	ctx, span := Tracer.Start(context,
		"Webhook::Hook.Service::processAttempt::DefaultHttpClient.Call",
		trace.WithLinks(trace.Link{
			SpanContext: parentSpan.SpanContext(),
		}),
		trace.WithAttributes(
			attribute.String("attempt-id", attempt.ID),
			attribute.String("hook-name", attempt.HookName),
			attribute.String("hook-endpoint", attempt.HookEndpoint),
			attribute.String("attempt-event", attempt.Event),
			attribute.String("attempt-payload", attempt.Payload),
			attribute.Int("attempt-nbTry", attempt.NbTry),
		),
	)

	defer span.End()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		hook.Endpoint,
		bytes.NewBuffer([]byte(attempt.Payload)))

	if err != nil {
		span.RecordError(err)
		return 0, err
	}

	ts := time.Now().UTC()
	timestamp := ts.Unix()
	signature, err := security.Sign(attempt.ID, timestamp, hook.Secret, []byte(attempt.Payload))

	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", "formance-webhooks/v2")
	req.Header.Set("formance-webhook-id", attempt.ID)
	req.Header.Set("formance-webhook-timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("formance-webhook-signature", signature)
	req.Header.Set("formance-webhook-test", fmt.Sprintf("%v", isTest))
	resp, err := dc.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)

		return 503, nil

	}

	span.SetAttributes(attribute.Int("attempt-statusCode", resp.StatusCode))

	return resp.StatusCode, nil

}

func NewDefaultHttpClient(httpClient *http.Client) DefaultHttpClient {
	return DefaultHttpClient{
		httpClient: httpClient,
	}
}
