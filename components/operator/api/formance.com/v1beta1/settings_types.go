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

type SettingsSpec struct {
	//+optional
	// Stacks on which the setting is applied. Can contain `*` to indicate a wildcard.
	Stacks []string `json:"stacks,omitempty"`
	// The setting Key. See the documentation of each module or [global settings](#global-settings) to discover them.
	Key    string   `json:"key"`
	// The value. It must have a specific format following the Key.
	Value  string   `json:"value"`
}

// Settings represents a configurable piece of the stacks.
//
// The purpose of this resource is to be able to configure some common settings between a set of stacks.
//
// Example :
// ```yaml
// apiVersion: formance.com/v1beta1
// kind: Settings
// metadata:
//   name: postgres-uri
// spec:
//   key: postgres.ledger.uri
//   stacks:
//   - stack0
//   value: postgresql://postgresql.formance.svc.cluster.local:5432
// ```
//
// This example create a setting named `postgres-uri` targeting the stack named `stack0` and the service `ledger` (see the key `postgres.ledger.uri`).
//
// Therefore, a [Database](#database) created for the stack `stack0` and the service named 'ledger' will use the uri `postgresql://postgresql.formance.svc.cluster.local:5432`.
//
// Settings allow to use wildcards in keys and in stacks list.
//
// For example, if you want to use the same database server for all the modules of a specific stack, you can write :
// ```yaml
// apiVersion: formance.com/v1beta1
// kind: Settings
// metadata:
//   name: postgres-uri
// spec:
//   key: postgres.*.uri # There, we use a wildcard to indicate we want to use that setting of all services of the stack `stack0`
//   stacks:
//   - stack0
//   value: postgresql://postgresql.formance.svc.cluster.local:5432
// ```
//
// Also, we could use that setting for all of our stacks using :
// ```yaml
// apiVersion: formance.com/v1beta1
// kind: Settings
// metadata:
//   name: postgres-uri
// spec:
//   key: postgres.*.uri # There, we use a wildcard to indicate we want to use that setting for all services of all stacks
//   stacks:
//   - * # There we select all the stacks
//   value: postgresql://postgresql.formance.svc.cluster.local:5432
// ```
//
// Some settings are really global, while some are used by specific module.
//
// Refer to the documentation of each module and resource to discover available Settings.
//
// ##### Global settings
// ###### AWS account
//
// A stack can use an AWS account for authentication.
//
// It can be used to connect to any AWS service we could use.
//
// It includes RDS, OpenSearch and MSK. To do so, you can create the following setting:
// ```yaml
// apiVersion: formance.com/v1beta1
// kind: Settings
// metadata:
//   name: aws-service-account
// spec:
//   key: aws.service-account
//   stacks:
//   - '*'
//   value: aws-access
// ```
// This setting instruct the operator than there is somewhere on the cluster a service account named `aws-access`.
//
// So, each time a service has the capability to use AWS, the operator will use this service account.
//
// The service account could look like that :
// ```yaml
// apiVersion: v1
// kind: ServiceAccount
// metadata:
//   annotations:
//     eks.amazonaws.com/role-arn: arn:aws:iam::************:role/staging-eu-west-1-hosting-stack-access
//   labels:
//     formance.com/stack: any
//   name: aws-access
// ```
// You can note two things :
// 1. We have an annotation indicating the role arn used to connect to AWS. Refer to the AWS documentation to create this role
// 2. We have a label `formance.com/stack=any` indicating we are targeting all stacks.
//    Refer to the documentation of [ResourceReference](#resourcereference) for further information.
//
// ###### JSON logging
//
// You can use the setting `logging.json` with the value `true` to configure elligible service to log as json.
// Example:
// ```yaml
// apiVersion: formance.com/v1beta1
// kind: Settings
// metadata:
//   name: json-logging
// spec:
//   key: logging.json
//   stacks:
//   - '*'
//   value: "true"
// ```
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Key",type=string,JSONPath=".spec.key",description="Key"
//+kubebuilder:printcolumn:name="Value",type=string,JSONPath=".spec.value",description="Value"
type Settings struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SettingsSpec   `json:"spec,omitempty"`
}

func (in *Settings) GetStacks() []string {
	return in.Spec.Stacks
}

func (in *Settings) IsWildcard() bool {
	return len(in.Spec.Stacks) == 1 && in.Spec.Stacks[0] == "*"
}

//+kubebuilder:object:root=true

// SettingsList contains a list of Settings
type SettingsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Settings `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Settings{}, &SettingsList{})
}
