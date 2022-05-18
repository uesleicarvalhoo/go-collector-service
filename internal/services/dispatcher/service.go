package dispatcher

import (
	"github.com/uesleicarvalhoo/go-collector-service/internal/services"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/sender"
)

// Create and manage sender services, one service is created binding each config.
type Dispatcher struct {
	workerPool []*sender.Sender
}

func New(
	config Config, storage services.Storage, fileServer services.FileServer, broker services.Broker,
) (*Dispatcher, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	workerPool := []*sender.Sender{}

	for senderID, cfg := range config.SenderConfig {
		worker, err := sender.New(senderID+1, cfg, storage, fileServer, broker)
		if err != nil {
			return nil, err
		}

		workerPool = append(workerPool, worker)
	}

	return &Dispatcher{workerPool: workerPool}, nil
}

func (d *Dispatcher) Start() {
	go func() {
		for _, worker := range d.workerPool {
			worker.Start()
		}
	}()
}

func (d *Dispatcher) Stop() {
	go func() {
		for _, worker := range d.workerPool {
			worker.Stop()
		}
	}()
}
