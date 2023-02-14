package v1beta2

type PostgresConfigCreateDatabase struct {
	PostgresConfigWithDatabase `json:",inline"`
	CreateDatabase             bool `json:"createDatabase"`
}

type CollectorConfig struct {
	KafkaConfig `json:",inline"`
	Topic       string `json:"topic"`
}
