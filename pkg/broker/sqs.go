package broker

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSClient struct {
	Channel chan Event
}

func NewSQSClient(cfg Config) (*SQSClient, error) {
	sqs.New(
		session.New(),
		&aws.Config{
			Endpoint:   aws.String(""),
			DisableSSL: aws.Bool(true),
		},
	)

	client := &SQSClient{
		Channel: make(chan Event, 50),
	}

	return client, nil
}

func (s *SQSClient) Close() {
}

func (s *SQSClient) SendEvent(event Event) error {
	s.Channel <- event

	return nil
}

func (s *SQSClient) CreateQueue(name string) error {
	return nil
}
