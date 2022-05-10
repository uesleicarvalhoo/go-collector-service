package config

type StorageConfig struct {
	URL    string `env:"STORAGE_HOST,default=http://localhost.localstack.cloud:4566"`
	Bucket string `env:"STORAGE_BUCKET,default=collector-files"`
}
