package interpreter

import (
	"slices"
	"testing"
)

type ReconcileTestCase struct {
	Senders     []Sender
	Receivers   []Receiver
	Expected    []Posting
	ExpectedErr error
}

func runReconcileTestCase(t *testing.T, tc ReconcileTestCase) {
	got, err := Reconcile(tc.Senders, tc.Receivers)

	if !slices.Equal(got, tc.Expected) {
		t.Fatalf(`Expected = %#v; got = %#v`, tc.Expected, got)
	}

	if err != tc.ExpectedErr {
		t.Fatalf(`Expected err = %#v; got = %#v`, tc.ExpectedErr, err)
	}
}

func TestReconcileEmpty(t *testing.T) {
	runReconcileTestCase(t, ReconcileTestCase{})
}

func TestReconcileSingletonExactMatch(t *testing.T) {
	runReconcileTestCase(t, ReconcileTestCase{
		Senders:   []Sender{{"src", 10}},
		Receivers: []Receiver{{"dest", 10}},
		Expected:  []Posting{{"src", "dest", 10}},
	})
}

func TestNoReceiversLeft(t *testing.T) {
	runReconcileTestCase(t, ReconcileTestCase{
		Senders: []Sender{{"src", 10}},
	})
}

func TestNoSendersLeft(t *testing.T) {
	runReconcileTestCase(t, ReconcileTestCase{
		Receivers:   []Receiver{{"dest", 10}},
		ExpectedErr: ReconcileError{},
	})
}

func TestReconcileSendersRemainder(t *testing.T) {
	runReconcileTestCase(t, ReconcileTestCase{
		Senders:   []Sender{{"src", 100}},
		Receivers: []Receiver{{"d1", 70}, {"d2", 30}},
		Expected: []Posting{
			{"src", "d1", 70},
			{"src", "d2", 30},
		},
	})
}

func TestReconcileWhenSendersAreSplit(t *testing.T) {
	runReconcileTestCase(t, ReconcileTestCase{
		Senders:   []Sender{{"s1", 20}, {"s2", 30}},
		Receivers: []Receiver{{"d", 50}},
		Expected: []Posting{
			{"s1", "d", 20},
			{"s2", "d", 30},
		},
	})
}
