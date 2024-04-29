package internal

import (
	"slices"

	"github.com/formancehq/stack/components/agent/internal/generated"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

type InMemoryStacksModules map[string][]string

func GetExpectedModules(stackName string, stacks InMemoryStacksModules) []string {
	if _, ok := stacks[stackName]; !ok {
		return []string{}
	}
	return stacks[stackName]
}

func NewStackEventHandler(logger sharedlogging.Logger, membershipClient MembershipClient, stacks InMemoryStacksModules) cache.ResourceEventHandlerFuncs {
	sendStatus := func(interpretedStatus generated.StackStatus, stackName string, status *structpb.Struct) {
		if err := membershipClient.Send(&generated.Message{
			Message: &generated.Message_StatusChanged{
				StatusChanged: &generated.StatusChanged{
					ClusterName: stackName,
					Status:      interpretedStatus,
					Statuses:    status,
				},
			},
		}); err != nil {
			logger.Errorf("Unable to send stack status to server: %s", err)
		}
	}

	InterpretedStatus := func(stack *unstructured.Unstructured) generated.StackStatus {
		disabled, found, err := unstructured.NestedBool(stack.Object, "spec", "disabled")
		if !found || err != nil {
			panic(err)
		}
		if disabled {
			return generated.StackStatus_Disabled
		}

		ready, found, err := unstructured.NestedBool(stack.Object, "status", "ready")
		if !found || err != nil || !ready {
			return generated.StackStatus_Progressing
		}

		modules, _, err := unstructured.NestedStringSlice(stack.Object, "status", "modules")
		if err != nil {
			panic(err)
		}
		for _, module := range GetExpectedModules(stack.GetName(), stacks) {
			if !slices.Contains(modules, module) {
				return generated.StackStatus_Progressing
			}
		}
		return generated.StackStatus_Ready
	}

	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			stack := obj.(*unstructured.Unstructured)
			if _, ok := stacks[stack.GetName()]; !ok {
				logger.Debugf("Stack '%s' not initialized in memory", stack.GetName())
				return
			}

			status, err := getStatus(stack)
			if err != nil {
				logger.Errorf("Unable to generate message stack update: %s", err)
				return
			}

			logger.Infof("Stack '%s' added", stack.GetName())
			sendStatus(InterpretedStatus(stack), stack.GetName(), status)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {

			newStack := newObj.(*unstructured.Unstructured)

			if _, ok := stacks[newStack.GetName()]; !ok {
				logger.Debugf("Stack '%s' not initialized in memory", newStack.GetName())
				return
			}

			status, err := getStatus(newStack)
			if err != nil {
				logger.Errorf("Unable to generate message stack update: %s", err)
				return
			}

			logger.Infof("Stack '%s' updated", newStack.GetName())
			sendStatus(InterpretedStatus(newStack), newStack.GetName(), status)
		},
		DeleteFunc: func(obj interface{}) {
			stack := obj.(*unstructured.Unstructured)
			if err := membershipClient.Send(&generated.Message{
				Message: &generated.Message_StackDeleted{
					StackDeleted: &generated.DeletedStack{
						ClusterName: stack.GetName(),
					},
				},
			}); err != nil {
				logger.Errorf("Unable to send stack delete to server: %s", err)
			}

			if _, ok := stacks[stack.GetName()]; !ok {
				logger.Debugf("Stack '%s' not initialized in memory", stack.GetName())
				return
			}

			logger.Infof("Stack '%s' deleted", stack.GetName())
			status, err := getStatus(stack)
			if err != nil {
				logger.Errorf("Unable to generate message stack update: %s", err)
				return
			}
			sendStatus(InterpretedStatus(stack), stack.GetName(), status)

		},
	}
}
