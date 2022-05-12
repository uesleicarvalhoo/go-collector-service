package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

const (
	maxConnectionRetries = 5
	retryConnectionDelay = 1
)

type RabbitMQClient struct {
	cfg        Config
	connection *amqp.Connection
	channel    *amqp.Channel
	errChannel chan *amqp.Error
	sync.Mutex
}

func NewRabbitMqClient(cfg Config, topics ...CreateTopicInput) (*RabbitMQClient, error) {
	client := &RabbitMQClient{
		cfg:        cfg,
		errChannel: make(chan *amqp.Error, 1),
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	for _, topic := range topics {
		if err := client.DeclareTopic(topic); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (mq *RabbitMQClient) Close() {
	mq.channel.Close()
	mq.connection.Close()
}

func (mq *RabbitMQClient) SendEvent(event Event) error {
	if err := mq.connect(); err != nil {
		return err
	}

	logger.Infof("Event received, %+v", event)

	body, err := json.Marshal(event.Data)
	if err != nil {
		logger.Infof("Couldn't decode event data: %s", err)

		return err
	}

	err = mq.channel.Publish(event.Topic, event.Key, false, false, amqp.Publishing{
		Body: body,
	})
	if err != nil {
		if errors.Is(err, amqp.ErrClosed) {
			logger.Infof("Connection error, retrying to send event %+v", event)

			if retryErr := mq.SendEvent(event); retryErr != nil {
				return retryErr
			}

			return nil
		}

		logger.Infof("Failed to publish event, %s", err)

		return err
	}

	return nil
}

func (mq *RabbitMQClient) DeclareTopic(payload CreateTopicInput) error {
	channel, err := mq.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	exchangeType, ok := payload.Attributes["type"]
	if !ok {
		exchangeType = "topic"
	}

	return channel.ExchangeDeclare(payload.Name, exchangeType, true, false, false, false, amqp.Table{})
}

func (mq *RabbitMQClient) connect() error {
	if mq.connection != nil && !mq.connection.IsClosed() {
		return nil
	}

	uri := fmt.Sprintf("amqp://%s:%s@%s", mq.cfg.User, mq.cfg.Password, net.JoinHostPort(mq.cfg.Host, mq.cfg.Port))

	con, err := amqp.Dial(uri)
	if err != nil {
		return err
	}

	channel, err := con.Channel()
	if err != nil {
		return err
	}

	mq.errChannel = channel.NotifyClose(make(chan *amqp.Error))

	con.IsClosed()
	mq.connection = con
	mq.channel = channel

	go mq.handleConnectionError()

	return nil
}

func (mq *RabbitMQClient) handleConnectionError() {
	for range mq.errChannel {
		mq.Lock()
		defer mq.Unlock()

		logger.Error("RabbitMQ connection is closed, trying stablish a new connection..")

		mq.connection = nil
		for i := 0; i < maxConnectionRetries; i++ {
			if err := mq.connect(); err == nil {
				logger.Error("RabbitMQ connection re-established with success")

				return
			}

			logger.Errorf("Failed to re-connect, trying again in %d seconds..", retryConnectionDelay)
			time.Sleep(time.Second * retryConnectionDelay)
		}

		logger.Panicf("Couldn't reconnect to RabbitMQ")
	}
}
