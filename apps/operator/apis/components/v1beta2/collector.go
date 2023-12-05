package v1beta2

type KafkaSASLConfig struct {
	// +optional
	Username string `json:"username,omitempty"`
	// +optional
	UsernameFrom *ConfigSource `json:"usernameFrom,omitempty"`
	// +optional
	Password string `json:"password,omitempty"`
	// +optional
	PasswordFrom *ConfigSource `json:"passwordFrom,omitempty"`
	Mechanism    string        `json:"mechanism"`
	ScramSHASize string        `json:"scramSHASize"`
}

type KafkaConfig struct {
	// +optional
	Brokers []string `json:"brokers"`
	// +optional
	BrokersFrom *ConfigSource `json:"brokersFrom"`
	// +optional
	TLS bool `json:"tls"`
	// +optional
	SASL *KafkaSASLConfig `json:"sasl,omitempty"`
}

type NatsConfig struct {
	URL string `json:"url"`
}
