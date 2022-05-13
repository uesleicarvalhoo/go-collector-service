package fileserver

import (
	"github.com/pkg/errors"
	"github.com/uesleicarvalhoo/go-collector-service/internal/config"
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

type Config = config.FileServerConfig

var ErrConnectionFailed = errors.New("Couldn't connect")

type Locker = models.Locker
