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

// DatabaseSpec defines the desired state of Database
type DatabaseSpec struct {
	StackDependency `json:",inline"`
	Service         string `json:"service"`
	Debug           bool   `json:"debug,omitempty"`
}

// DatabaseStatus defines the observed state of Database
type DatabaseStatus struct {
	CommonStatus `json:",inline"`
	//+optional
	URI *URI `json:"uri,omitempty"`
	//+optional
	Database string `json:"database,omitempty"`
	//+optional
	OutOfSync bool `json:"outOfSync,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Stack",type=string,JSONPath=".spec.stack",description="Stack"
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.ready",description="Ready"
//+kubebuilder:printcolumn:name="Out of sync",type=string,JSONPath=".status.outOfSync",description="Is the databse configuration out of sync"
//+kubebuilder:printcolumn:name="Info",type=string,JSONPath=".status.info",description="Info"

// Database is the Schema for the databases API
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSpec   `json:"spec,omitempty"`
	Status DatabaseStatus `json:"status,omitempty"`
}

func (in *Database) SetReady(b bool) {
	in.Status.SetReady(b)
}

func (in *Database) IsReady() bool {
	return in.Status.Ready
}

func (in *Database) SetError(s string) {
	in.Status.SetError(s)
}

func (a *Database) GetStack() string {
	return a.Spec.Stack
}

func (a *Database) isResource() {}

//+kubebuilder:object:root=true

// DatabaseList contains a list of Database
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}
