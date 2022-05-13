package broker

import (
	"github.com/uesleicarvalhoo/go-collector-service/internal/config"
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/schemas"
)

type (
	CreateTopicInput = schemas.CreateTopicInput
	Config           = config.BrokerConfig
	Event            = models.Event
)
