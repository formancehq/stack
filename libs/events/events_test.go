package events

import (
	"fmt"
	"testing"
	"time"

	"github.com/formancehq/stack/libs/events/payments"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTtoto(t *testing.T) {
	event := &Event{
		CreatedAt: timestamppb.New(time.Now()),
		Event: &Event_ResetConnector{
			ResetConnector: &payments.ResetConnector{
				CreatedAt: timestamppb.New(time.Now()),
				Provider:  payments.ConnectorProvider_CONNECTOR_PROVIDER_BANKING_CIRCLE,
			},
		},
	}

	fmt.Println(event)

	b, err := proto.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	event_test := &Event{}
	if err := proto.Unmarshal(b, event_test); err != nil {
		t.Fatal(err)
	}

	fmt.Println(event_test)

	// protojson.MarshalOptions{
	// 	NoUnkeyedLiterals: pragma.NoUnkeyedLiterals{},
	// 	Multiline:         false,
	// 	Indent:            "",
	// 	AllowPartial:      false,
	// 	UseProtoNames:     false,
	// 	UseEnumNumbers:    false,
	// 	EmitUnpopulated:   false,
	// 	Resolver:          nil,
	// }

	b, err = protojson.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))

	t.Fatal("stop")
}
