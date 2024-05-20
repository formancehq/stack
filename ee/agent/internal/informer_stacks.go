package internal

import (
	"reflect"

	"github.com/formancehq/stack/components/agent/internal/generated"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

type StackEventHandler struct {
	logger sharedlogging.Logger
	client MembershipClient
}

func (h *StackEventHandler) sendStatus(stackName string, status *structpb.Struct) {
	if err := h.client.Send(&generated.Message{
		Message: &generated.Message_StatusChanged{
			StatusChanged: &generated.StatusChanged{
				ClusterName: stackName,
				Statuses:    status,
			},
		},
	}); err != nil {
		h.logger.Errorf("Unable to send stack status to server: %s", err)
	}
}

func (h *StackEventHandler) AddFunc(obj interface{}) {
	stack := obj.(*unstructured.Unstructured)
	status, err := getStatus(stack)
	if err != nil {
		h.logger.Errorf("Unable to generate message stack add: %s", err)
		return
	}

	if status == nil {
		return
	}

	h.sendStatus(stack.GetName(), status)
	h.logger.Infof("Stack '%s' added", stack.GetName())

}

func (h *StackEventHandler) UpdateFunc(oldObj, newObj interface{}) {
	oldStack := oldObj.(*unstructured.Unstructured)
	newStack := newObj.(*unstructured.Unstructured)

	oldStatus, err := getStatus(oldStack)
	if err != nil {
		h.logger.Errorf("Unable to get old stack status update: %s", err)
	}

	newStatus, err := getStatus(newStack)
	if err != nil {
		h.logger.Errorf("Unable to get new stack status update: %s", err)
		return
	}

	oldDisabled, _, err := unstructured.NestedBool(oldStack.Object, "spec", "disabled")
	if err != nil {
		h.logger.Errorf("Unable to get new stack `spec.disabled` update: %s", err)
		return
	}
	newDisabled, _, err := unstructured.NestedBool(newStack.Object, "spec", "disabled")
	if err != nil {
		h.logger.Errorf("Unable to get new stack `spec.disabled` update: %s", err)
		return
	}

	// There is no status
	// The status has not changed and stack is not been disabled or enabled
	if newStatus == nil || (reflect.DeepEqual(oldStatus, newStatus) && oldDisabled == newDisabled) {
		return
	}

	h.sendStatus(newStack.GetName(), newStatus)
	h.logger.Infof("Stack '%s' updated", newStack.GetName())

}

func (h *StackEventHandler) DeleteFunc(obj interface{}) {
	stack := obj.(*unstructured.Unstructured)
	if err := h.client.Send(&generated.Message{
		Message: &generated.Message_StackDeleted{
			StackDeleted: &generated.DeletedStack{
				ClusterName: stack.GetName(),
			},
		},
	}); err != nil {
		h.logger.Errorf("Unable to send stack delete to server: %s", err)
		return
	}
	h.logger.Infof("Stack '%s' deleted", stack.GetName())
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
