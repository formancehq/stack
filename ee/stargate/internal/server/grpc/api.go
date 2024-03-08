package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/formancehq/stack/components/stargate/internal/generated"
	"github.com/formancehq/stack/components/stargate/internal/server/grpc/metrics"
	"github.com/formancehq/stack/components/stargate/internal/utils"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	generated.UnimplementedStargateServiceServer

	logger          logging.Logger
	natsConn        *nats.Conn
	metricsRegistry metrics.MetricsRegistry
}

func NewServer(
	logger logging.Logger,
	natsConn *nats.Conn,
	metricsRegistry metrics.MetricsRegistry,
) *Server {
	return &Server{
		logger:          logger,
		natsConn:        natsConn,
		metricsRegistry: metricsRegistry,
	}
}

type waitingResponse struct {
	msg    *nats.Msg
	sendAt time.Time
}

type waitingPingResponse struct {
	resp chan struct{}
}

func (s *Server) Stargate(stream generated.StargateService_StargateServer) error {
	ctx := stream.Context()
	organizationID, stackID, err := orgaAndStackIDFromIncomingContext(ctx)
	if err != nil {
		s.metricsRegistry.StreamErrors().Add(ctx, 1, metric.WithAttributes([]attribute.KeyValue{
			attribute.Int("code", int(codes.InvalidArgument)),
		}...))
		return status.Errorf(codes.InvalidArgument, "cannot get organization and stack id from contex metadata: %v", err)
	}

	logger := s.logger.WithFields(map[string]any{
		"organization_id": organizationID,
		"stack_id":        stackID,
	})
	metrics.ClientsConnected.Add(1)
	defer metrics.ClientsConnected.Add(-1)

	logger.Infof("[GRPC] new stargate connection")
	defer logger.Infof("[GRPC] stargate connection closed")

	waitingResponses := sync.Map{}
	waitingPingResponses := sync.Map{}

	subject := utils.GetNatsSubject(organizationID, stackID)
	logger.Debugf("[GRPC] subscribing to nats subject %s", subject)
	sub, err := s.natsConn.QueueSubscribeSync(subject, subject)
	if err != nil {
		s.metricsRegistry.StreamErrors().Add(ctx, 1, metric.WithAttributes([]attribute.KeyValue{
			attribute.Int("code", int(codes.Internal)),
		}...))
		return status.Errorf(codes.Internal, "cannot subscribe to nats subject")
	}
	defer sub.Unsubscribe()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		for {
			msg, err := sub.NextMsgWithContext(ctx)
			if err != nil {
				if err == context.Canceled {
					return nil
				} else {
					return err
				}
			}

			var request generated.StargateServerMessage
			if err := proto.Unmarshal(msg.Data, &request); err != nil {
				return err
			}

			switch ev := request.Event.(type) {
			case *generated.StargateServerMessage_ApiCall:
				logger.WithFields(map[string]any{
					"path": ev.ApiCall.Path,
				}).Debug("[GRPC] stream api call request received")
			default:
			}

			correlationID := uuid.New().String()

			waitingResponses.Store(correlationID, waitingResponse{
				msg:    msg,
				sendAt: time.Now(),
			})
			request.CorrelationId = correlationID

			if err := stream.Send(&request); err != nil {
				return err
			}
		}
	})

	eg.Go(func() error {
		for {
			in, err := stream.Recv()
			if err != nil {
				if err == io.EOF || status.Code(err) == codes.Canceled {
					return nil
				}

				return err
			}

			switch in.Event.(type) {
			case *generated.StargateClientMessage_ApiCallResponse:
				logger.Debugf("[GRPC] stream api call response received")

				entry, ok := waitingResponses.LoadAndDelete(in.CorrelationId)
				if !ok {
					s.metricsRegistry.CorrelationIDNotFound().Add(ctx, 1)
					continue
				}
				wr := entry.(waitingResponse)

				s.metricsRegistry.GRPCLatencies().Record(ctx, time.Since(wr.sendAt).Milliseconds(), metric.WithAttributes([]attribute.KeyValue{
					attribute.String("name", fmt.Sprintf("%s.%s", organizationID, stackID)),
				}...))

				data, err := proto.Marshal(in)
				if err != nil {
					return err
				}

				if err := wr.msg.Respond(data); err != nil {
					return err
				}
			case *generated.StargateClientMessage_Pong_:
				entry, ok := waitingPingResponses.LoadAndDelete(in.CorrelationId)
				if !ok {
					continue
				}

				wpr := entry.(waitingPingResponse)
				close(wpr.resp)
			}
		}
	})

	// We have to implement a ping/pong system to detect dead connections.
	// We cannot use grpc keepalive to do that because AWS ALB does not support
	// raw HTTP/2 frames.
	// c.f.: https://stackoverflow.com/questions/66818645/http2-ping-frames-over-aws-alb-grpc-keepalive-ping
	eg.Go(func() error {
		for {
			select {
			case <-time.After(10 * time.Second):
				correlationID := uuid.New().String()
				resp := make(chan struct{})
				waitingPingResponses.Store(correlationID, waitingPingResponse{
					resp: resp,
				})

				if err := stream.Send(&generated.StargateServerMessage{
					CorrelationId: correlationID,
					Event:         &generated.StargateServerMessage_Ping_{Ping: &generated.StargateServerMessage_Ping{}},
				}); err != nil {
					return err
				}

				select {
				case <-time.After(10 * time.Second):
					logger.Debugf("[GRPC] ping timeout")
					return status.Errorf(codes.DeadlineExceeded, "ping timeout")
				case <-resp:
					// Pong received, do nothing
				}
			case <-ctx.Done():
				return nil
			}
		}
	})

	if err := eg.Wait(); err != nil {
		s.metricsRegistry.StreamErrors().Add(ctx, 1, metric.WithAttributes([]attribute.KeyValue{
			attribute.Int("code", int(codes.Internal)),
		}...))
		return status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return nil
}

func orgaAndStackIDFromIncomingContext(ctx context.Context) (string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", errors.New("no metadata from incoming context")
	}

	organizationID := md.Get("organization-id")
	if len(organizationID) == 0 {
		return "", "", errors.New("no organization-id in metadata")
	}

	stackID := md.Get("stack-id")
	if len(stackID) == 0 {
		return "", "", errors.New("no stack-id in metadata")
	}

	return organizationID[0], stackID[0], nil
}
