package internal

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal/generated"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"k8s.io/client-go/tools/cache"
)

func sendStatus(logger sharedlogging.Logger, membershipClient MembershipClient, stack string, status generated.StackStatus) error {
	if err := membershipClient.Send(&generated.Message{
		Message: &generated.Message_StatusChanged{
			StatusChanged: &generated.StatusChanged{
				ClusterName: stack,
				Status:      status,
			},
		},
	}); err != nil {
		logger.Errorf("Unable to send stack status to server: %s", err)
		return err
	}

	return nil
}

func evaluateStackStatus(isReady, isDisabled bool) generated.StackStatus {
	if isDisabled {
		return generated.StackStatus_Disabled
	}

	if isReady {
		return generated.StackStatus_Ready
	}

	return generated.StackStatus_Progressing

}

func sendStatusFromStack(logger sharedlogging.Logger, membershipClient MembershipClient, stack *v1beta1.Stack) error {
	return sendStatus(logger, membershipClient, stack.Name, evaluateStackStatus(stack.Status.Ready, stack.Spec.Disabled))
}

func StacksEventHandler(logger sharedlogging.Logger, membershipClient MembershipClient) cache.ResourceEventHandlerFuncs {

	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			stack := convertUnstructured[*v1beta1.Stack](obj)
			logger.Infof("Stack '%s' added", stack.Name)
			sendStatusFromStack(logger, membershipClient, stack)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {

			oldStack := convertUnstructured[*v1beta1.Stack](oldObj)
			newStack := convertUnstructured[*v1beta1.Stack](newObj)

			logger.Debugf("Stack OLD(%s) status ready ?'%t' updated, disabled ? %t", oldStack.Name, oldStack.Status.Ready, oldStack.Spec.Disabled)
			logger.Debugf("Stack(%s) status ready ? '%t' updated, disabled ? %t", newStack.Name, newStack.Status.Ready, newStack.Spec.Disabled)

			if oldStack.Spec.Disabled == newStack.Spec.Disabled && oldStack.Status.Ready == newStack.Status.Ready {
				return
			}

			logger.Infof("Stack '%s' updated", newStack.Name)
			sendStatusFromStack(logger, membershipClient, newStack)
		},
		DeleteFunc: func(obj interface{}) {
			stack := convertUnstructured[*v1beta1.Stack](obj)
			logger.Infof("Stack '%s' deleted", stack.Name)
			sendStatus(logger, membershipClient, stack.Name, generated.StackStatus_Deleted)
		},
	}
}
