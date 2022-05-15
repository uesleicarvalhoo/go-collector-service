package publisher

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

type Publisher struct {
	ID           int
	EventTopic   string
	storage      services.Storage
	waitGroup    *sync.WaitGroup
	fileChannel  <-chan models.File
	eventChannel chan models.Event
	quit         chan bool
}

func New(
	publisherID int,
	eventTopic string,
	storage services.Storage,
	fileCh <-chan models.File,
	eventCh chan models.Event,
	waitGroup *sync.WaitGroup,
) *Publisher {
	return &Publisher{
		ID:           publisherID,
		EventTopic:   eventTopic,
		storage:      storage,
		fileChannel:  fileCh,
		waitGroup:    waitGroup,
		eventChannel: eventCh,
	}
}

func (p *Publisher) Start() {
	go func() {
		for {
			select {
			case file := <-p.fileChannel:
				logger.Debugf("[Publisher %d] File %+v received", p.ID, file)

				err := p.processFile(context.Background(), file)
				if err == nil {
					p.notifyResult("success", map[string]string{"file_key": file.Key})
					logger.Info("[Publisher %d] File %s uploaded with success", p.ID, file.FilePath)
				} else {
					p.notifyResult("error", map[string]string{"file_path": file.FilePath, "error": err.Error()})
					logger.Errorf("[Publisher %d] Failed to upload file, %s", p.ID, err)
				}

			case <-p.quit:
				logger.Info("[Publisher %d] Quit signal received, stopping worker..", p.ID)

				return
			}
		}
	}()
}

func (p *Publisher) Stop() {
	go func() {
		logger.Infof("[Publisher %d] Stopping to publish files..", p.ID)
		p.quit <- true
	}()
}

func (p *Publisher) processFile(ctx context.Context, file models.File) error {
	defer p.waitGroup.Done()

	ctx, span := trace.NewSpan(ctx, "publisher.processFile")
	defer span.End()

	trace.AddSpanTags(
		span,
		map[string]string{
			"fileKey":   file.Key,
			"fileName:": file.Name,
			"filePath:": file.FilePath,
		},
	)

	err := file.Lock(ctx)
	if err != nil {
		logger.Errorf("[Publisher %d] Error on acquire file lock '%s': '%s'", p.ID, file.FilePath, err)
		trace.AddSpanTags(span, map[string]string{"result": "lock-error"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Error on acquire file lock")

		return err
	}

	err = p.publishFile(ctx, file)
	if err != nil {
		logger.Errorf("[Publisher %d] Error on publish file '%s': '%s'", p.ID, file.FilePath, err)
		trace.AddSpanTags(span, map[string]string{"result": "fail"})
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Error on publish file")

		return err
	}

	_ = file.Unlock(ctx)

	p.moveFile(ctx, file)

	return nil
}

// Publish File at Storage.
func (p *Publisher) publishFile(ctx context.Context, file models.File) error {
	span := trace.SpanFromContext(ctx)

	trace.AddSpanEvents(
		span,
		"publisher.publishFile",
		map[string]string{
			"filename": file.Name,
			"filepath": file.FilePath,
		})

	reader, err := file.Open(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	err = p.storage.SendFile(ctx, file.Key, reader)
	if err != nil {
		return err
	}

	return nil
}

// Move file to ./send<filename><timestamp><fileextension>.
func (p *Publisher) moveFile(ctx context.Context, file models.File) {
	span := trace.SpanFromContext(ctx)

	fileDir, name := path.Split(file.FilePath)
	ext := path.Ext(name)
	baseName := strings.TrimSuffix(file.Name, ext)
	newName := fmt.Sprintf("%s-%s%s", baseName, time.Now().Format(time.RFC3339Nano), ext)
	newPath := path.Join(fileDir, "sent", newName)

	trace.AddSpanEvents(
		span,
		"publisher.moveFile",
		map[string]string{
			"filename": file.Name,
			"filepath": file.FilePath,
			"newpath":  newPath,
		})

	if err := file.Move(ctx, newPath); err != nil {
		trace.AddSpanError(span, err)
		trace.FailSpan(span, "Failed to move file")
		logger.Errorf("Failed to move file, %s", err)

		return
	}
}

func (p *Publisher) notifyResult(result string, data any) {
	event, err := models.NewEvent(p.EventTopic, result, data)
	if err != nil {
		logger.Errorf("Failed to create event, %s", err)
	}
	p.eventChannel <- event
}
