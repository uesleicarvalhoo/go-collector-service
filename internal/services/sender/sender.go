package sender

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
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

func (s *Sender) Consume(collector Collector) {
	logrus.Info("Start to collect files")

	for {
		files, err := collector.GetFiles()
		if err != nil {
			continue
		}

		for _, file := range files {
			key, err := s.ProcessFile(context.Background(), file)
			if err == nil {
				logrus.Infof("File published at '%s'", key)
			} else {
				logrus.Errorf("Error on publish file '%s': '%s'", file, err)
			}
		}

		time.Sleep(time.Microsecond * 100)
	}
}

func (s *Sender) ProcessFile(ctx context.Context, file models.File) (string, error) {
	ctx, span := trace.NewSpan(ctx, "sender.process_file")
	defer span.End()

	fileKey, err := s.PublishFile(ctx, s.getFileKey(file), file)
	if err != nil {
		return "", err
	}

	return fileKey, s.RemoveFile(ctx, file)
}

func (s *Sender) PublishFile(ctx context.Context, fileKey string, file models.File) (string, error) {
	span := trace.SpanFromContext(ctx)

	reader, err := file.GetReader()
	if err != nil {
		logrus.Infof("Error on get file reader, %s\n", err)

		return "", err
	}
	defer reader.Close()

	err = s.storage.SendFile(ctx, fileKey, reader)
	if err != nil {
		logrus.Errorf("Error on sendfile, %s\n", err)
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
		logrus.Errorf("Error on publish event %s\n", err)
		trace.AddSpanError(span, err)

		return "", err
	}

	return fileKey, nil
}

func (s *Sender) RemoveFile(ctx context.Context, file models.File) error {
	span := trace.SpanFromContext(ctx)

	data := map[string]string{
		"filename": file.Name,
		"filepath": file.FilePath,
	}

	trace.AddSpanEvents(span, "sender.remove_file", data)

	if err := file.Delete(); err != nil {
		trace.AddSpanError(span, err)

		return err
	}

	return nil
}

func (s *Sender) getFileKey(file models.File) string {
	ext := filepath.Ext(file.Name)
	baseName := strings.TrimSuffix(file.Name, ext)

	return fmt.Sprintf("%s-%s%s", baseName, time.Now().Format(time.RFC3339Nano), ext)
}
