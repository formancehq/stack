package v1beta3

// +kubebuilder:object:generate=true
type WalletsSpec struct {
	CommonServiceProperties `json:",inline"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}
