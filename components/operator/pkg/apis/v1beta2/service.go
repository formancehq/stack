package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConditionTypeServiceReady = "ServiceReady"
)

func SetServiceReady(object Object, msg ...string) {
	SetCondition(object, ConditionTypeServiceReady, metav1.ConditionTrue, msg...)
}

func SetServiceError(object Object, msg ...string) {
	SetCondition(object, ConditionTypeServiceReady, metav1.ConditionFalse, msg...)
}
