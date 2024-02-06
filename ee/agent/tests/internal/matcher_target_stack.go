package internal

import (
	"encoding/json"
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	gomegaTypes "github.com/onsi/gomega/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type targetStackMatcher struct {
	stack *v1beta1.Stack
}

func (s targetStackMatcher) Match(actual interface{}) (success bool, err error) {
	object, ok := actual.(client.Object)
	if !ok {
		return false, fmt.Errorf("expect object of type client.Object")
	}

	data, err := json.Marshal(object)
	if err != nil {
		return false, err
	}

	asMap := make(map[string]any)
	if err := json.Unmarshal(data, &asMap); err != nil {
		return false, err
	}

	if _, ok := asMap["spec"]; !ok {
		return false, nil
	}

	if _, ok := asMap["spec"].(map[string]any); !ok {
		return false, nil
	}

	if _, ok := asMap["spec"].(map[string]any)["stack"]; !ok {
		return false, nil
	}

	return true, nil
}

func (s targetStackMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("object %s should target stack %s",
		actual.(client.Object).GetName(), s.stack.GetName())
}

func (s targetStackMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("object %s should not target stack %s",
		actual.(client.Object).GetName(), s.stack.GetName())
}

var _ gomegaTypes.GomegaMatcher = (*targetStackMatcher)(nil)

func TargetStack(stack *v1beta1.Stack) gomegaTypes.GomegaMatcher {
	return &targetStackMatcher{
		stack: stack,
	}
}
