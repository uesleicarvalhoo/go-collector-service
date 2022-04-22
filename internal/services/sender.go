package services

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

type Sender struct {
	collector Collector
	streamer  Streamer
	storage   Storage
}

func NewSender(streamer Streamer, storage Storage, collector Collector) *Sender {
	return &Sender{
		streamer:  streamer,
		storage:   storage,
		collector: collector,
	}
}

func (s *Sender) Run() {
	for {
		files, err := s.collector.ListFiles()
		if err != nil {
			continue
		}

		for _, file := range files {
			s.ProcessFile(context.Background(), file)
		}

		time.Sleep(time.Microsecond * 100)
	}
}

func (s *Sender) ProcessFile(ctx context.Context, fileInfo models.FileInfo) error {
	ctx, span := trace.NewSpan(ctx, "sender.process_file")
	defer span.End()

	body, err := s.collector.GetFileData(fileInfo)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		log.Printf("File with no data '%s', removed", fileInfo.Name)
		s.RemoveFile(ctx, fileInfo)

		return ErrNoData
	}

	file := models.File{
		FileInfo: fileInfo,
		Key:      s.getFileKey(fileInfo),
		Data:     body,
	}

	err = s.PublishFile(ctx, file)
	if err != nil {
		return err
	}

	return s.RemoveFile(ctx, file.FileInfo)
}

func (s *Sender) PublishFile(ctx context.Context, file models.File) error {
	span := trace.SpanFromContext(ctx)

	err := s.storage.SendFile(ctx, file)
	if err != nil {
		log.Printf("Error on sendfile, %s\n", err)
		trace.AddSpanError(span, err)

		return err
	}

	data := map[string]string{
		"key":   file.Key,
		"name:": file.Name,
		"path:": file.FilePath,
	}

	trace.AddSpanEvents(span, "publish_file", data)

	err = s.streamer.NotifyPublishedFile(file)
	if err != nil {
		trace.AddSpanError(span, err)

		return err
	}

	return nil
}

func (s *Sender) RemoveFile(ctx context.Context, file models.FileInfo) error {
	span := trace.SpanFromContext(ctx)

	data := map[string]string{
		"filename": file.Name,
		"filepath": file.FilePath,
	}

	trace.AddSpanEvents(span, "sender.remove_file", data)

	err := s.collector.RemoveFile(file)
	if err != nil {
		trace.AddSpanError(span, err)

		return err
	}

	return nil
}

func (s *Sender) getFileKey(file models.FileInfo) string {
	ext := filepath.Ext(file.Name)
	baseName := strings.TrimSuffix(file.Name, ext)

	return fmt.Sprintf("%s-%s%s", baseName, time.Now().Format(time.RFC3339Nano), ext)
}
