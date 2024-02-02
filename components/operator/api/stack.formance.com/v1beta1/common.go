package v1beta1

type DevProperties struct {
	// +optional
	Debug bool `json:"debug"`
	// +optional
	Dev bool `json:"dev"`
}

type CommonServiceProperties struct {
	DevProperties `json:",inline"`
	// +optional
	//+kubebuilder:default:="latest"
	Version string `json:"version,omitempty"`
}
