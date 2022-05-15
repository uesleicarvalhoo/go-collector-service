package sender

import (
	"sync"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/collector"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/publisher"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/streamer"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

type Sender struct {
	ID            int
	config        Config
	storage       services.Storage
	collector     *collector.Collector
	streamer      *streamer.Streamer
	publisherPool []*publisher.Publisher
	fileChannel   chan models.File
	eventChannel  chan models.Event
	waitGroup     *sync.WaitGroup
}

func New(
	processID int, config Config, storage services.Storage, fileServer services.FileServer, broker services.Broker,
) (*Sender, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	waitGroup := &sync.WaitGroup{}
	fileChan := make(chan models.File, config.Workers)
	eventChannel := make(chan models.Event, config.Workers)

	collector, err := collector.New(processID, config.CollectorCfg, fileServer, fileChan, waitGroup)
	if err != nil {
		return nil, err
	}

	eventStreamer, err := streamer.New(broker, eventChannel)
	if err != nil {
		return nil, err
	}

	return &Sender{
		ID:            processID,
		config:        config,
		storage:       storage,
		collector:     collector,
		streamer:      eventStreamer,
		publisherPool: []*publisher.Publisher{},
		fileChannel:   fileChan,
		eventChannel:  eventChannel,
		waitGroup:     waitGroup,
	}, nil
}

func (s *Sender) Start() {
	logger.Infof("[Sender %d] Starting with %d workers", s.ID, s.config.Workers)

	go func() {
		s.collector.Start()
		s.streamer.Start()

		for workerID := len(s.publisherPool); workerID < s.config.Workers; workerID++ {
			s.newPublisher(workerID + 1)
		}
	}()
}

func (s *Sender) Stop() {
	go func() {
		s.collector.Stop()
		s.streamer.Stop()

		for _, publisher := range s.publisherPool {
			publisher.Stop()
		}
	}()
}

func (s *Sender) newPublisher(workerID int) {
	publisher := publisher.New(workerID, s.config.EventTopic, s.storage, s.fileChannel, s.eventChannel, s.waitGroup)
	s.publisherPool = append(s.publisherPool, publisher)
	publisher.Start()
}
