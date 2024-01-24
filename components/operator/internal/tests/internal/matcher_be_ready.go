package internal

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"

	gomegaTypes "github.com/onsi/gomega/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type beReadyMatcher struct{}

func (s beReadyMatcher) Match(actual interface{}) (success bool, err error) {
	object, ok := actual.(v1beta1.Object)
	if !ok {
		return false, fmt.Errorf("expect object of type core.Object")
	}
	return object.IsReady(), nil
}

func (s beReadyMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("object %s should be ready",
		actual.(client.Object).GetName())
}

func (s beReadyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("object %s should not be ready",
		actual.(client.Object).GetName())
}

var _ gomegaTypes.GomegaMatcher = (*beReadyMatcher)(nil)

func BeReady() gomegaTypes.GomegaMatcher {
	return &beReadyMatcher{}
}
