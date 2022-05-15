package collector

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

type Collector struct {
	ID          int
	cfg         Config
	server      services.FileServer
	waitGroup   *sync.WaitGroup
	collectorWg sync.WaitGroup
	fileChannel chan models.File
	quit        chan bool
}

func New(
	processID int, config Config, fileServer services.FileServer, fileCh chan models.File, waitGroup *sync.WaitGroup,
) (*Collector, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Collector{
		ID:          processID,
		cfg:         config,
		server:      fileServer,
		fileChannel: fileCh,
		waitGroup:   waitGroup,
		quit:        make(chan bool),
		collectorWg: sync.WaitGroup{},
	}, nil
}

func (c *Collector) Start() {
	logger.Infof("[Collector %d] Starting files collect", c.ID)

	go func() {
		for {
			select {
			case <-c.quit:
				return
			default:
				for _, pattern := range c.cfg.MatchPatterns {
					c.collectorWg.Add(1)

					go c.collectFiles(pattern)
				}

				time.Sleep(time.Second * time.Duration(c.cfg.CollectDelay))
				c.collectorWg.Wait()
				c.waitGroup.Wait()
			}
		}
	}()
}

func (c *Collector) Stop() {
	go func() {
		logger.Infof("[Collector %d] Stopping file collect..", c.ID)
		c.quit <- true
	}()
}

// Responsabilidades do collector
// Coletar os arquivos da fonte e enviar através do Channel.
func (c *Collector) collectFiles(pattern string) {
	defer c.collectorWg.Done()

	ctx, span := trace.NewSpan(context.Background(), "collector.collectFiles")
	defer span.End()

	trace.AddSpanTags(span, map[string]string{"pattern": pattern})
	logger.Infof("[Collector %d] Collecting files with pattern: %s", c.ID, pattern)

	collectedFiles, err := c.server.Glob(ctx, pattern)
	if err != nil {
		trace.AddSpanError(span, err)
		trace.FailSpan(span, fmt.Sprintf("Error on collect files with pattern: %s", pattern))
		logger.Errorf("[Collector %d] Error on collect files with pattern %s, %s", c.ID, pattern, err)

		return
	}

	if len(collectedFiles) > 0 {
		trace.AddSpanTags(span, map[string]string{"filesCount": strconv.Itoa(len(collectedFiles))})
	}

	sendedCount := 0

	for _, fp := range collectedFiles {
		model, err := c.createFileModel(fp)
		if err != nil {
			trace.AddSpanError(span, err)
			trace.FailSpan(span, "Failed to create FileModel")
			logger.Errorf("[Collector %d] Failed to create FileModel, %s", c.ID, err)

			continue
		}

		c.waitGroup.Add(1)
		c.fileChannel <- model
		logger.Debugf("[Collector %d] File collected %+v", c.ID, model)

		sendedCount++
		if c.cfg.MaxCollectBatchSize > 0 && sendedCount == c.cfg.MaxCollectBatchSize {
			return
		}
	}
}

// Insert a timestamp at end of file name maintaining same file extension.
func (c *Collector) createFileModel(filePath string) (models.File, error) {
	fileName := path.Base(filePath)

	return models.NewFile(fileName, filePath, fileName, c.server)
}
