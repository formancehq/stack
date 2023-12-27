package v1beta1

type DevProperties struct {
	// +optional
	Debug bool `json:"debug"`
	// +optional
	Dev bool `json:"dev"`
}

func (p DevProperties) IsDebug() bool {
	return p.Debug
}

func (p DevProperties) IsDev() bool {
	return p.Dev
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
	//+optional
	Version string `json:"version,omitempty"`
	//+optional
	ResourceProperties *ResourceProperties `json:"resourceProperties,omitempty"`
}
