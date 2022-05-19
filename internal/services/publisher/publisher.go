package publisher

import (
	"context"
	"path"
	"strconv"
	"sync"

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
	eventChannel chan models.Event
}

func New(
	publisherID int,
	eventTopic string,
	storage services.Storage,
	eventCh chan models.Event,
	waitGroup *sync.WaitGroup,
) *Publisher {
	return &Publisher{
		ID:           publisherID,
		EventTopic:   eventTopic,
		storage:      storage,
		waitGroup:    waitGroup,
		eventChannel: eventCh,
	}
}

func (p *Publisher) Handle(ctx context.Context, fileChannel chan models.File) {
	go func() {
		for file := range fileChannel {
			err := p.processFile(ctx, file)
			if err == nil {
				logger.Infof("[Publisher %d] File %+v uploaded with success", p.ID, file.FileInfo)
				p.notifyResult("success", map[string]string{"file_key": file.Key})
			} else {
				logger.Errorf("[Publisher %d] Failed to upload file '%+v', %s", p.ID, file.FileInfo, err)
				p.notifyResult("error", map[string]string{"file_path": file.FilePath, "error": err.Error()})
			}
		}
	}()
}

func (p *Publisher) processFile(ctx context.Context, file models.File) error {
	defer p.waitGroup.Done()

	ctx, span := trace.NewSpan(ctx, "publisher.processFile")
	defer span.End()

	trace.AddSpanTags(
		span,
		map[string]string{
			"fileKey":     file.Key,
			"fileName:":   file.Name,
			"filePath:":   file.FilePath,
			"fileSize":    strconv.Itoa(int(file.Size)),
			"fileModTime": file.ModTime.Format("2006-01-02 15:04:05"),
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

	if file.Size == 0 {
		return ErrEmptyFile
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

// Move file to ./sent/<filename>.
func (p *Publisher) moveFile(ctx context.Context, file models.File) {
	span := trace.SpanFromContext(ctx)

	fileDir, fileName := path.Split(file.FilePath)
	newPath := path.Join(fileDir, "sent", fileName)

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
