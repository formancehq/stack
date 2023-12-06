package v1beta1

// +kubebuilder:object:generate=true
type SearchSpec struct {
	ImageHolder `json:",inline"`

	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`

	//+optional
	Ingress *IngressConfig `json:"ingress"`
}
