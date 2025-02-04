package config

import (
	"github.com/caarlos0/env/v11"
)

type S3 struct {
	Bucket string `env:"S3_BUCKET"`
}

type SQS struct {
	QueueURL string `env:"SQS_QUEUE_URL"`
}

type Config struct {
	S3  S3
	SQS SQS
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(&cfg)
	return cfg, err
}
