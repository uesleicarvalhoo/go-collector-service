package collector

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"sync"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

type Collector struct {
	ID           int
	cfg          Config
	server       services.FileServer
	collectGroup *sync.WaitGroup
	processGroup *sync.WaitGroup
}

func New(
	processID int,
	config Config,
	fileServer services.FileServer,
	collectWaitGroup *sync.WaitGroup,
	proccessGroup *sync.WaitGroup,
) (*Collector, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Collector{
		ID:           processID,
		cfg:          config,
		server:       fileServer,
		collectGroup: collectWaitGroup,
		processGroup: proccessGroup,
	}, nil
}

func (c *Collector) CollectFiles(ctx context.Context, channel chan models.File) {
	ctx, span := trace.NewSpan(ctx, "collector.collectFiles")
	defer span.End()

	for _, pattern := range c.cfg.MatchPatterns {
		c.collectGroup.Add(1)

		go c.collectFilesWithPattern(ctx, channel, pattern)
	}
}

// Responsabilidades do collector
// Coletar os arquivos da fonte e enviar através do Channel.
func (c *Collector) collectFilesWithPattern(
	ctx context.Context, channel chan models.File, pattern string,
) {
	defer c.collectGroup.Done()

	span := trace.SpanFromContext(ctx)

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
		trace.AddSpanTags(span, map[string]string{"matchCount": strconv.Itoa(len(collectedFiles))})
	}

	sendedCount := 0

	defer func() {
		trace.AddSpanTags(span, map[string]string{"sendedCount": strconv.Itoa(sendedCount)})
	}()

	for _, fp := range collectedFiles {
		model, err := c.createFileModel(fp)
		if err != nil {
			trace.AddSpanError(span, err)
			trace.FailSpan(span, "Failed to create FileModel")
			logger.Errorf("[Collector %d] Failed to create FileModel, %s", c.ID, err)

			continue
		}

		c.processGroup.Add(1)
		channel <- model

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
