package sender

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

type Sender struct {
	streamer Streamer
	storage  Storage
}

func NewSender(streamer Streamer, storage Storage) *Sender {
	return &Sender{
		streamer: streamer,
		storage:  storage,
	}
}

// Get files from collector and send file to Storage
// if it was success, send notification to MessageBroker and delete file.
func (s *Sender) Consume(collector Collector) {
	logger.Info("Start to collect files")

	for {
		files, err := collector.GetFiles()
		if err != nil {
			continue
		}

		for _, file := range files {
			ctx := context.Background()

			key, err := s.PublishFile(ctx, s.getFileKey(file), file)
			if err == nil {
				logger.Infof("File published at '%s'", key)
			} else {
				logger.Errorf("Error on publish file '%s': '%s'", file, err)
			}

			_ = s.RemoveFile(ctx, file, collector)
		}

		time.Sleep(time.Microsecond * 100)
	}
}

// Publish File at Storage and if it was success, send a message with fileKey to MessageBroker.
func (s *Sender) PublishFile(ctx context.Context, fileKey string, file models.File) (string, error) {
	span := trace.SpanFromContext(ctx)

	reader, err := file.GetReader()
	if err != nil {
		logger.Infof("Error on get file reader, %s\n", err)

		return "", err
	}
	defer reader.Close()

	err = s.storage.SendFile(ctx, fileKey, reader)
	if err != nil {
		logger.Errorf("Error on sendfile, %s\n", err)
		trace.AddSpanError(span, err)

		return "", err
	}

	data := map[string]string{
		"key":   fileKey,
		"name:": file.Name,
		"path:": file.FilePath,
	}

	trace.AddSpanEvents(span, "publish_file", data)

	err = s.streamer.NotifyPublishedFile(fileKey, file)
	if err != nil {
		logger.Errorf("Error on publish event %s\n", err)
		trace.AddSpanError(span, err)

		return "", err
	}

	return fileKey, nil
}

// Delete file from origin.
func (s *Sender) RemoveFile(ctx context.Context, file models.File, collector Collector) error {
	span := trace.SpanFromContext(ctx)

	data := map[string]string{
		"filename": file.Name,
		"filepath": file.FilePath,
	}

	trace.AddSpanEvents(span, "sender.remove_file", data)

	if err := collector.RemoveFile(file); err != nil {
		trace.AddSpanError(span, err)

		return err
	}

	return nil
}

// Insert a timestamp at end of file name maintaining same file extension.
func (s *Sender) getFileKey(file models.File) string {
	ext := filepath.Ext(file.Name)
	baseName := strings.TrimSuffix(file.Name, ext)

	return fmt.Sprintf("%s-%s%s", baseName, time.Now().Format(time.RFC3339Nano), ext)
}
