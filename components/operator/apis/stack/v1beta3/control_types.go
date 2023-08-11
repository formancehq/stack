package v1beta3

// +kubebuilder:object:generate=true
type ControlSpec struct {
	// +optional
	Annotations AnnotationsServicesSpec `json:"annotations,omitempty"`
}
