/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceReferenceSpec struct {
	StackDependency  `json:",inline"`
	GroupVersionKind *metav1.GroupVersionKind `json:"gvk"`
	Name             string                   `json:"name"`
}

type ResourceReferenceStatus struct {
	Status `json:",inline"`
	//+optional
	SyncedResource string `json:"syncedResource,omitempty"`
	//+optional
	Hash string `json:"hash,omitempty"`
}

// ResourceReference is a special resources used to refer to externally created resources.
//
// It includes k8s service accounts and secrets.
//
// Why? Because the operator create a namespace by stack, so, a stack does not have access to secrets and service
// accounts created externally.
//
// A ResourceReference is created by other resource who need to use a specific secret or service account.
// For example, if you want to use a secret for your database connection (see [Database](#database), you will
// create a setting indicating a secret name. You will need to create this secret yourself, and you will put this
// secret inside the namespace you want (`default` maybe).
//
// The Database reconciler will create a ResourceReference looking like that :
// ```
// apiVersion: formance.com/v1beta1
// kind: ResourceReference
// metadata:
//
//	name: jqkuffjxcezj-qlii-auth-postgres
//	ownerReferences:
//	- apiVersion: formance.com/v1beta1
//	  blockOwnerDeletion: true
//	  controller: true
//	  kind: Database
//	  name: jqkuffjxcezj-qlii-auth
//	  uid: 2cc4b788-3ffb-4e3d-8a30-07ed3941c8d2
//
// spec:
//
//	gvk:
//	  group: ""
//	  kind: Secret
//	  version: v1
//	name: postgres
//	stack: jqkuffjxcezj-qlii
//
// status:
//
//	...
//
// ```
// This reconciler behind this ResourceReference will search, in all namespaces, for a secret named "postgres".
// The secret must have a label `formance.com/stack` with the value matching either a specific stack or `any` to target any stack.
//
// Once the reconciler has found the secret, it will copy it inside the stack namespace, allowing the ResourceReconciler owner to use it.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
type ResourceReference struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceReferenceSpec   `json:"spec,omitempty"`
	Status ResourceReferenceStatus `json:"status,omitempty"`
}

func (in *ResourceReference) SetReady(b bool) {
	in.Status.SetReady(b)
}

func (in *ResourceReference) IsReady() bool {
	return in.Status.Ready
}

func (in *ResourceReference) SetError(s string) {
	in.Status.SetError(s)
}

func (in *ResourceReference) GetStack() string {
	return in.Spec.Stack
}

func (in *ResourceReference) GetConditions() *Conditions {
	return &in.Status.Conditions
}

//+kubebuilder:object:root=true

// ResourceReferenceList contains a list of ResourceReference
type ResourceReferenceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceReference `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceReference{}, &ResourceReferenceList{})
}
