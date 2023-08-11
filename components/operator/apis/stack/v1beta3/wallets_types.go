package v1beta3

// +kubebuilder:object:generate=true
type WalletsSpec struct {
	// +optional
	Annotations AnnotationsServicesSpec `json:"annotations,omitempty"`
}
