package testserver

import (
	"context"
	"fmt"
	"github.com/formancehq/stack/ee/ingester/internal"
	"github.com/onsi/gomega/types"
)

type HaveStateMatcher struct {
	state         ingester.StateLabel
	lastSeenState string
	srv           *Server
}

func (s *HaveStateMatcher) Match(actual interface{}) (success bool, err error) {
	id, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("expected input type %T, was %T", id, actual)
	}

	pipeline, err := GetPipeline(context.Background(), s.srv, id)
	if err != nil {
		return false, err
	}
	s.lastSeenState = pipeline.State.Label

	return pipeline.State.Label == string(s.state), nil
}

func (s *HaveStateMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("pipeline %s should have state %s but has %s", actual.(string), s.state, s.lastSeenState)
}

func (s *HaveStateMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("pipeline %s should not have state %s", actual.(string), s.state)
}

var _ types.GomegaMatcher = (*HaveStateMatcher)(nil)

func HaveState(srv *Server, state ingester.StateLabel) *HaveStateMatcher {
	return &HaveStateMatcher{
		state: state,
		srv:   srv,
	}
}
