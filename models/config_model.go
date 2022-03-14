package models

type Config struct {
	Kafka KafkaConfig
	Log   struct {
		Level string `yaml:"level"`
	}
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}
