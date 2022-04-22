package streamer

import "github.com/uesleicarvalhoo/go-collector-service/internal/models"

type Broker interface {
	CreateQueue(name string) error
	SendEvent(models.Event) error
}
