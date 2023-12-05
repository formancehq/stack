package v1beta3

type DevProperties struct {
	// +optional
	Debug bool `json:"debug"`
	// +optional
	Dev bool `json:"dev"`
}

type resource struct {
	Cpu    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type ResourceProperties struct {
	// +optional
	Request *resource `json:"request,omitempty"`
	// +optional
	Limits *resource `json:"limits,omitempty"`
}

type CommonServiceProperties struct {
	DevProperties `json:",inline"`
	// +optional
	Disabled *bool `json:"disabled,omitempty"`
}
