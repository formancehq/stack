package internal

import (
	"slices"
	"strings"
	"sync"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	"github.com/formancehq/stack/components/agent/internal/generated"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

type InMemoryStacksModules struct {
	mu     sync.Mutex
	stacks map[string][]string
	onPush map[string]func(modules []string)
}

func (m *InMemoryStacksModules) Get(stackName string) ([]string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ret, ok := m.stacks[stackName]
	return ret, ok
}

func (m *InMemoryStacksModules) Push(stackName string, modules []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stacks[stackName] = modules
	fn, ok := m.onPush[stackName]
	if ok {
		fn(modules)
	}
	delete(m.onPush, stackName)
}

func (m *InMemoryStacksModules) GetExpectedModules(stackName string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.stacks[stackName]; !ok {
		return []string{}
	}
	return m.stacks[stackName]
}

func (m *InMemoryStacksModules) OnPush(stackName string, fn func(modules []string)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	modules, ok := m.stacks[stackName] // Added while locking
	if ok {
		fn(modules)
		return
	}
	m.onPush[stackName] = fn
}

func (m *InMemoryStacksModules) Contains(stack, module string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	modules, ok := m.stacks[stack] // Added while locking
	if !ok {
		return false
	}

	return collectionutils.Contains(modules, module)
}

func NewInMemoryStacksModules() *InMemoryStacksModules {
	return &InMemoryStacksModules{
		stacks: make(map[string][]string),
		onPush: make(map[string]func([]string)),
	}
}

func NewStackEventHandler(logger sharedlogging.Logger, membershipClient MembershipClient, stacks *InMemoryStacksModules) cache.ResourceEventHandlerFuncs {
	sendStatus := func(interpretedStatus generated.StackStatus, stackName string, status *structpb.Struct) {
		logger.Infof("Send status %s for stack %s", interpretedStatus, stackName)
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

	InterpretedStatus := func(stack *unstructured.Unstructured, expectedModules []string) generated.StackStatus {
		logger := logger.WithField("stack", stack.GetName())

		disabled, found, err := unstructured.NestedBool(stack.Object, "spec", "disabled")
		if !found || err != nil {
			panic(err)
		}
		if disabled {
			logger.Infof("Set status as disabled as the stack is disabled")
			return generated.StackStatus_Disabled
		}

		ready, found, err := unstructured.NestedBool(stack.Object, "status", "ready")
		if !found || err != nil || !ready {
			logger.Infof("Set status as progressing as the stack is not ready")
			return generated.StackStatus_Progressing
		}

		stackModules, _, err := unstructured.NestedStringSlice(stack.Object, "status", "modules")
		if err != nil {
			panic(err)
		}
		stackModules = collectionutils.Map(stackModules, func(v string) string {
			return strings.ToLower(v)
		})
		for _, module := range expectedModules {
			if !slices.Contains(stackModules, module) {
				logger.Infof("Set status as progressing as expected module '%s' is not in stack modules, have %s", module, stackModules)
				return generated.StackStatus_Progressing
			}
		}
		return generated.StackStatus_Ready
	}

	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			stack := obj.(*unstructured.Unstructured)

			status, err := getStatus(stack)
			if err != nil {
				logger.Errorf("Unable to generate message stack update: %s", err)
				return
			}

			if status == nil {
				return
			}

			modules, ok := stacks.Get(stack.GetName())
			if !ok {
				logger.Debugf("Stack '%s' not initialized in memory", stack.GetName())
				stacks.OnPush(stack.GetName(), func(modules []string) {
					logger.Debugf("Stack '%s' finally sent by membership, update status", stack.GetName())
					sendStatus(InterpretedStatus(stack, modules), stack.GetName(), status)
				})
				return
			}

			logger.Infof("Stack '%s' added", stack.GetName())
			sendStatus(InterpretedStatus(stack, modules), stack.GetName(), status)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {

			newStack := newObj.(*unstructured.Unstructured)

			status, err := getStatus(newStack)
			if err != nil {
				logger.Errorf("Unable to generate message stack update: %s", err)
				return
			}

			if status == nil {
				return
			}

			modules, ok := stacks.Get(newStack.GetName())
			if !ok {
				logger.Debugf("Stack '%s' not initialized in memory", newStack.GetName())
				stacks.OnPush(newStack.GetName(), func(modules []string) {
					logger.Debugf("Stack '%s' finally sent by membership, update status", newStack.GetName())
					sendStatus(InterpretedStatus(newStack, modules), newStack.GetName(), status)
				})
				return
			}

			logger.Infof("Stack '%s' updated", newStack.GetName())
			sendStatus(InterpretedStatus(newStack, modules), newStack.GetName(), status)
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

			modules, ok := stacks.Get(stack.GetName())
			if !ok {
				logger.Debugf("Stack '%s' not initialized in memory", stack.GetName())
				return
			}

			logger.Infof("Stack '%s' deleted", stack.GetName())
			status, err := getStatus(stack)
			if err != nil {
				logger.Errorf("Unable to generate message stack update: %s", err)
				stacks.OnPush(stack.GetName(), func(modules []string) {
					logger.Debugf("Stack '%s' finally sent by membership, update status", stack.GetName())
					sendStatus(InterpretedStatus(stack, modules), stack.GetName(), status)
				})
				return
			}
			if status == nil {
				return
			}
			sendStatus(InterpretedStatus(stack, modules), stack.GetName(), status)
		},
	}
}
