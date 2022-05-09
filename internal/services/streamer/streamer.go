package streamer

import (
	"context"

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

// Send event to Message Broker to event topic with body {"file_key": <fileKey>} and routing-key "published".
func (s *Streamer) NotifyPublishedFile(ctx context.Context, fileKey string, file models.File) error {
	event, err := models.NewEvent(s.eventTopic, "published", map[string]string{"file_key": fileKey})
	if err != nil {
		return err
	}

	return s.broker.SendEvent(event)
}

// Send event to Message Broker to event topic with body {"file_path": <file.FilePath>} and routing-key "invalid".
func (s *Streamer) NotifyInvalidFile(ctx context.Context, file models.File) error {
	event, err := models.NewEvent(s.eventTopic, "invalid", map[string]string{"file_path": file.FilePath})
	if err != nil {
		return err
	}

	return s.broker.SendEvent(event)
}
