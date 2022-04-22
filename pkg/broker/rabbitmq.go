package broker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	cfg        Config
	connection *amqp.Connection
	Channel    chan Event
}

func NewRabbitMqClient(cfg Config) (*RabbitMQClient, error) {
	con, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}

	client := &RabbitMQClient{
		cfg:        cfg,
		connection: con,
		Channel:    make(chan Event, 50),
	}

	go client.startPublish()

	return client, nil
}

func (mq *RabbitMQClient) Close() {
	mq.connection.Close()
}

func (mq *RabbitMQClient) startPublish() {
	ch, err := mq.connection.Channel()
	defer ch.Close()

	if err != nil {
		log.Fatal("Couldn't start to consume events,", err)
	}

	for !mq.connection.IsClosed() {
		event := <-mq.Channel

		log.Printf("Event received, %+v\n", event)

		body, err := json.Marshal(event.Data)
		if err != nil {
			log.Printf("Couldn't decode event data: %s\n", err)

			continue
		}

		err = ch.Publish(event.Topic, event.Key, false, false, amqp.Publishing{
			Body: body,
		})
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Stopping to send messages to broker, connection was closed.")
}

func (mq *RabbitMQClient) SendEvent(event Event) error {
	mq.Channel <- event

	return nil
}

func (mq *RabbitMQClient) CreateQueue(name string) error {
	ch, err := mq.connection.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	return ch.ExchangeDeclare(name, "direct", true, false, false, false, amqp.Table{})
}
