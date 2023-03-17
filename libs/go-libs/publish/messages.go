package publish

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func NewMessage(ctx context.Context, m proto.Message) *message.Message {
	data, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}
	msg := message.NewMessage(uuid.NewString(), data)
	msg.SetContext(ctx)
	return msg
}
