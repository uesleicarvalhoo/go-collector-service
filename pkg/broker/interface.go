package broker

import (
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/schemas"
	"github.com/uesleicarvalhoo/go-collector-service/internal/infra/config"
)

type (
	CreateTopicInput = schemas.CreateTopicInput
	Config           = config.BrokerConfig
	Event            = models.Event
)
