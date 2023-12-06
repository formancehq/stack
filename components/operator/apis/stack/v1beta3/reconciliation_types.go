package v1beta3

type ReconciliationSpec struct {
	CommonServiceProperties `json:",inline"`

	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}
