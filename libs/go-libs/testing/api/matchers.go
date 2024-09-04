package api

import (
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/onsi/gomega/types"
	"github.com/pkg/errors"
)

type HaveErrorCodeMatcher struct {
	lastSeen string
	expected string
}

func (s *HaveErrorCodeMatcher) Match(actual interface{}) (success bool, err error) {
	err, ok := actual.(error)
	if !ok {
		return false, fmt.Errorf("expected input type error, was %T", actual)
	}

	errorResponse := api.ErrorResponse{}
	if !errors.As(err, &errorResponse) {
		return false, nil
	}
	s.lastSeen = errorResponse.ErrorCode

	return errorResponse.ErrorCode == s.expected, nil
}

func (s *HaveErrorCodeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("error should have code %s but have %s", s.expected, s.lastSeen)
}

func (s *HaveErrorCodeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("error should not have code %s", s.expected)
}

var _ types.GomegaMatcher = (*HaveErrorCodeMatcher)(nil)

func HaveErrorCode(expected string) *HaveErrorCodeMatcher {
	return &HaveErrorCodeMatcher{
		expected: expected,
	}
}
