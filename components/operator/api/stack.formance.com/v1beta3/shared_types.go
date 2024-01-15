package v1beta3

type DatabaseSpec struct {
	// +optional
	Url string `json:"url,omitempty"`
	// +optional
	Type string `json:"type,omitempty"`
}

type IngressConfig struct {
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}
