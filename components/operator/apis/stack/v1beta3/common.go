package v1beta3

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

func (p CommonServiceProperties) GetVersion() string {
	return p.Version
}
