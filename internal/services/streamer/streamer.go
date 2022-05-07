package streamer

import (
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/schemas"
)

type Streamer struct {
	eventTopic string
	broker     Broker
}

func NewStreamer(broker Broker, eventTopicInput schemas.CreateTopicInput) (*Streamer, error) {
	streamer := &Streamer{broker: broker, eventTopic: eventTopicInput.Name}

	err := broker.DeclareTopic(eventTopicInput)
	if err != nil {
		return nil, err
	}

	return streamer, nil
}

func (s *Streamer) NotifyPublishedFile(fileKey string, file models.File) error {
	event := models.Event{
		Topic: s.eventTopic,
		Key:   "published",
		Data:  map[string]string{"file_key": fileKey},
	}

	return s.broker.SendEvent(event)
}

func (s *Streamer) NotifyInvalidFile(file models.File) error {
	event := models.Event{
		Topic: s.eventTopic,
		Key:   "invalid",
		Data:  map[string]string{"file_path": file.FilePath},
	}

	return s.broker.SendEvent(event)
}
