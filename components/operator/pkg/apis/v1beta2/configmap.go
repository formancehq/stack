package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConditionTypeConfigMapReady = "ConfigMapReady"
)

func SetConfigMapReady(object Object, msg ...string) {
	SetCondition(object, ConditionTypeConfigMapReady, metav1.ConditionTrue, msg...)
}

func SetConfigMapError(object Object, msg ...string) {
	SetCondition(object, ConditionTypeConfigMapReady, metav1.ConditionFalse, msg...)
}
