package messages

import (
	"net/http"
	"net/url"
	"time"

	"github.com/formancehq/stack/libs/go-libs/publish"
)

const (
	EventVersion   = "v1"
	EventApp       = "gateway"
	EventTypeAudit = "AUDIT"
)

type HttpRequest struct {
	Method string      `json:"method"`
	Url    *url.URL    `json:"path"`
	Host   string      `json:"host"`
	Header http.Header `json:"header"`
	Body   string      `json:"body"`
}

type HttpResponse struct {
	StatusCode int         `json:"status_code"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body"`
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
	Request  HttpRequest  `json:"request"`
	Response HttpResponse `json:"response"`
}

func NewAuditMessagePayload(
	request HttpRequest,
	response HttpResponse,
) publish.EventMessage {

	payload := Payload{
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
