package messages

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	Identity string       `json:"identity"`
	Request  HttpRequest  `json:"request"`
	Response HttpResponse `json:"response"`
}

func NewAuditMessagePayload(
	logger *zap.Logger,
	request HttpRequest,
	response HttpResponse,
) publish.EventMessage {
	identity := ""

	if request.Header != nil {
		if authorizationHeader := request.Header.Get("Authorization"); authorizationHeader != "" && strings.HasPrefix(strings.ToLower(authorizationHeader), "bearer ") {

			tokenString := strings.Replace(strings.Replace(authorizationHeader, "Bearer ", "", 1), "bearer ", "", 1)
			token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
			if err != nil {
				logger.Error(fmt.Sprintf("error for Parse %s", err))
			}
			if token != nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					identity = fmt.Sprint(claims["sub"])
				} else {
					logger.Error(fmt.Sprintf("error get claims JWT token: %s", err))
				}
			}
		}

		request.Header.Del("Authorization")
	}

	if request.Path == "/api/auth/oauth/token" {
		response.Body = ""
	}

	payload := Payload{
		ID:       uuid.New().String(),
		Identity: identity,
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
