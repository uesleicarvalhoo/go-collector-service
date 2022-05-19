package config

type StorageConfig struct {
	URL    string `envconfig:"STORAGE_HOST" default:"http://localhost.localstack.cloud:4566"`
	User   string `envconfig:"STORAGE_USER"`
	Key    string `envconfig:"STORAGE_KEY"`
	Bucket string `envconfig:"STORAGE_BUCKET" default:"collector-files"`
}
