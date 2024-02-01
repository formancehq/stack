package v1beta1

type CollectorConfig struct {
	KafkaConfig `json:",inline"`
	Topic       string `json:"topic"`
}
