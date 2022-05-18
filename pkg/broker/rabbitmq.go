package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
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
}

func NewRabbitMqClient(cfg Config) (*RabbitMQClient, error) {
	client := &RabbitMQClient{
		cfg:        cfg,
		errChannel: make(chan *amqp.Error, 1),
	}

	if err := client.connect(); err != nil {
		return nil, err
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

	logger.Debugf("Event received, %+v", event)

	body, err := json.Marshal(event.Data)
	if err != nil {
		logger.Errorf("Couldn't decode event data: %s", err)

		return err
	}

	err = mq.channel.Publish(event.Topic, event.Key, false, false, amqp.Publishing{
		Body: body,
	})
	if err != nil {
		if errors.Is(err, amqp.ErrClosed) {
			logger.Warningf("Connection error, retrying to send event %+v", event)

			if retryErr := mq.SendEvent(event); retryErr != nil {
				return retryErr
			}

			return nil
		}

		logger.Errorf("Failed to publish event, %s", err)

		return err
	}

	return nil
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
		logger.Error("RabbitMQ connection is closed, trying stablish a new connection..")

		for i := 0; i < maxConnectionRetries; i++ {
			mq.connection = nil

			err := mq.connect()
			if err == nil {
				logger.Error("RabbitMQ connection re-established with success")

				break
			}

			logger.Errorf("Failed to re-connect, '%s', trying again in %d seconds..", err, retryConnectionDelay)
			time.Sleep(time.Second * retryConnectionDelay)
		}

		if mq.channel == nil {
			logger.Panicf("Couldn't reconnect to RabbitMQ")
		}
	}
}
