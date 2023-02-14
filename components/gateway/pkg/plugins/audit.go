package plugins

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/formancehq/stack/components/gateway/internal/audit/messages"
	"github.com/formancehq/stack/libs/go-libs/logging/logginglogrus"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(Audit{})
	httpcaddyfile.RegisterHandlerDirective("audit", parseAuditCaddyfile)
}

type Audit struct {
	logger    *logrus.Logger    `json:"-"`
	bufPool   *sync.Pool        `json:"-"`
	publisher message.Publisher `json:"-"`

	PublisherKafkaBroker           string   `json:"publisher_kafka_broker,omitempty"`
	PublisherKafkaTopics           []string `json:"publisher_kafka_topics,omitempty"`
	PublisherKafkaEnabled          bool     `json:"publisher_kafka_enabled,omitempty"`
	PublisherKafkaTLSEnabled       bool     `json:"publisher_kafka_tls_enabled,omitempty"`
	PublisherKafkaSASLEnabled      bool     `json:"publisher_kafka_sasl_enabled,omitempty"`
	PublisherKafkaSASLUsername     string   `json:"publisher_kafka_sasl_username,omitempty"`
	PublisherKafkaSASLPassword     string   `json:"publisher_kafka_sasl_password,omitempty"`
	PublisherKafkaSASLMechanism    string   `json:"publisher_kafka_sasl_mechanism,omitempty"`
	PublisherKafkaSASLScramSHASize int      `json:"publisher_kafka_sasl_scram_sha_size,omitempty"`
}

// Implements the caddy.Module interface.
func (Audit) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.audit",
		New: func() caddy.Module { return new(Audit) },
	}
}

func parseAuditCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	a := new(Audit)
	for h.Next() {
		for h.NextBlock(0) {
			key := h.Val()
			switch key {
			case "publisher_kafka_enabled":
				var err error
				a.PublisherKafkaEnabled, err = parseBool(h.Dispenser)
				if err != nil {
					return nil, h.Errf("failed to parse publisher_kafka_enabled: %v", err)
				}

			case "publisher_kafka_broker":
				if !h.AllArgs(&a.PublisherKafkaBroker) {
					return nil, h.Errf("expected one string value for kafka_broker")
				}

			case "publisher_kafka_topics":
				a.PublisherKafkaTopics = h.RemainingArgs()

			case "publisher_kafka_tls_enabled":
				var err error
				a.PublisherKafkaTLSEnabled, err = parseBool(h.Dispenser)
				if err != nil {
					return nil, h.Errf("failed to parse publisher_kafka_tls_enabled: %v", err)
				}

			case "publisher_kafka_sasl_enabled":
				var err error
				a.PublisherKafkaSASLEnabled, err = parseBool(h.Dispenser)
				if err != nil {
					return nil, h.Errf("failed to parse publisher_kafka_sasl_enabled: %v", err)
				}

			case "publisher_kafka_sasl_username":
				if !h.AllArgs(&a.PublisherKafkaSASLUsername) {
					return nil, h.Errf("expected one string value for publisher kafka sasl username")
				}

			case "publisher_kafka_sasl_password":
				if !h.AllArgs(&a.PublisherKafkaSASLPassword) {
					return nil, h.Errf("expected one string value for publisher kafka sasl password")
				}

			case "publisher_kafka_sasl_mechanism":
				if !h.AllArgs(&a.PublisherKafkaSASLMechanism) {
					return nil, h.Errf("expected one string value for publisher kafka sasl mechanism")
				}

			case "publisher_kafka_sasl_scram_sha_size":
				var publisherKafkaSaslScramShaSize string
				if !h.AllArgs(&publisherKafkaSaslScramShaSize) {
					return nil, h.Errf("expected one boolean value")
				}

				res, err := strconv.ParseInt(publisherKafkaSaslScramShaSize, 10, 32)
				if err != nil {
					return nil, h.Errf("failed to parse publisher_kafka_sasl_scram_sha_size: %v", err)
				}
				a.PublisherKafkaSASLScramSHASize = int(res)

			default:
				return nil, h.Errf("unrecognized option: %s", key)
			}
		}
	}

	return a, nil
}

func parseBool(d *caddyfile.Dispenser) (bool, error) {
	var b string
	if !d.AllArgs(&b) {
		return false, d.Errf("expected one boolean value")
	}

	res, err := strconv.ParseBool(b)
	if err != nil {
		return false, d.Errf("expected boolean value")
	}

	return res, nil
}

// Implements the caddy.Provisioner interface.
func (a *Audit) Provision(ctx caddy.Context) error {
	// We have to instantiate a different logger here instead of using
	// caddy's logger because we need a different interface for the watermill
	// logger.
	a.logger = logrus.New()
	a.bufPool = &sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}

	var err error
	a.publisher, err = publish.NewKafkaPublisher(
		publish.NewWatermillLoggerAdapter(logginglogrus.New(a.logger)),
		publish.NewSaramaConfig(
			publish.ClientId("gateway"),
			sarama.V1_0_0_0,
			publish.BuildSaramaOption(
				a.PublisherKafkaTLSEnabled,
				a.PublisherKafkaSASLEnabled,
				a.PublisherKafkaSASLUsername,
				a.PublisherKafkaSASLPassword,
				a.PublisherKafkaSASLMechanism,
				a.PublisherKafkaSASLScramSHASize,
			)...),
		kafka.DefaultMarshaler{},
		a.PublisherKafkaBroker,
	)
	if err != nil {
		a.logger.Error("failed to create kafka publisher", zap.Error(err))
		return err
	}

	return nil
}

func (a Audit) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	request, err := httputil.DumpRequest(r, true)
	if err != nil {
		a.logger.Error("unable to dump request: %v", err)
		return err
	}

	buf := a.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer a.bufPool.Put(buf)

	rww := NewResponseWriterWrapper(w, buf)
	if err := next.ServeHTTP(rww, r); err != nil {
		return err
	}

	response, err := json.Marshal(messages.NewAuditResponseMessage(
		*rww.statusCode,
		rww.Header(),
		rww.body.Bytes(),
	))
	if err != nil {
		a.logger.Error("failed to create audit message: %v", err)
		return err
	}

	if err := a.publisher.Publish(
		messages.TopicAudit,
		publish.NewMessage(
			r.Context(),
			messages.NewAuditMessagePayload(
				request,
				response,
			),
		),
	); err != nil {
		a.logger.Error("failed to publish audit message: %v", err)
	}

	return nil
}

// Interface Guards
var (
	_ caddy.Provisioner           = (*Audit)(nil)
	_ caddy.Module                = (*Audit)(nil)
	_ caddyhttp.MiddlewareHandler = (*Audit)(nil)
)

//------------------------------------------------------------------------------

// ResponseWriterWrapper is a wrapper for the http.ResponseWriter, it captures
// the response body and status code to be used in the audit log.
type ResponseWriterWrapper struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode *int
}

// NewResponseWriterWrapper static function creates a wrapper for the
// http.ResponseWriter
func NewResponseWriterWrapper(w http.ResponseWriter, buf *bytes.Buffer) ResponseWriterWrapper {
	statusCode := 200
	return ResponseWriterWrapper{
		ResponseWriter: w,
		body:           buf,
		statusCode:     &statusCode, // Default status code
	}
}

func (rww ResponseWriterWrapper) Write(buf []byte) (int, error) {
	rww.body.Write(buf)
	return rww.ResponseWriter.Write(buf)
}

// Header function overwrites the http.ResponseWriter Header() function
func (rww ResponseWriterWrapper) Header() http.Header {
	return rww.ResponseWriter.Header()

}

// WriteHeader function overwrites the http.ResponseWriter WriteHeader() function
func (rww ResponseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.statusCode) = statusCode
	rww.ResponseWriter.WriteHeader(statusCode)
}
