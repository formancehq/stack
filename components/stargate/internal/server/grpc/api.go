package grpc

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/formancehq/stack/components/stargate/internal/api"
	"github.com/formancehq/stack/components/stargate/internal/server/grpc/opentelemetry"
	"github.com/formancehq/stack/components/stargate/internal/utils"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	api.UnimplementedStargateServiceServer

	logger          logging.Logger
	natsConn        *nats.Conn
	metricsRegistry opentelemetry.MetricsRegistry
}

func NewServer(
	logger logging.Logger,
	natsConn *nats.Conn,
	metricsRegistry opentelemetry.MetricsRegistry,
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

func (s *Server) Stargate(stream api.StargateService_StargateServer) error {
	ctx := stream.Context()
	organizationID, stackID, err := orgaAndStackIDFromIncomingContext(ctx)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "cannot get organization and stack id from contex metadata: %v", err)
	}

	logger := s.logger.WithFields(map[string]any{
		"organization_id": organizationID,
		"stack_id":        stackID,
	})

	logger.Infof("[GRPC] new stargate connection")
	defer logger.Infof("[GRPC] stargate connection closed")

	waitingResponses := sync.Map{}
	waitingPingResponses := sync.Map{}

	subject := utils.GetNatsSubject(organizationID, stackID)
	logger.Debugf("[GRPC] subscribing to nats subject %s", subject)
	sub, err := s.natsConn.QueueSubscribeSync(subject, subject)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot subscribe to nats subject")
	}

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

			var request api.StargateServerMessage
			if err := proto.Unmarshal(msg.Data, &request); err != nil {
				return err
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
			case *api.StargateClientMessage_ApiCallResponse:
				logger.Debugf("[GRPC] stream api call response received")

				entry, ok := waitingResponses.LoadAndDelete(in.CorrelationId)
				if !ok {
					s.metricsRegistry.CorrelationIDNotFound().Add(ctx, 1)
					continue
				}
				wr := entry.(waitingResponse)

				s.metricsRegistry.GRPCLatencies().Record(ctx, time.Since(wr.sendAt).Milliseconds())

				data, err := proto.Marshal(in)
				if err != nil {
					return err
				}

				if err := wr.msg.Respond(data); err != nil {
					return err
				}
			case *api.StargateClientMessage_Pong_:
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

				if err := stream.Send(&api.StargateServerMessage{
					CorrelationId: correlationID,
					Event:         &api.StargateServerMessage_Ping_{Ping: &api.StargateServerMessage_Ping{}},
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
		// TODO(polo): should we expose the error here ?
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
