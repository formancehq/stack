package internal

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal/generated"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"k8s.io/client-go/tools/cache"
)

func StacksEventHandler(logger sharedlogging.Logger, membershipClient MembershipClient) cache.ResourceEventHandlerFuncs {
	sendStatus := func(stack string, status generated.StackStatus) {
		if err := membershipClient.Send(&generated.Message{
			Message: &generated.Message_StatusChanged{
				StatusChanged: &generated.StatusChanged{
					ClusterName: stack,
					Status:      status,
				},
			},
		}); err != nil {
			logger.Errorf("Unable to send stack status to server: %s", err)
		}
	}

	sendStatusFromStack := func(stack *v1beta1.Stack) {
		sendStatus(stack.Name, func() generated.StackStatus {
			if stack.Spec.Disabled {
				return generated.StackStatus_Disabled
			}

			if stack.Status.Ready {
				return generated.StackStatus_Ready
			} else {
				return generated.StackStatus_Progressing
			}
		}())
	}

	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			stack := convertUnstructured[*v1beta1.Stack](obj)
			logger.Infof("Stack '%s' added", stack.Name)
			sendStatusFromStack(stack)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {

			oldStack := convertUnstructured[*v1beta1.Stack](oldObj)
			newStack := convertUnstructured[*v1beta1.Stack](newObj)

			if oldStack.Spec.Disabled == newStack.Spec.Disabled && oldStack.Status.Ready == newStack.Status.Ready {
				return
			}

			logger.Infof("Stack '%s' updated", newStack.Name)
			sendStatusFromStack(newStack)
		},
		DeleteFunc: func(obj interface{}) {
			stack := convertUnstructured[*v1beta1.Stack](obj)
			logger.Infof("Stack '%s' deleted", stack.Name)
			sendStatus(stack.Name, generated.StackStatus_Deleted)
		},
	}
}
