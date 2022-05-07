package broker

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

type SQSClient struct {
	Channel chan Event
	session *session.Session
}

func NewSQSClient(cfg Config, region string) (*SQSClient, error) {
	uri := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	client := &SQSClient{
		Channel: make(chan Event, 50),
		session: session.Must(session.NewSession(&aws.Config{
			Region:   aws.String(region),
			Endpoint: aws.String(uri),
		})),
	}

	return client, nil
}

func (svc *SQSClient) Close() {
}

func (svc *SQSClient) SendEvent(event Event) error {
	eventBody, err := svc.getEventBody(event.Data)
	if err != nil {
		return err
	}

	sqsSvc := sqs.New(svc.session)

	queueURL, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: &event.Topic})
	if err != nil {
		return err
	}

	_, err = sqsSvc.SendMessage(
		&sqs.SendMessageInput{
			MessageBody: eventBody,
			QueueUrl:    queueURL.QueueUrl,
		},
	)

	return err
}

func (svc *SQSClient) DeclareTopic(payload CreateTopicInput) error {
	queueAtributes := map[string]*string{}

	for k, v := range payload.Attributes {
		queueAtributes[k] = aws.String(v)
	}

	_, err := sqs.New(svc.session).CreateQueue(
		&sqs.CreateQueueInput{
			QueueName:  aws.String(payload.Name),
			Attributes: queueAtributes,
		},
	)

	return err
}

func (svc *SQSClient) getEventBody(data any) (*string, error) {
	eventData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Couldn't decode event data: %s\n", err)

		return nil, err
	}

	return aws.String(string(eventData)), nil
}
