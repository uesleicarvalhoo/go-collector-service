package sender

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/trace"
)

const processChannelBuffer = 100

type Sender struct {
	broker         Broker
	storage        Storage
	fileServer     FileServer
	cfg            Config
	processChannel chan models.File
	processWg      sync.WaitGroup
	collectWg      sync.WaitGroup
	isRunning      bool
	sync.Mutex     // Used for shutdown
}

func New(cfg Config, storage Storage, broker Broker, fileServer FileServer) (*Sender, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	sender := &Sender{
		cfg:            cfg,
		broker:         broker,
		storage:        storage,
		fileServer:     fileServer,
		processWg:      sync.WaitGroup{},
		collectWg:      sync.WaitGroup{},
		processChannel: make(chan models.File, processChannelBuffer),
	}

	for i := 0; i < sender.cfg.ParalelUploads; i++ {
		go sender.fileProcessWorker()
	}

	return sender, nil
}

// Get files from FileServer and send file to Storage
// if it was success, send notification to MessageBroker and delete file.
func (s *Sender) Start() {
	if s.isRunning {
		return
	}

	s.isRunning = true
	logger.Infof(
		"Starting file collection with %d workers and bind %d pattenrs", s.cfg.ParalelUploads, len(s.cfg.MatchPatterns),
	)

	for s.isRunning {
		for _, pattern := range s.cfg.MatchPatterns {
			s.collectWg.Add(1)

			go s.collectFiles(pattern)
		}

		s.collectWg.Wait()
		s.processWg.Wait()
		time.Sleep(time.Second * time.Duration(s.cfg.CollectDelay))
	}

	logger.Infof("File collection stoped.")
}

// Close Channels and stop to collect files.
func (s *Sender) Shutdown() {
	logger.Info("Stopping sender..")
	s.Lock()
	defer s.Unlock()

	s.isRunning = false
	close(s.processChannel)
}

// Collect files from FileServer that match with <pattern> and send to processChannel
// each file just is sended one time.
func (s *Sender) collectFiles(pattern string) {
	defer s.collectWg.Done()

	ctx, span := trace.NewSpan(context.Background(), "sender.collectFiles")
	defer span.End()

	trace.AddSpanTags(span, map[string]string{"pattern": pattern})
	logger.Infof("Collecting files with pattern: %s", pattern)

	collectedFiles, err := s.fileServer.Glob(ctx, pattern)
	if err != nil {
		trace.AddSpanError(span, err)
		logger.Errorf("Error on collect files with pattern %s, %s", pattern, err)

		return
	}

	sendedCount := 0

	for _, fp := range collectedFiles {
		model, err := s.createFileModel(fp)
		if err != nil {
			trace.AddSpanError(span, err)
			logger.Errorf("Failed to create FileModel, %s", err)

			continue
		}

		s.processWg.Add(1)
		s.processChannel <- model

		sendedCount++
		if s.cfg.MaxCollectBatchSize > 0 && sendedCount == s.cfg.MaxCollectBatchSize {
			return
		}
	}
}

// Publish File at Storage.
func (s *Sender) publishFile(ctx context.Context, file models.File) error {
	span := trace.SpanFromContext(ctx)

	trace.AddSpanEvents(
		span,
		"sender.publishFile",
		map[string]string{
			"filename": file.Name,
			"filepath": file.FilePath,
		})

	reader, err := s.fileServer.Open(ctx, file.FilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	err = s.storage.SendFile(ctx, file.Key, reader)
	if err != nil {
		return err
	}

	return nil
}

// Move file to ./send<filename><timestamp><fileextension>.
func (s *Sender) moveFile(ctx context.Context, file models.File) {
	span := trace.SpanFromContext(ctx)

	trace.AddSpanEvents(
		span,
		"sender.moveFile",
		map[string]string{
			"filename": file.Name,
			"filepath": file.FilePath,
		})

	fileDir, name := filepath.Split(file.FilePath)
	ext := filepath.Ext(name)
	baseName := strings.TrimSuffix(file.Name, ext)
	newName := fmt.Sprintf("%s-%s%s", baseName, time.Now().Format(time.RFC3339Nano), ext)
	moveDir := filepath.Join(fileDir, "sent", newName)

	if err := s.fileServer.MoveFile(ctx, file.FilePath, moveDir); err != nil {
		trace.AddSpanError(span, err)
		logger.Errorf("Failed to move file, %s", err)

		return
	}
}

// Receive files from processChannel and process file.
func (s *Sender) fileProcessWorker() {
	for file := range s.processChannel {
		_ = s.processFile(context.Background(), file)
		s.processWg.Done()
	}
}

// Publish file, and send a message to broker Event{Key: <success/fail>, Data: {"file_key": <file.Key>}}.
func (s *Sender) processFile(ctx context.Context, file models.File) error {
	ctx, span := trace.NewSpan(ctx, "sender.processFile")
	defer span.End()

	trace.AddSpanTags(
		span,
		map[string]string{
			"fileKey":   file.Key,
			"fileName:": file.Name,
			"filePath:": file.FilePath,
		},
	)

	lock, err := s.fileServer.Lock(ctx, file.FilePath)
	if err != nil {
		logger.Errorf("Error on acquire file lock '%s': '%s'", file.FilePath, err)
		trace.AddSpanTags(span, map[string]string{"result": "lock-error"})
		trace.AddSpanError(span, err)

		return err
	}

	err = s.publishFile(ctx, file)
	if err != nil {
		logger.Errorf("Error on publish file '%s': '%s'", file.FilePath, err)
		trace.AddSpanTags(span, map[string]string{"result": "fail"})
		trace.AddSpanError(span, err)
		s.notifyPublishFileError(ctx, file, err)

		return err
	}

	_ = lock.Unlock()

	s.notifyPublishedFile(ctx, file)
	trace.AddSpanTags(span, map[string]string{"result": "success"})
	logger.Infof("File published at '%s'", file.Key)

	s.moveFile(ctx, file)

	return nil
}

// Insert a timestamp at end of file name maintaining same file extension.
func (s *Sender) createFileModel(filePath string) (models.File, error) {
	fileName := filepath.Base(filePath)

	return models.NewFile(fileName, filePath, fileName)
}

// Send event to Topic with RoutingKey "published" and body {"file_key": "<file.Key>"}.
func (s *Sender) notifyPublishedFile(ctx context.Context, file models.File) {
	routingKey := "published"
	span := trace.SpanFromContext(ctx)
	trace.AddSpanEvents(
		span,
		"sender.notifyPublishedFile",
		map[string]string{"topic": s.cfg.EventTopic, "routing-key": routingKey},
	)

	event, err := models.NewEvent(s.cfg.EventTopic, routingKey, map[string]string{"file_key": file.Key})
	if err != nil {
		logger.Errorf("Invalid event, %s", err)
	}

	err = s.broker.SendEvent(event)
	if err != nil {
		logger.Errorf("Failed to send event, %s", err)
	}
}

// Send event to topic with RoutingKey "error" and body {"file_path": "<file.FilePath>", "error": <err.Error()>}.
func (s *Sender) notifyPublishFileError(ctx context.Context, file models.File, err error) {
	routingKey := "error"
	data := map[string]string{
		"file_path": file.FilePath,
		"error":     err.Error(),
	}

	span := trace.SpanFromContext(ctx)
	trace.AddSpanEvents(
		span,
		"sender.notifyPublishedFileError",
		map[string]string{"topic": s.cfg.EventTopic, "routing-key": routingKey},
	)

	event, err := models.NewEvent(s.cfg.EventTopic, routingKey, data)
	if err != nil {
		logger.Errorf("Invalid event, %s", err)
	}

	err = s.broker.SendEvent(event)
	if err != nil {
		logger.Errorf("Failed to send event, %s", err)
	}
}