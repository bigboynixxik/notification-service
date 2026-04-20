package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string `env:"APP_ENV" envDefault:"dev"`

	KafkaBrokers    string `env:"KAFKA_BROKERS"`
	KafkaTopic      string `env:"KAFKA_TOPIC"`
	KafkaGroup      string `env:"KAFKA_GROUP_ID"`
	AuthServiceAddr string `env:"AUTH_SERVICE_ADDR"`

	Token string `env:"TOKEN"`
}

func LoadConfig(path string) (*Config, error) {
	_ = godotenv.Load(path)
	//if err != nil {
	//	return nil, fmt.Errorf("config.LoadConfig: %w", err)
	//}
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig failed to parse config: %w", err)
	}
	return &cfg, nil
}
