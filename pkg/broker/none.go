package broker

import (
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

type NoneBroker struct{}

func NewNoneBroker() *NoneBroker {
	logger.Warning("Using NoneBroker, all events are ignored")

	return &NoneBroker{}
}

func (n *NoneBroker) Close() {
	logger.Info("Stop to consume")
}

func (n *NoneBroker) SendEvent(event Event) error {
	logger.Warningf("Ignored event: %+v", event)

	return nil
}

func (n *NoneBroker) DeclareTopic(payload CreateTopicInput) error {
	return nil
}
