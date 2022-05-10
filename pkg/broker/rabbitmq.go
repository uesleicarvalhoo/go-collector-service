package broker

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/streadway/amqp"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

type RabbitMQClient struct {
	cfg        Config
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMqClient(cfg Config) (*RabbitMQClient, error) {
	con, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s", cfg.User, cfg.Password, net.JoinHostPort(cfg.Host, cfg.Port)))
	if err != nil {
		return nil, err
	}

	channel, err := con.Channel()
	if err != nil {
		return nil, err
	}

	client := &RabbitMQClient{
		cfg:        cfg,
		connection: con,
		channel:    channel,
	}

	return client, nil
}

func (mq *RabbitMQClient) Close() {
	mq.channel.Close()
	mq.connection.Close()
}

func (mq *RabbitMQClient) SendEvent(event Event) error {
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
