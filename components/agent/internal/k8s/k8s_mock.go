package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type K8SClientMock struct {
	Stacks map[string]*v1beta3.Stack

	readiness map[*v1beta3.Stack]bool

	fakeWatch *K8SClientMockWatcher
}

func NewK8SClientMock() *K8SClientMock {
	return &K8SClientMock{
		fakeWatch: NewK8SClientMockWatcher(),
		Stacks:    make(map[string]*v1beta3.Stack),
		readiness: make(map[*v1beta3.Stack]bool),
	}
}

func (c *K8SClientMock) Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta3.Stack, error) {
	stack, ok := c.Stacks[name]
	if !ok {
		return nil, fmt.Errorf("stack not found")
	}
	return stack, nil
}

func (c *K8SClientMock) Create(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error) {
	if _, ok := c.Stacks[stack.Name]; ok {
		return nil, fmt.Errorf("stack already exists")
	}
	c.Stacks[stack.Name] = stack

	c.fakeWatch.Add(stack)

	return stack, nil
}

func (c *K8SClientMock) Update(ctx context.Context, stack *v1beta3.Stack) (*v1beta3.Stack, error) {
	if _, ok := c.Stacks[stack.Name]; !ok {
		return nil, fmt.Errorf("stack not found")
	}
	c.Stacks[stack.Name] = stack
	c.fakeWatch.Modify(stack)
	return stack, nil
}
func (c *K8SClientMock) Delete(ctx context.Context, name string) error {
	if _, ok := c.Stacks[name]; !ok {
		return fmt.Errorf("stack not found")
	}
	delete(c.Stacks, name)
	c.fakeWatch.Delete(c.Stacks[name])
	return nil
}
func (c *K8SClientMock) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {

	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				for _, stack := range c.Stacks {
					if stack.Spec.Disabled && c.readiness[stack] {
						c.readiness[stack] = false
						c.fakeWatch.Modify(stack)
						continue
					}

					if !c.readiness[stack] {
						stack.Status.Conditions = append(stack.Status.Conditions, v1beta3.Condition{
							Type: v1beta3.ConditionTypeProgressing,
						})
					}

					c.fakeWatch.Modify(stack)
				}
			case <-time.After(10 * time.Second):
				for _, stack := range c.Stacks {
					if stack.Spec.Disabled {
						continue
					}
					if !c.readiness[stack] {
						stack.Status.Conditions = append(stack.Status.Conditions, v1beta3.Condition{
							Type: v1beta3.ConditionTypeReady,
						})
						c.readiness[stack] = true
						c.fakeWatch.Modify(stack)
					}
				}
			}
		}
	}()

	return nil, c.fakeWatch
}

type K8SClientMockWatcher struct {
	watch.RaceFreeFakeWatcher
}

// Error implements error.
func (*K8SClientMockWatcher) Error() string {
	panic("unimplemented")
}

func NewK8SClientMockWatcher() *K8SClientMockWatcher {
	return &K8SClientMockWatcher{
		RaceFreeFakeWatcher: *watch.NewRaceFreeFake(),
	}
}
