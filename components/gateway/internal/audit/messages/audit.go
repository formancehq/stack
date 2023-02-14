package messages

import (
	"net/http"
	"time"
)

type auditResponseMessage struct {
	StatusCode int         `json:"status_code"`
	Headers    http.Header `json:"headers"`
	Body       []byte      `json:"body"`
}

func NewAuditResponseMessage(
	statusCode int,
	headers http.Header,
	body []byte,
) auditResponseMessage {
	return auditResponseMessage{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       body,
	}
}

type auditMessagePayload struct {
	Request  []byte `json:"request"`
	Response []byte `json:"response"`
}

func NewAuditMessagePayload(
	request []byte,
	response []byte,
) EventMessage {
	payload := auditMessagePayload{
		Request:  request,
		Response: response,
	}

	return EventMessage{
		Date:    time.Now().UTC(),
		App:     EventApp,
		Version: EventVersion,
		Type:    EventTypeAudit,
		Payload: payload,
	}
}
