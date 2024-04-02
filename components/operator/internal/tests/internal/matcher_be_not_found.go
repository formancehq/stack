package internal

import (
	"fmt"

	gomegaTypes "github.com/onsi/gomega/types"
	controllererrors "k8s.io/apimachinery/pkg/api/errors"
)

type beNotFound struct{}

func (b beNotFound) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, nil
	}
	err, ok := actual.(error)
	if !ok {
		return false, fmt.Errorf("expected error type, got %T", actual)
	}
	return controllererrors.IsNotFound(err), nil
}

func (b beNotFound) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("should be not found, got %v", actual)
}

func (b beNotFound) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("should be found, got %v", actual)
}

var _ gomegaTypes.GomegaMatcher = (*beNotFound)(nil)

func BeNotFound() gomegaTypes.GomegaMatcher {
	return &beNotFound{}
}
