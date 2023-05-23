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

func (s *Server) Stargate(stream api.StargateService_StargateServer) error {
	ctx := stream.Context()
	organizationID, stackID, err := orgaAndStackIDFromIncomingContext(ctx)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "cannot get organization and stack id from contex metadata: %v", err)
	}

	s.logger.WithFields(map[string]any{
		"organization_id": organizationID,
		"stack_id":        stackID,
	}).Infof("new stargate connection")
	defer s.logger.WithFields(map[string]any{
		"organization_id": organizationID,
		"stack_id":        stackID,
	}).Infof("stargate connection closed")

	waitingResponses := sync.Map{}
	subject := utils.GetNatsSubject(organizationID, stackID)
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
