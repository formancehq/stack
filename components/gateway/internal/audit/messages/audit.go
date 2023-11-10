package messages

import (
	"net/http"
	"time"

	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/google/uuid"
)

const (
	EventVersion   = "v1"
	EventApp       = "gateway"
	EventTypeAudit = "AUDIT"
)

type HttpRequest struct {
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Host   string      `json:"host"`
	Header http.Header `json:"header"`
	Body   string      `json:"body,omitempty"`
}

type HttpResponse struct {
	StatusCode int         `json:"status_code"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body,omitempty"`
}

func NewHttpResponse(
	statusCode int,
	headers http.Header,
	body string,
) HttpResponse {
	return HttpResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       body,
	}
}

type Payload struct {
	ID       string       `json:"id"`
	Request  HttpRequest  `json:"request"`
	Response HttpResponse `json:"response"`
}

func NewAuditMessagePayload(
	request HttpRequest,
	response HttpResponse,
) publish.EventMessage {

	payload := Payload{
		ID:       uuid.New().String(),
		Request:  request,
		Response: response,
	}

	return publish.EventMessage{
		Date:    time.Now().UTC(),
		App:     EventApp,
		Version: EventVersion,
		Type:    EventTypeAudit,
		Payload: payload,
	}
}
