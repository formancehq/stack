package controllers

import (
	"io"
	"net/http"
	"strings"

	service "github.com/formancehq/stack/components/stargate/internal/api"
	"github.com/formancehq/stack/components/stargate/internal/utils"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/proto"
)

func (s *StargateController) HandleCalls(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organizationID, stackID, err := getOrganizationAndStackID(r)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("organization_id", organizationID),
		attribute.String("stack_id", stackID),
		attribute.String("path", r.URL.Path),
	}
	s.metricsRegistry.ReceivedHTTPCallByPath().Add(ctx, 1, attrs...)

	msg, err := requestToProto(r)
	if err != nil {
		ResponseError(w, r, errors.Wrapf(err, "failed to parse request"))
		return
	}

	buf, err := proto.Marshal(msg)
	if err != nil {
		ResponseError(w, r, errors.Wrapf(err, "failed to marshal message"))
		return
	}

	subject := utils.GetNatsSubject(organizationID, stackID)
	// requestCtx, cancel := context.WithTimeout(ctx, s.config.natsRequestTimeout)
	// defer cancel()
	resp, err := s.natsConn.Request(subject, buf, s.config.natsRequestTimeout)
	if err != nil {
		ResponseError(w, r, errors.Wrapf(err, "failed to send message"))
		return
	}

	var response service.StargateClientMessage
	if err = proto.Unmarshal(resp.Data, &response); err != nil {
		ResponseError(w, r, errors.Wrapf(err, "failed to unmarshal response"))
		return
	}

	switch ev := response.Event.(type) {
	case *service.StargateClientMessage_ApiCallResponse:
		for k, v := range ev.ApiCallResponse.Headers {
			for _, vv := range v.Values {
				w.Header().Add(k, vv)
			}
		}

		api.WriteResponse(w, int(ev.ApiCallResponse.StatusCode), ev.ApiCallResponse.Body)
		return
	default:
		ResponseError(w, r, errors.Wrapf(err, "invalid response from client"))
		return
	}
}

func getOrganizationAndStackID(r *http.Request) (string, string, error) {
	authHeader := r.Header.Get("authorization")
	if authHeader == "" {
		return "", "", errors.Wrapf(ErrValidation, "missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// The token will be verified by the client directly, we just need here to
	// parse it and check if the organizationID and stackID are present.
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", "", errors.Wrapf(ErrValidation, "invalid authorization header")
	}

	mapClaims := token.Claims.(jwt.MapClaims)
	organizationID, ok := mapClaims["organization_id"].(string)
	if !ok {
		return "", "", errors.Wrapf(ErrValidation, "organization_id not found in token")
	}

	stackID, ok := mapClaims["stack_id"].(string)
	if !ok {
		return "", "", errors.Wrapf(ErrValidation, "stack_id not found in token")
	}

	return organizationID, stackID, nil
}

func requestToProto(r *http.Request) (*service.StargateServerMessage, error) {
	urlQuery := make(map[string]*service.Values)
	for k, v := range r.URL.Query() {
		urlQuery[k] = &service.Values{Values: v}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]*service.Values)
	for k, v := range r.Header {
		headers[k] = &service.Values{Values: v}
	}

	return &service.StargateServerMessage{
		Event: &service.StargateServerMessage_ApiCall{
			ApiCall: &service.StargateServerMessage_APICall{
				Method:  r.Method,
				Path:    r.URL.Path,
				Query:   urlQuery,
				Body:    body,
				Headers: headers,
			},
		},
	}, nil
}
