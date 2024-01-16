package internal

import (
	"fmt"
	"reflect"

	gomegaTypes "github.com/onsi/gomega/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type targetDeployment struct {
	deployment *appsv1.Deployment
}

func (t targetDeployment) Match(actual interface{}) (success bool, err error) {
	service, ok := actual.(*corev1.Service)
	if !ok {
		return false, mismatchTypeError(&corev1.Service{}, actual)
	}
	selector := service.Spec.Selector
	if !reflect.DeepEqual(t.deployment.Spec.Selector.MatchLabels, selector) {
		return false, nil
	}
	if !reflect.DeepEqual(t.deployment.Spec.Template.Labels, selector) {
		return false, nil
	}
	return true, nil
}

func (t targetDeployment) FailureMessage(actual interface{}) (message string) {
	service := actual.(*corev1.Service)
	selector := service.Spec.Selector
	if !reflect.DeepEqual(t.deployment.Spec.Selector.MatchLabels, selector) {
		return fmt.Sprintf("expected .spec.selector to match %s, got %s", t.deployment.Spec.Selector.MatchLabels, selector)
	}
	if !reflect.DeepEqual(t.deployment.Spec.Template.Labels, selector) {
		return fmt.Sprintf("expected .spec.selector to match %s, got %s", t.deployment.Spec.Template.Labels, selector)
	}
	return ""
}

func (t targetDeployment) NegatedFailureMessage(actual interface{}) (message string) {
	service := actual.(*corev1.Service)
	selector := service.Spec.Selector
	if !reflect.DeepEqual(t.deployment.Spec.Selector.MatchLabels, selector) {
		return fmt.Sprintf("expected .spec.selector to not match %s, got %s", t.deployment.Spec.Selector.MatchLabels, selector)
	}
	if !reflect.DeepEqual(t.deployment.Spec.Template.Labels, selector) {
		return fmt.Sprintf("expected .spec.selector to not match %s, got %s", t.deployment.Spec.Template.Labels, selector)
	}
	return ""
}

func TargetDeployment(spec *appsv1.Deployment) gomegaTypes.GomegaMatcher {
	return targetDeployment{
		deployment: spec,
	}
}
