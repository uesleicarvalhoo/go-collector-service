package sender

import (
	"context"
	"sync"
	"time"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/collector"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/publisher"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/streamer"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

type Sender struct {
	ID               int
	config           Config
	storage          services.Storage
	collector        *collector.Collector
	streamer         *streamer.Streamer
	publisherPool    []*publisher.Publisher
	eventChannel     chan models.Event
	collectWaitGroup *sync.WaitGroup
	processWaitGroup *sync.WaitGroup
	quit             chan bool
}

func New(
	processID int, config Config, storage services.Storage, fileServer services.FileServer, broker services.Broker,
) (*Sender, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	collectWaitGroup := &sync.WaitGroup{}
	processWaitGroup := &sync.WaitGroup{}
	eventChannel := make(chan models.Event, config.Workers)

	collector, err := collector.New(processID, config.CollectorCfg, fileServer, collectWaitGroup, processWaitGroup)
	if err != nil {
		return nil, err
	}

	eventStreamer, err := streamer.New(broker, eventChannel)
	if err != nil {
		return nil, err
	}

	eventStreamer.Start()

	return &Sender{
		ID:               processID,
		config:           config,
		storage:          storage,
		collector:        collector,
		streamer:         eventStreamer,
		publisherPool:    []*publisher.Publisher{},
		eventChannel:     eventChannel,
		collectWaitGroup: collectWaitGroup,
		processWaitGroup: processWaitGroup,
	}, nil
}

func (s *Sender) loop() {
	for {
		select {
		case <-s.quit:
			s.Stop()

			return

		default:
			startTime := time.Now()

			fileChannel := make(chan models.File, s.config.Workers)
			ctx, span := trace.NewSpan(context.Background(), "sender.loop")

			for _, worker := range s.publisherPool {
				worker.Handle(ctx, fileChannel)
			}

			s.collector.CollectFiles(ctx, fileChannel)

			s.collectWaitGroup.Wait()
			s.processWaitGroup.Wait()

			close(fileChannel)
			span.End()

			took := time.Since(startTime)
			logger.Infof("[Sender %d] Took %s", s.ID, took.String())

			time.Sleep((time.Duration(s.config.CollectDelay) * time.Second) - took)
		}
	}
}

func (s *Sender) Start() {
	logger.Infof("[Sender %d] Starting with %d workers", s.ID, s.config.Workers)

	go func() {
		s.streamer.Start()

		for workerID := len(s.publisherPool); workerID < s.config.Workers; workerID++ {
			s.newPublisher(workerID + 1)
		}

		s.loop()
	}()
}

func (s *Sender) Stop() {
	go func() {
		s.streamer.Stop()
	}()
}

func (s *Sender) newPublisher(workerID int) {
	publisher := publisher.New(workerID, s.config.EventTopic, s.storage, s.eventChannel, s.processWaitGroup)
	s.publisherPool = append(s.publisherPool, publisher)
}
