package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConditionTypeDeploymentReady = "DeploymentReady"
)

func SetDeploymentReady(object Object, msg ...string) {
	SetCondition(object, ConditionTypeDeploymentReady, metav1.ConditionTrue, msg...)
}

func SetDeploymentError(object Object, msg ...string) {
	SetCondition(object, ConditionTypeDeploymentReady, metav1.ConditionFalse, msg...)
}
