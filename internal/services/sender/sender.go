package sender

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

type Sender struct {
	streamer        Streamer
	storage         Storage
	fileServer      FileServer
	pendingFiles    chan models.File
	processingFiles map[string]bool
	sync.Mutex
}

func NewSender(streamer Streamer, storage Storage, fileServer FileServer) *Sender {
	return &Sender{
		streamer:        streamer,
		storage:         storage,
		fileServer:      fileServer,
		pendingFiles:    make(chan models.File, 50),
		processingFiles: map[string]bool{},
	}
}

// Get files from FileServer and send file to Storage
// if it was success, send notification to MessageBroker and delete file.
func (s *Sender) Start(patterns ...string) {
	go s.process()

	for _, pattern := range patterns {
		go s.collect(pattern)
	}
}

func (s *Sender) collect(pattern string) {
	ctx, span := trace.NewSpan(context.Background(), "sender.process_files")
	defer span.End()

	collectedFiles, err := s.fileServer.Glob(ctx, pattern)
	if err != nil {
		return
	}

	for _, fp := range collectedFiles {
		model, err := s.createFileModel(fp)
		if err != nil {
			logger.Errorf("Failed to create FileModel, %s\n", err)
			trace.AddSpanError(span, err)

			continue
		}

		if _, ok := s.processingFiles[model.FilePath]; ok {
			continue
		}

		s.pendingFiles <- model
		s.processingFiles[model.FilePath] = true
	}

	// TODO: Avaliar se precisa colocar algum intervalo entre as coletas
}

func (s *Sender) process() {
	for file := range s.pendingFiles {
		ctx := context.TODO()

		err := s.ProcessFile(ctx, file)
		if err != nil {
			if err = s.streamer.NotifyInvalidFile(ctx, file); err != nil {
				logger.Errorf("Erro on notify invalid file, %v", err)
			}
		}
	}
}

func (s *Sender) ProcessFile(ctx context.Context, file models.File) error {
	ctx, span := trace.NewSpan(ctx, "sender.process_file")
	defer span.End()

	trace.AddSpanTags(span, map[string]string{"fileKey": file.Key})

	key, err := s.PublishFile(ctx, file)
	if err != nil {
		logger.Errorf("Error on publish file '%s': '%s'", file.FilePath, err)

		return err
	}

	logger.Infof("File published at '%s'", key)

	return s.RemoveFile(ctx, file)
}

// Publish File at Storage and if it was success, send a message with fileKey to MessageBroker.
func (s *Sender) PublishFile(ctx context.Context, file models.File) (string, error) {
	span := trace.SpanFromContext(ctx)

	reader, err := s.fileServer.Open(ctx, file.FilePath)
	if err != nil {
		logger.Errorf("Error on get file reader, %s\n", err)
		trace.AddSpanError(span, err)

		return "", err
	}
	defer reader.Close()

	err = s.storage.SendFile(ctx, file.Key, reader)
	if err != nil {
		logger.Errorf("Error on sendfile, %s\n", err)
		trace.AddSpanError(span, err)

		return "", err
	}

	trace.AddSpanEvents(
		span,
		"publish_file",
		map[string]string{
			"key":   file.Key,
			"name:": file.Name,
			"path:": file.FilePath,
		})

	err = s.streamer.NotifyPublishedFile(ctx, file.Key, file)
	if err != nil {
		logger.Errorf("Error on publish event %s\n", err)
		trace.AddSpanError(span, err)

		return "", err
	}

	return file.Key, nil
}

// Delete file from origin.
func (s *Sender) RemoveFile(ctx context.Context, file models.File) error {
	span := trace.SpanFromContext(ctx)

	trace.AddSpanEvents(
		span,
		"sender.remove_file",
		map[string]string{
			"filename": file.Name,
			"filepath": file.FilePath,
		})

	if err := s.fileServer.Remove(ctx, file.FilePath); err != nil {
		trace.AddSpanError(span, err)

		return err
	}

	s.Lock()
	defer s.Unlock()

	delete(s.processingFiles, file.FilePath)

	return nil
}

func (s *Sender) getFiles(ctx context.Context, pattern string) ([]models.File, error) {
	_, span := trace.NewSpan(ctx, "list-files")
	defer span.End()

	files := make([]models.File, 0)

	collectedFiles, err := s.fileServer.Glob(ctx, pattern)
	if err != nil {
		return nil, err
	}

	for _, fp := range collectedFiles {
		model, err := s.createFileModel(fp)
		if err != nil {
			logger.Errorf("Failed to create FileModel, %s\n", err)
			trace.AddSpanError(span, err)

			continue
		}

		files = append(files, model)
	}

	return files, nil
}

// Insert a timestamp at end of file name maintaining same file extension.
func (s *Sender) createFileModel(filePath string) (models.File, error) {
	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, ext)

	key := fmt.Sprintf("%s-%s%s", baseName, time.Now().Format(time.RFC3339Nano), ext)

	return models.NewFile(fileName, filePath, key)
}
