package testing

import (
	"fmt"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	gomegaTypes "github.com/onsi/gomega/types"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type beOwnedByMatcher struct {
	owner client.Object
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
			APIVersion:         groupVersionsKinds[0].GroupVersion().String(),
			Kind:               groupVersionsKinds[0].Kind,
			Name:               s.owner.GetName(),
			UID:                s.owner.GetUID(),
			Controller:         pointer.For(true),
			BlockOwnerDeletion: pointer.For(true),
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

func BeOwnedBy(owner client.Object) gomegaTypes.GomegaMatcher {
	return &beOwnedByMatcher{
		owner: owner,
	}
}
