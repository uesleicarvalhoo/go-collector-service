package broker

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	cfg        Config
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMqClient(cfg Config) (*RabbitMQClient, error) {
	con, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port))
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
	logrus.Infof("Event received, %+v\n", event)

	body, err := json.Marshal(event.Data)
	if err != nil {
		logrus.Infof("Couldn't decode event data: %s\n", err)

		return err
	}

	err = mq.channel.Publish(event.Topic, event.Key, false, false, amqp.Publishing{
		Body: body,
	})
	if err != nil {
		logrus.Infof("Failed to publish event, %s\n", err)

		return err
	}

	return nil
}

func (mq *RabbitMQClient) DeclareTopic(payload CreateTopicInput) error {
	ch, err := mq.connection.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	exchangeType, ok := payload.Attributes["type"]
	if !ok {
		exchangeType = "topic"
	}

	return ch.ExchangeDeclare(payload.Name, exchangeType, true, false, false, false, amqp.Table{})
}
