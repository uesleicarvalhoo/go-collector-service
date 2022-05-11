package config

type StorageConfig struct {
	URL    string `envconfig:"STORAGE_HOST" default:"http://localhost.localstack.cloud:4566"`
	Bucket string `envconfig:"STORAGE_BUCKET" default:"collector-files"`
}
