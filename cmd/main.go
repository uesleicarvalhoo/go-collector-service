package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/uesleicarvalhoo/go-collector-service/internal/config"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/dispatcher"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/broker"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/fileserver"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/storage"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

func main() {
	cfg := config.Settings{}
	if err := cfg.LoadFromEnv(); err != nil {
		panic(err)
	}

	if err := logger.InitLogger(cfg.LoggerConfig); err != nil {
		panic(err)
	}

	// Tracer
	provider, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: fmt.Sprintf("%s/api/traces", cfg.TraceURL),
		ServiceName:    cfg.TraceServiceName,
		ServiceVersion: cfg.ServiceVersion,
		Environment:    cfg.Env,
		Disabled:       !cfg.TraceEnable,
	})
	if err != nil {
		logger.Fatal(err)
	}
	defer provider.Close(context.Background())

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
	fileServer, err := fileserver.NewLocalFileServer(cfg.FileServerConfig)
	if err != nil {
		panic(err)
	}

	// Run service
	var dispatcherCfg dispatcher.Config
	if err := dispatcherCfg.LoadFromYaml("./config.yaml"); err != nil {
		panic(err)
	}

	dispatcher, err := dispatcher.New(dispatcherCfg, storage, fileServer, brokerService)
	if err != nil {
		panic(err)
	}

	dispatcher.Start()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	dispatcher.Stop()
}
