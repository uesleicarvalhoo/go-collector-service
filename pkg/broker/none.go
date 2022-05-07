package broker

import (
	"github.com/sirupsen/logrus"
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
	logrus.Info("Stop to consume")
}

func (n *NoneBroker) startPublish() {
	logrus.Info("Start to Consume")
}

func (n *NoneBroker) SendEvent(event Event) error {
	n.Channel <- event

	return nil
}

func (n *NoneBroker) DeclareTopic(payload CreateTopicInput) error {
	return nil
}
