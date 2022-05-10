package fileserver

import (
	"github.com/pkg/errors"
	"github.com/uesleicarvalhoo/go-collector-service/internal/infra/config"
)

type Config = config.FileServerConfig

var ErrConnectionFailed = errors.New("Couldn't connect")
