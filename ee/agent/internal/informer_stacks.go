package internal

import (
	"reflect"

	sharedlogging "github.com/formancehq/go-libs/logging"
	"github.com/formancehq/stack/components/agent/internal/generated"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

type StackEventHandler struct {
	logger sharedlogging.Logger
	client MembershipClient
}

func (h *StackEventHandler) sendStatus(stackName string, status *structpb.Struct) error {
	if err := h.client.Send(&generated.Message{
		Message: &generated.Message_StatusChanged{
			StatusChanged: &generated.StatusChanged{
				ClusterName: stackName,
				Statuses:    status,
			},
		},
	}); err != nil {
		h.logger.Errorf("Unable to send stack status to server: %s", err)
		return err
	}
	return nil
}

func (h *StackEventHandler) AddFunc(obj interface{}) {
	stack := obj.(*unstructured.Unstructured)
	logger := h.logger.WithField("func", "Add").WithField("stack", stack.GetName())

	status, err := getStatus(stack)
	if err != nil {
		logger.Errorf("Unable to generate message stack add: %s", err)
		return
	}

	if status == nil {
		return
	}

	if err := h.sendStatus(stack.GetName(), status); err != nil {
		return
	}
	logger.Infof("Stack '%s' added", stack.GetName())

}

func (h *StackEventHandler) UpdateFunc(oldObj, newObj interface{}) {
	oldStack := oldObj.(*unstructured.Unstructured)
	newStack := newObj.(*unstructured.Unstructured)

	logger := h.logger.WithField("func", "Update").WithField("stack", newStack.GetName())

	oldStatus, err := getStatus(oldStack)
	if err != nil {
		logger.Errorf("Unable to get old stack status update: %s", err)
	}

	newStatus, err := getStatus(newStack)
	if err != nil {
		logger.Errorf("Unable to get new stack status update: %s", err)
		return
	}

	oldDisabled, _, err := unstructured.NestedBool(oldStack.Object, "spec", "disabled")
	if err != nil {
		logger.Errorf("Unable to get new stack `spec.disabled` update: %s", err)
		return
	}
	newDisabled, _, err := unstructured.NestedBool(newStack.Object, "spec", "disabled")
	if err != nil {
		logger.Errorf("Unable to get new stack `spec.disabled` update: %s", err)
		return
	}

	// There is no status
	// The status has not changed and stack is not been disabled or enabled
	if newStatus == nil || (reflect.DeepEqual(oldStatus, newStatus) && oldDisabled == newDisabled) {
		return
	}

	if err := h.sendStatus(newStack.GetName(), newStatus); err != nil {
		return
	}
	logger.Infof("Stack '%s' updated", newStack.GetName())

}

func (h *StackEventHandler) DeleteFunc(obj interface{}) {
	stack := obj.(*unstructured.Unstructured)
	logger := h.logger.WithField("func", "Delete").WithField("stack", stack.GetName())

	if err := h.client.Send(&generated.Message{
		Message: &generated.Message_StackDeleted{
			StackDeleted: &generated.DeletedStack{
				ClusterName: stack.GetName(),
			},
		},
	}); err != nil {
		logger.Errorf("Unable to send stack delete to server: %s", err)
		return
	}
	logger.Infof("Stack '%s' deleted", stack.GetName())
}

func NewStackEventHandler(logger sharedlogging.Logger, membershipClient MembershipClient) cache.ResourceEventHandlerFuncs {
	stackEventHandler := &StackEventHandler{
		logger: logger,
		client: membershipClient,
	}

	return cache.ResourceEventHandlerFuncs{
		AddFunc:    stackEventHandler.AddFunc,
		UpdateFunc: stackEventHandler.UpdateFunc,
		DeleteFunc: stackEventHandler.DeleteFunc,
	}
}
