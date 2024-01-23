package v1beta3

type DevProperties struct {
	// +optional
	Debug bool `json:"debug"`
	// +optional
	Dev bool `json:"dev"`
}

type Resource struct {
	Cpu    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type ResourceProperties struct {
	// +optional
	Request *Resource `json:"request,omitempty"`
	// +optional
	Limits *Resource `json:"limits,omitempty"`
}

type CommonServiceProperties struct {
	DevProperties `json:",inline"`
	// +optional
	Disabled *bool `json:"disabled,omitempty"`
}
