package main

import (
	"context"
	"fmt"

	"github.com/uesleicarvalhoo/go-collector-service/internal/infra/config"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/sender"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/streamer"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/broker"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/fileserver"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/storage"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

func main() {
	logger.Initialize()

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
		logger.Fatal(err)
	}
	defer provider.Close(ctx)

	// Broker
	brokerService, err := broker.NewRabbitMqClient(env.BrokerConfig)
	defer brokerService.Close()

	if err != nil {
		panic(err)
	}

	// Streamer
	topic := broker.CreateTopicInput{Name: env.BrokerConfig.EventTopic}

	streamerService, err := streamer.NewStreamer(brokerService, topic)
	if err != nil {
		panic(err)
	}

	// Storage
	storage := storage.NewS3Storage(env.StorageConfig, env.AwsRegion)

	// FileSerrver
	fileServer, err := fileserver.NewSFTP(env.FileServerConfig)
	if err != nil {
		panic(err)
	}

	// Run service
	senderService := sender.NewSender(streamerService, storage, fileServer)

	senderService.Start(env.MatchPattern)
}
