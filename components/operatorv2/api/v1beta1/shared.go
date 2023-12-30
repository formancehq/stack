package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type DevProperties struct {
	// +optional
	Debug bool `json:"debug"`
	// +optional
	Dev bool `json:"dev"`
}

func (p DevProperties) IsDebug() bool {
	return p.Debug
}

func (p DevProperties) IsDev() bool {
	return p.Dev
}

type resource struct {
	Cpu    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type ResourceProperties struct {
	// +optional
	Request *resource `json:"request,omitempty"`
	// +optional
	Limits *resource `json:"limits,omitempty"`
}

type CommonServiceProperties struct {
	DevProperties `json:",inline"`
	//+optional
	Version string `json:"version,omitempty"`
	//+optional
	ResourceProperties *ResourceProperties `json:"resourceProperties,omitempty"`
}

// Condition contains details for one aspect of the current state of this API Resource.
// ---
// This struct is intended for direct use as an array at the field path .status.conditions.  For example,
//
//	type FooStatus struct{
//	    // Represents the observations of a foo's current state.
//	    // Known .status.conditions.type are: "Available", "Progressing", and "Degraded"
//	    // +patchMergeKey=type
//	    // +patchStrategy=merge
//	    // +listType=map
//	    // +listMapKey=type
//	    CommonStatus []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
//
//	    // other fields
//	}
type Condition struct {
	// type of condition in CamelCase or in foo.example.com/CamelCase.
	// ---
	// Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be
	// useful (see .node.status.conditions), the ability to deconflict is important.
	// The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt)
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$`
	// +kubebuilder:validation:MaxLength=316
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`
	// status of the condition, one of True, False, Unknown.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status metav1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status"`
	// observedGeneration represents the .metadata.generation that the condition was set based upon.
	// For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
	// with respect to the current state of the instance.
	// +optional
	// +kubebuilder:validation:Minimum=0
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`
	// lastTransitionTime is the last time the condition transitioned from one status to another.
	// This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// message is a human readable message indicating details about the transition.
	// This may be an empty string.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=32768
	Message string `json:"message" protobuf:"bytes,6,opt,name=message"`
}

type CommonStatus struct {
	//+optional
	Conditions []Condition `json:"conditions,omitempty"`
	//+optional
	Ready bool `json:"ready"`
	//+optional
	Error string `json:"error,omitempty"`
}

func (c *CommonStatus) SetStatus(status bool, error string) {
	c.Ready = true
}

func (c *CommonStatus) DeleteCondition(t string) {
	for i, existingCondition := range c.Conditions {
		if existingCondition.Type == t {
			if i < len(c.Conditions)-1 {
				c.Conditions = append(c.Conditions[:i], c.Conditions[i+1:]...)
			} else {
				c.Conditions = c.Conditions[:i]
			}
			return
		}
	}
}

func (c *CommonStatus) SetCondition(condition Condition) {
	c.DeleteCondition(condition.Type)
	c.Conditions = append(c.Conditions, condition)
}

type StackDependency struct {
	Stack string `json:"stack,omitempty" yaml:"-"`
}

func (d StackDependency) GetStack() string {
	return d.Stack
}
