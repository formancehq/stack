package v1beta1

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/formancehq/stack/libs/go-libs/pointer"
	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:object:generate=false
type EventPublisher interface {
	isEventPublisher()
}

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
	// reason contains a programmatic identifier indicating the reason for the condition's last transition.
	// Producers of specific condition types may define expected values and meanings for this field,
	// and whether the values are considered a guaranteed API.
	// The value should be a CamelCase string.
	// This field may not be empty.
	// +optional
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:Pattern=`^([A-Za-z]([A-Za-z0-9_,:]*[A-Za-z0-9_])?)?$`
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
}

type CommonStatus struct {
	//+optional
	Ready bool `json:"ready"`
	//+optional
	Info string `json:"info,omitempty"`
}

func (c *CommonStatus) SetReady(ready bool) {
	c.Ready = ready
}

func (c *CommonStatus) SetError(err string) {
	c.Info = err
}

type StatusWithConditions struct {
	CommonStatus `json:",inline"`
	//+optional
	Conditions []Condition `json:"conditions,omitempty"`
}

func (c *StatusWithConditions) DeleteCondition(t, reason string) {
	for i, existingCondition := range c.Conditions {
		if existingCondition.Type == t && existingCondition.Reason == reason {
			if i < len(c.Conditions)-1 {
				c.Conditions = append(c.Conditions[:i], c.Conditions[i+1:]...)
			} else {
				c.Conditions = c.Conditions[:i]
			}
			return
		}
	}
}

func (c *StatusWithConditions) SetCondition(condition Condition) {
	c.DeleteCondition(condition.Type, condition.Reason)
	c.Conditions = append(c.Conditions, condition)
}

type ModuleStatus struct {
	StatusWithConditions `json:",inline"`
}

type AuthConfig struct {
	// +optional
	ReadKeySetMaxRetries int `json:"readKeySetMaxRetries"`
	// +optional
	CheckScopes bool `json:"checkScopes"`
}

// +kubebuilder:object:generate=false
type Module interface {
	Dependent
	GetVersion() string
	GetConditions() []Condition
	IsDebug() bool
	IsDev() bool
	IsEE() bool
}

type ModuleProperties struct {
	DevProperties `json:",inline"`
	//+optional
	Version string `json:"version,omitempty"`
}

func (in *ModuleProperties) CompareVersion(stack *Stack, version string) int {
	actualVersion := in.Version
	if actualVersion == "" {
		actualVersion = stack.Spec.Version
	}
	if !semver.IsValid(actualVersion) {
		return 1
	}

	return semver.Compare(actualVersion, version)
}

// +kubebuilder:object:generate=false
type Dependent interface {
	Object
	GetStack() string
}

type StackDependency struct {
	Stack string `json:"stack,omitempty" yaml:"-"`
}

func (d StackDependency) GetStack() string {
	return d.Stack
}

// +kubebuilder:object:generate=false
type Object interface {
	client.Object
	SetReady(bool)
	IsReady() bool
	SetError(string)
}

// +kubebuilder:object:generate=false
type Resource interface {
	Dependent
	isResource()
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
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}

	v, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	*u = URI{
		URL: v,
	}
	return nil
}

func (in *URI) WithoutQuery() *URI {
	cp := *in.URL
	cp.ForceQuery = false
	cp.RawQuery = ""
	return &URI{
		URL: &cp,
	}
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

func init() {
	if err := equality.Semantic.AddFunc(func(a, b *URI) bool {
		if a == nil && b != nil {
			return false
		}
		if a != nil && b == nil {
			return false
		}
		if a == nil && b == nil {
			return true
		}
		return a.String() == b.String()
	}); err != nil {
		panic(err)
	}
}

const (
	StackLabel = "formance.com/stack"
	SkipLabel  = "formance.com/skip"
)
