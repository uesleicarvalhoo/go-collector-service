package broker

import (
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

type NoneBroker struct {
	Channel chan Event
}

func NewNoneBroker() *NoneBroker {
	return &NoneBroker{
		Channel: make(chan Event, 50),
	}
}

func (n *NoneBroker) Close() {
	logger.Info("Stop to consume")
}

func (n *NoneBroker) startPublish() {
	logger.Info("Start to Consume")
}

func (n *NoneBroker) SendEvent(event Event) error {
	n.Channel <- event

	return nil
}

func (n *NoneBroker) DeclareTopic(payload CreateTopicInput) error {
	return nil
}
