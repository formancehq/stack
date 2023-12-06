/*
Copyright 2022.

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

package v1beta3

import (
	"reflect"

	"github.com/iancoleman/strcase"
	"golang.org/x/mod/semver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VersionsSpec defines the desired state of Versions
type VersionsSpec struct {
	// +optional
	Control string `json:"control"`
	// +optional
	Ledger string `json:"ledger"`
	// +optional
	Payments string `json:"payments"`
	// +optional
	Search string `json:"search"`
	// +optional
	Auth string `json:"auth"`
	// +optional
	Webhooks string `json:"webhooks"`
	// +optional
	Wallets string `json:"wallets"`
	// +optional
	Orchestration string `json:"orchestration"`
	// +optional
	Gateway string `json:"gateway"`
	// +optional
	Stargate string `json:"stargate"`
	// +optional
	Reconciliation string `json:"reconciliation"`
}

// VersionsStatus defines the observed state of Versions
type VersionsStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
//+kubebuilder:resource:scope=Cluster

// Versions is the Schema for the versions API
type Versions struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VersionsSpec   `json:"spec,omitempty"`
	Status VersionsStatus `json:"status,omitempty"`
}

type VersionDiff int

const (
	Lower VersionDiff = iota
	Equals
	Higher
	Undefined
)

func (v *Versions) GetVersion(module string) string {
	return reflect.ValueOf(v.Spec).FieldByName(strcase.ToCamel(module)).Interface().(string)
}

func (v *Versions) Compare(module, version string) VersionDiff {
	if !semver.IsValid(version) {
		panic("provided version is not semver")
	}

	versionInConfig := v.GetVersion(module)
	if !semver.IsValid(versionInConfig) {
		return Undefined
	}

	switch semver.Compare(version, versionInConfig) {
	case 1:
		return Lower
	case 0:
		return Equals
	case -1:
		return Higher
	}
	return Undefined
}

func (v *Versions) IsLower(module, version string) bool {
	return v.Compare(module, version) == Lower
}

func (v *Versions) IsHigher(module, version string) bool {
	return v.Compare(module, version) >= Higher
}

func (v *Versions) IsHigherOrEqual(module, version string) bool {
	cmp := v.Compare(module, version)
	return cmp >= Equals
}

func (*Versions) Hub() {}

func (in *Versions) GetFromModuleName(s string) string {
	fieldByName := reflect.ValueOf(in.Spec).FieldByName(strcase.ToCamel(s))
	if fieldByName.String() == "" {
		return "latest"
	}
	return fieldByName.String()
}

//+kubebuilder:object:root=true

// VersionsList contains a list of Versions
type VersionsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Versions `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Versions{}, &VersionsList{})
}
