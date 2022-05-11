package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uesleicarvalhoo/go-collector-service/internal/infra/config"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/sender"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/broker"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/fileserver"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/storage"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

func main() {
	logger.Initialize()

	ctx := context.Background()
	cfg := config.LoadAppSettingsFromEnv()

	// Tracer
	provider, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: fmt.Sprintf("%s/api/traces", cfg.TraceURL),
		ServiceName:    cfg.TraceServiceName,
		ServiceVersion: cfg.ServiceVersion,
		Environment:    cfg.Env,
		Disabled:       false,
	})
	if err != nil {
		logger.Fatal(err)
	}
	defer provider.Close(ctx)

	// Broker
	brokerService, err := broker.NewRabbitMqClient(
		cfg.BrokerConfig, broker.CreateTopicInput{Name: cfg.BrokerConfig.EventTopic},
	)
	if err != nil {
		panic(err)
	}
	defer brokerService.Close()

	// Storage
	storage := storage.NewS3Storage(cfg.StorageConfig, cfg.AwsRegion)

	// FileSerrver
	fileServer, err := fileserver.NewSFTP(cfg.FileServerConfig)
	if err != nil {
		panic(err)
	}

	// Run service
	senderCfg := sender.Config{
		MatchPatterns: cfg.MatchPatterns,
		Workers:       cfg.ParalelUploads,
		EventTopic:    cfg.EventTopic,
		Delay:         time.Second * time.Duration(cfg.CollectDelay),
	}

	senderService, err := sender.New(senderCfg, storage, brokerService, fileServer)
	if err != nil {
		panic(err)
	}

	go senderService.Start()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	senderService.Shutdown()
}
