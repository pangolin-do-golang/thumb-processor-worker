package config

import (
	"github.com/caarlos0/env/v11"
)

type S3 struct {
	Bucket string `env:"S3_BUCKET,required"`
}

type SQS struct {
	QueueURL string `env:"SQS_QUEUE_URL,required"`
}

type ThumbAPI struct {
	URL string `env:"THUMB_API_URL,required"`
}

type Config struct {
	S3       S3
	SQS      SQS
	ThumbAPI ThumbAPI
}

func Load() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
