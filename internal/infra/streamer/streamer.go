package streamer

import "github.com/uesleicarvalhoo/go-collector-service/internal/models"

type Streamer struct {
	broker Broker
}

func NewStreamer(broker Broker) (*Streamer, error) {
	streamer := &Streamer{broker: broker}

	err := broker.CreateQueue("collector.files")
	if err != nil {
		return nil, err
	}

	return streamer, nil
}

func (s *Streamer) NotifyPublishedFile(file models.File) error {
	event := models.Event{
		Topic: "collector.files",
		Key:   "published",
		Data:  map[string]string{"file_key": file.Key},
	}

	return s.broker.SendEvent(event)
}

func (s *Streamer) NotifyInvalidFile(file models.File) error {
	event := models.Event{
		Topic: "collector.files",
		Key:   "invalid",
		Data:  map[string]string{"file_path": file.FilePath},
	}

	return s.broker.SendEvent(event)
}
