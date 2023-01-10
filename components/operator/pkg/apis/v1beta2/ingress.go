// +kubebuilder:object:generate=true
package v1beta2

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressSpec struct {
	// +optional
	Annotations map[string]string `json:"annotations"`
	Path        string            `json:"path"`
	Host        string            `json:"host"`
	// +optional
	TLS *IngressTLS `json:"tls"`
}

const (
	ConditionTypeIngressReady = "IngressReady"
)

func SetIngressReady(object Object, msg ...string) {
	SetCondition(object, ConditionTypeIngressReady, metav1.ConditionTrue, msg...)
}

func SetIngressError(object Object, msg ...string) {
	SetCondition(object, ConditionTypeIngressReady, metav1.ConditionFalse, msg...)
}

func RemoveIngressCondition(object Object) {
	object.GetConditions().Remove(ConditionTypeIngressReady)
}

type IngressTLS struct {
	// SecretName is the name of the secret used to terminate TLS traffic on
	// port 443. Field is left optional to allow TLS routing based on SNI
	// hostname alone. If the SNI host in a listener conflicts with the "Host"
	// header field used by an IngressRule, the SNI host is used for termination
	// and value of the Host header is used for routing.
	// +optional
	SecretName string `json:"secretName,omitempty" protobuf:"bytes,2,opt,name=secretName"`
}

func (t *IngressTLS) AsK8SIngressTLSSlice() []networkingv1.IngressTLS {
	if t == nil {
		return nil
	}
	return []networkingv1.IngressTLS{{
		//SecretName: t.SecretName,
	}}
}
