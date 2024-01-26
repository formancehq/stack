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
	"fmt"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/url"
)

// DatabaseSpec defines the desired state of Database
type DatabaseSpec struct {
	StackDependency `json:",inline"`
	Service         string `json:"service"`
}

// +k8s:openapi-gen=true
// +kubebuilder:validation:Type=string
type URI struct {
	*url.URL `json:"-"`
}

func (u URI) String() string {
	if u.URL == nil {
		return "nil"
	}
	return u.URL.String()
}

func (u URI) IsZero() bool {
	return u.URL == nil
}

func (u *URI) DeepCopyInto(v *URI) {
	cp := *u.URL
	if u.User != nil {
		cp.User = pointer.For(*u.User)
	}
	v.URL = pointer.For(cp)
}

func (u *URI) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, u.String())), nil
}

func (u *URI) UnmarshalJSON(data []byte) error {
	v, err := url.Parse(string(data[1 : len(data)-1]))
	if err != nil {
		panic(err)
	}

	*u = URI{
		URL: v,
	}
	return nil
}

func ParseURL(v string) (*URI, error) {
	ret, err := url.Parse(v)
	if err != nil {
		return nil, err
	}
	return &URI{
		URL: ret,
	}, nil
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

func (a Database) GetStack() string {
	return a.Spec.Stack
}

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
