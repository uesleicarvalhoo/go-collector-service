package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/uesleicarvalhoo/go-collector-service/internal/infra/config"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/collector"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/sender"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/streamer"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/broker"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/storage"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

func main() {
	ctx := context.Background()
	env := config.LoadAppSettingsFromEnv()

	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Tracer
	provider, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: fmt.Sprintf("%s/api/traces", env.TraceURL),
		ServiceName:    env.TraceServiceName,
		ServiceVersion: config.ServiceVersion,
		Environment:    env.Env,
		Disabled:       false,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer provider.Close(ctx)

	// Broker
	brokerService, err := broker.NewRabbitMqClient(env.BrokerConfig)
	defer brokerService.Close()

	if err != nil {
		panic(err)
	}

	// Streamer
	streamerService, err := streamer.NewStreamer(brokerService, broker.CreateTopicInput{Name: env.BrokerConfig.EventTopic})
	if err != nil {
		panic(err)
	}

	// Storage
	storage := storage.NewS3Storage(env.StorageConfig, env.AwsRegion)

	// Collector
	fileCollector, err := collector.NewLocalCollector("/home/uescarvalho/go-collector-service/data/*.json")
	if err != nil {
		panic(err)
	}

	// Run service
	senderService := sender.NewSender(streamerService, storage)

	senderService.Consume(fileCollector)
}
