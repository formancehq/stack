package internal

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal/generated"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"k8s.io/client-go/tools/cache"
)

func StacksEventHandler(logger sharedlogging.Logger, membershipClient MembershipClient) cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {

			oldStack := convertUnstructured[*v1beta1.Stack](oldObj)
			newStack := convertUnstructured[*v1beta1.Stack](newObj)

			if oldStack.Status.Ready == newStack.Status.Ready {
				return
			}

			logger.Infof("Stack '%s' updated", newStack.Name)

			if err := membershipClient.Send(&generated.Message{
				Message: &generated.Message_StatusChanged{
					StatusChanged: &generated.StatusChanged{
						ClusterName: newStack.Name,
						Status: func() generated.StackStatus {
							if newStack.Status.Ready {
								return generated.StackStatus_Ready
							} else {
								return generated.StackStatus_Progressing
							}
						}(),
					},
				},
			}); err != nil {
				logger.Errorf("Unable to send stack status to server: %s", err)
			}
		},
		DeleteFunc: func(obj interface{}) {
			stack := convertUnstructured[*v1beta1.Stack](obj)

			logger.Infof("Stack '%s' deleted", stack.Name)

			if err := membershipClient.Send(&generated.Message{
				Message: &generated.Message_StatusChanged{
					StatusChanged: &generated.StatusChanged{
						ClusterName: stack.Name,
						Status:      generated.StackStatus_Deleted,
					},
				},
			}); err != nil {
				logger.Errorf("Unable to send stack status to server: %s", err)
			}
		},
	}
}
