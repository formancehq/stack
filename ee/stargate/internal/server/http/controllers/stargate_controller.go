package controllers

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	service "github.com/formancehq/stack/components/stargate/internal/api"
	"github.com/formancehq/stack/components/stargate/internal/opentelemetry"
	"github.com/formancehq/stack/components/stargate/internal/utils"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
)

var IsLetter = regexp.MustCompile(`^[a-z]+$`).MatchString

func (s *StargateController) HandleCalls(w http.ResponseWriter, r *http.Request) {
	s.logger.Debugf("[HTTP] received call: %s", r.URL.Path)

	ctx := r.Context()
	organizationID, stackID, err := getOrganizationAndStackID(r)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	var status int
	attrs := []attribute.KeyValue{
		attribute.String("organization_id", organizationID),
		attribute.String("stack_id", stackID),
		attribute.String("path", r.URL.Path),
	}
	defer func() {
		attrs = append(attrs, attribute.Int("status", status))
		s.metricsRegistry.ReceivedHTTPCallByPath().Add(ctx, 1, metric.WithAttributes(attrs...))
	}()

	msg, err := requestToProto(r)
	if err != nil {
		status = ResponseError(w, r, errors.Wrapf(err, "failed to parse request"))
		return
	}

	buf, err := proto.Marshal(msg)
	if err != nil {
		status = ResponseError(w, r, errors.Wrapf(err, "failed to marshal message"))
		return
	}

	subject := utils.GetNatsSubject(organizationID, stackID)
	s.logger.Debugf("[HTTP] sending message to %s with path: %s", subject, r.URL.Path)
	resp, err := s.natsConn.Request(subject, buf, s.config.natsRequestTimeout)
	if err != nil {
		s.logger.Errorf("[HTTP] error sending message to %s with path: %s: %v", subject, r.URL.Path, err)
		status = ResponseError(w, r, ErrNoResponders)
		return
	}

	var response service.StargateClientMessage
	if err = proto.Unmarshal(resp.Data, &response); err != nil {
		status = ResponseError(w, r, errors.Wrapf(err, "failed to unmarshal response"))
		return
	}

	s.logger.Debugf("[HTTP] received response for %s with path: %s", subject, r.URL.Path)

	switch ev := response.Event.(type) {
	case *service.StargateClientMessage_ApiCallResponse:
		for k, v := range ev.ApiCallResponse.Headers {
			for _, vv := range v.Values {
				w.Header().Add(k, vv)
			}
		}

		status = int(ev.ApiCallResponse.StatusCode)
		api.WriteResponse(w, int(ev.ApiCallResponse.StatusCode), ev.ApiCallResponse.Body)
		return
	default:
		status = ResponseError(w, r, errors.Wrapf(err, "invalid response from client"))
		return
	}
}

func getOrganizationAndStackID(r *http.Request) (string, string, error) {
	paths := strings.Split(r.URL.Path, "/")
	if len(paths) < 3 {
		return "", "", errors.Wrapf(ErrValidation, "invalid path, missing organizationID and stackID")
	}

	organizationID := paths[1]
	if len(organizationID) != 12 || !IsLetter(organizationID) {
		return "", "", errors.Wrapf(ErrValidation, "invalid organizationID")
	}

	stackID := paths[2]
	if len(stackID) != 4 || !IsLetter(stackID) {
		return "", "", errors.Wrapf(ErrValidation, "invalid stackID")
	}

	paths = append(paths[:1], paths[3:]...)

	r.URL.Path = strings.Join(paths, "/")
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

	carrier := propagation.MapCarrier{}
	opentelemetry.Propagator.Inject(r.Context(), carrier)

	return &service.StargateServerMessage{
		Event: &service.StargateServerMessage_ApiCall{
			ApiCall: &service.StargateServerMessage_APICall{
				Method:      r.Method,
				Path:        r.URL.Path,
				Query:       urlQuery,
				Body:        body,
				Headers:     headers,
				OtlpContext: carrier,
			},
		},
	}, nil
}
