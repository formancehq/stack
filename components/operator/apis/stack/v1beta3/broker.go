package v1beta3

type KafkaSASLConfig struct {
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Mechanism    string `json:"mechanism"`
	ScramSHASize string `json:"scramSHASize"`
}

type KafkaConfig struct {
	Brokers []string `json:"brokers"`
	// +optional
	TLS bool `json:"tls"`
	// +optional
	SASL *KafkaSASLConfig `json:"sasl,omitempty"`
}

type NatsConfig struct {
	Hostname string `json:"hostname"`

	// +kubebuilder:default:=4222
	// +optional
	Port int32 `json:"port"`

	// +kubebuilder:default:=3
	// +optional
	Replicas int `json:"replicas"`

	// +kubebuilder:default:=8222
	// +optional
	MonitoringPort int32 `json:"monitoringPort,omitempty"`
}
