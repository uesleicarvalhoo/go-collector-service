package main

import (
	"context"
	"fmt"
	"log"

	"github.com/uesleicarvalhoo/go-collector-service/internal/infra/collector"
	"github.com/uesleicarvalhoo/go-collector-service/internal/infra/config"
	"github.com/uesleicarvalhoo/go-collector-service/internal/infra/streamer"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/broker"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/storage"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

func main() {
	ctx := context.Background()
	env := config.LoadAppSettingsFromEnv()

	// Tracer
	provider, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: fmt.Sprintf("%s/api/traces", env.TraceURL),
		ServiceName:    env.TraceServiceName,
		ServiceVersion: config.ServiceVersion,
		Environment:    env.Env,
		Disabled:       false,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer provider.Close(ctx)

	// Broker
	eventBroker, err := broker.NewRabbitMqClient(env.BrokerConfig)
	if err != nil {
		panic(err)
	}
	defer eventBroker.Close()

	// Streamer
	streamer, err := streamer.NewStreamer(eventBroker)
	if err != nil {
		panic(err)
	}

	// Storage
	storage := storage.NewS3Storage(env.StorageConfig, env.AwsRegion)

	// Collector
	fileCollector := collector.NewLocalCollector("/home/uescarvalho/go-collector-service/tmp/*.json")

	// Run service
	sender := services.NewSender(streamer, storage, fileCollector)

	sender.Run()
}
