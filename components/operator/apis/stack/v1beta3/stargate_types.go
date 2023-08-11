package v1beta3

type StargateSpec struct {
	// +optional
	Annotations AnnotationsServicesSpec `json:"service"`
}
