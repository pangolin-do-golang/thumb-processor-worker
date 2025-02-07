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

type Smtp struct {
	Host     string `env:"SMTP_HOST"`
	Port     string `env:"SMTP_PORT"`
	Password string `env:"SMTP_PASSWORD"`
	From     string `env:"SMTP_FROM"`
}

type Config struct {
	S3   S3
	SQS  SQS
	Smtp Smtp
}

func Load() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
