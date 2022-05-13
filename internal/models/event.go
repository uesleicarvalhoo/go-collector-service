package models

import (
	"strings"
)

type Event struct {
	Topic string
	Key   string
	Data  any
}

func NewEvent(topic, key string, data any) (Event, error) {
	event := Event{
		Topic: topic,
		Key:   key,
		Data:  data,
	}

	if err := event.validate(); err != nil {
		return Event{}, err
	}

	return event, nil
}

func (e *Event) validate() error {
	validator := newValidator()
	if strings.TrimSpace(e.Topic) == "" {
		validator.AddError(ValidationErrorProps{Context: "event", Message: "topic should be informed"})
	}

	if strings.TrimSpace(e.Key) == "" {
		validator.AddError(ValidationErrorProps{Context: "event", Message: "key should be informed"})
	}

	if validator.HasErrors() {
		return validator.GetError()
	}

	return nil
}
