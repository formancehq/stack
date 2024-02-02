package internal

import (
	"fmt"
	"reflect"

	"github.com/formancehq/stack/libs/go-libs/pointer"
	gomegaTypes "github.com/onsi/gomega/types"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type beOwnedByOption func(matcher *beOwnedByMatcher)

func WithBlockOwnerDeletion() beOwnedByOption {
	return func(matcher *beOwnedByMatcher) {
		matcher.blockOwnerDeletion = true
	}
}

type beOwnedByMatcher struct {
	owner              client.Object
	controller         bool
	blockOwnerDeletion bool
}

func (s beOwnedByMatcher) Match(actual interface{}) (success bool, err error) {
	object, ok := actual.(client.Object)
	if !ok {
		return false, fmt.Errorf("expect object of type runtime.Object")
	}
	for _, reference := range object.GetOwnerReferences() {
		groupVersionsKinds, _, err := Client().Scheme().ObjectKinds(s.owner)
		if err != nil {
			return false, errors.Wrap(err, "searching object kinds")
		}
		expectedOwnerReference := metav1.OwnerReference{
			APIVersion: groupVersionsKinds[0].GroupVersion().String(),
			Kind:       groupVersionsKinds[0].Kind,
			Name:       s.owner.GetName(),
			UID:        s.owner.GetUID(),
		}
		if s.controller {
			expectedOwnerReference.Controller = pointer.For(true)
		}
		if s.blockOwnerDeletion {
			expectedOwnerReference.BlockOwnerDeletion = pointer.For(true)
		}
		if reflect.DeepEqual(reference, expectedOwnerReference) {
			return true, nil
		}
	}
	return false, nil
}

func (s beOwnedByMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("object %s should be owned by %s",
		actual.(client.Object).GetName(), (any)(s.owner))
}

func (s beOwnedByMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("object %s should not be owned by %s",
		actual.(client.Object).GetName(), (any)(s.owner))
}

var _ gomegaTypes.GomegaMatcher = (*beOwnedByMatcher)(nil)

func BeOwnedBy(owner client.Object, opts ...beOwnedByOption) gomegaTypes.GomegaMatcher {
	ret := &beOwnedByMatcher{
		owner: owner,
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
