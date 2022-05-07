package streamer

import (
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/schemas"
)

type Broker interface {
	DeclareTopic(schemas.CreateTopicInput) error
	SendEvent(models.Event) error
}
