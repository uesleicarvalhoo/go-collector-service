package streamer

import (
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

type Streamer struct {
	eventChannel chan models.Event
	broker       services.Broker
	quit         chan bool
}

func New(broker services.Broker, eventChannel chan models.Event) (*Streamer, error) {
	return &Streamer{
		broker:       broker,
		eventChannel: eventChannel,
	}, nil
}

func (s *Streamer) Start() {
	go func() {
		for {
			select {
			case <-s.quit:
				return
			case event := <-s.eventChannel:
				logger.Debugf("Event received %+v", event)
				_ = s.broker.SendEvent(event)
			}
		}
	}()
}

func (s *Streamer) Stop() {
	go func() {
		s.quit <- true
	}()
}
