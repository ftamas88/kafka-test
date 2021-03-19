package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

// See https://github.com/caarlos0/env for env parse options
type Config struct {
	HTTPPort                   int      `env:"HTTP_PORT" envDefault:"3000"`
	Servers                    []string `env:"KAFKA_SERVERS" envSeparator:"," envDefault:"localhost:9092"`
	CloudServer                string   `env:"KAFKA_CLOUD_SERVERS" envDefault:""`
	SchemaRegistryServers      []string `env:"KAFKA_SCHEMA_SERVERS" envSeparator:"," envDefault:"http://localhost:8081"`
	CloudSchemaRegistryServers []string `env:"KAFKA_CLOUD_SCHEMA_SERVERS" envSeparator:"," envDefault:""`
	SchemaPath                 string   `env:"KAFKA_SCHEMA_PATH" envDefault:"pkg/avro/transaction_with_amount.avsc"`
	Topic                      string   `env:"KAFKA_TOPIC" envDefault:"test"`
	KafkaCloudKey              string   `env:"KAFKA_CLOUD_KEY" envDefault:""`
	KafkaCloudSecret           string   `env:"KAFKA_CLOUD_SECRET" envDefault:""`
}

func New() (*Config, error) {
	conf := &Config{}
	err := env.Parse(conf)
	if err != nil {
		return nil, fmt.Errorf("config env parse: %w", err)
	}

	return conf, nil
}
