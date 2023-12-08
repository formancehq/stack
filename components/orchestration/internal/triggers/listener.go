package triggers

import (
	"encoding/json"
	"fmt"

	"go.temporal.io/api/enums/v1"

	"github.com/formancehq/stack/libs/go-libs/pointer"
	"go.temporal.io/api/serviceerror"

	"github.com/formancehq/stack/libs/go-libs/logging"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/client"
)

// Quick hack to filter already processed events
func getWorkflowIDFromEvent(event publish.EventMessage) *string {
	switch event.Type {
	case "SAVED_PAYMENT", "SAVED_ACCOUNT":
		data, err := json.Marshal(event.Payload)
		if err != nil {
			panic(err)
		}

		type object struct {
			ID string `json:"id"`
		}
		o := &object{}
		if err := json.Unmarshal(data, o); err != nil {
			panic(err)
		}

		return pointer.For(o.ID)
	default:
		return nil
	}
}

func handleMessage(logger logging.Logger, temporalClient client.Client, taskQueue string, msg *message.Message) error {
	logger = logger.WithFields(map[string]any{
		"event-id":  msg.UUID,
		"duplicate": "false",
	})

	var err error
	defer func() {
		if err != nil {
			logger = logger.WithField("err", err)
			logger.Errorf("Handle message")
		} else {
			logger.Infof("Handle message")
		}
	}()

	var event *publish.EventMessage
	event, err = publish.UnmarshalMessage(msg)
	if err != nil {
		return err
	}

	logger = logger.WithField("type", event.Type)

	options := client.StartWorkflowOptions{
		TaskQueue: taskQueue,
	}
	if ik := getWorkflowIDFromEvent(*event); ik != nil {
		options.ID = *ik
		options.WorkflowIDReusePolicy = enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE
		options.WorkflowExecutionErrorWhenAlreadyStarted = true
		logger = logger.WithField("ik", *ik)
	}

	var w client.WorkflowRun
	w, err = temporalClient.ExecuteWorkflow(msg.Context(), options, RunTrigger, ProcessEventRequest{
		MessageID: msg.UUID,
		Event:     *event,
	})
	if err != nil {
		_, ok := err.(*serviceerror.WorkflowExecutionAlreadyStarted)
		if ok {
			logger = logger.WithField("duplicate", "true")
			err = nil
			return nil
		}
	}
	logger = logger.WithFields(map[string]any{
		"id":     w.GetID(),
		"run-id": w.GetRunID(),
	})

	return errors.Wrap(err, "executing workflow")
}

func registerListener(logger logging.Logger, r *message.Router, s message.Subscriber, temporalClient client.Client, taskQueue string, topics []string) {
	for _, topic := range topics {
		r.AddNoPublisherHandler(fmt.Sprintf("listen-%s-events", topic), topic, s, func(msg *message.Message) error {
			if err := handleMessage(logger, temporalClient, taskQueue, msg); err != nil {
				logging.Errorf("Error executing workflow: %s", err)
				return err
			}
			return nil
		})
	}
}
