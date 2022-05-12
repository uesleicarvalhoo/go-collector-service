package sender

import (
	"context"
	"errors"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/broker"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/fileserver"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/storage"
)

func createTempFile(dir, fileName string) (models.File, error) {
	if dir == "" {
		tmpDir, err := ioutil.TempDir("", "*")
		if err != nil {
			panic(err)
		}

		dir = tmpDir
	}

	fp := filepath.Join(dir, fileName)

	err := ioutil.WriteFile(fp, []byte{}, fs.ModePerm)
	if err != nil {
		panic(err)
	}

	return models.NewFile(fileName, fp, fileName)
}

func newSut(patterns ...string) *Sender {
	if len(patterns) == 0 {
		patterns = append(patterns, "")
	}

	brokerService, err := broker.NewMemoryBroker()
	if err != nil {
		panic(err)
	}

	fs, err := fileserver.NewLocalFileServer(fileserver.Config{})
	if err != nil {
		panic(err)
	}

	cfg := Config{
		ParalelUploads: 1,
		EventTopic:     "collector.files",
		MatchPatterns:  patterns,
		CollectDelay:   1,
	}

	sender, err := New(cfg, storage.NewMemoryStorage(), brokerService, fs)
	if err != nil {
		panic(err)
	}

	return sender
}

func TestPublishFileSendFileToStorage(t *testing.T) {
	// Prepare
	sut := newSut()
	memoryStorage, _ := sut.storage.(*storage.MemoryStorage)

	// Arrange
	file, err := createTempFile("", "test_publish_file_send_file_to_storage.json")
	assert.Nil(t, err)

	// Action
	err = sut.publishFile(context.TODO(), file)
	assert.Nil(t, err)

	// Arrange
	assert.True(t, memoryStorage.FileExists(file.Key))
}

func TestProcessFileSendEventToStreamer(t *testing.T) {
	// Prepare
	sut := newSut()
	memoryBroker, _ := sut.broker.(*broker.MemoryBroker)

	// Arrange
	file, err := createTempFile("", "test_publish_file_send_file_to_storage.json")
	assert.Nil(t, err)

	expectedEvent := models.Event{
		Topic: sut.cfg.EventTopic,
		Key:   "published",
		Data:  map[string]string{"file_key": file.Key},
	}

	// Action
	err = sut.processFile(context.TODO(), file)
	assert.Nil(t, err)

	// Assert
	brokerEvents, ok := memoryBroker.Events[sut.cfg.EventTopic]
	assert.True(t, ok)

	sendedEvent := brokerEvents[len(brokerEvents)-1]

	assert.Equal(t, expectedEvent, sendedEvent)
}

func TestNotifyPublishedFileSendEventToStreamer(t *testing.T) {
	// Prepare
	sut := newSut()
	memoryBroker, _ := sut.broker.(*broker.MemoryBroker)

	// Arrange
	file, err := createTempFile("", "test_notify_published_file_send_event_to_streamer.json")
	assert.Nil(t, err)

	expectedEvent := models.Event{
		Topic: sut.cfg.EventTopic,
		Key:   "published",
		Data:  map[string]string{"file_key": file.Key},
	}

	// Action
	sut.notifyPublishedFile(context.TODO(), file)

	// Assert
	brokerEvents, ok := memoryBroker.Events[sut.cfg.EventTopic]
	assert.True(t, ok)

	sendedEvent := brokerEvents[len(brokerEvents)-1]

	assert.Equal(t, expectedEvent, sendedEvent)
}

func TestNotifyPublishedFileErrorSendEventToStreamer(t *testing.T) {
	// Prepare
	sut := newSut()
	memoryBroker, _ := sut.broker.(*broker.MemoryBroker)
	expetecdErrorMsg := "Invalid File"

	// Arrange
	file, err := createTempFile("", "test_notify_published_file_error_send_event_to_streamer.json")
	assert.Nil(t, err)

	expectedEvent := models.Event{
		Topic: sut.cfg.EventTopic,
		Key:   "error",
		Data:  map[string]string{"file_path": file.FilePath, "error": expetecdErrorMsg},
	}

	// Action
	sut.notifyPublishFileError(context.TODO(), file, errors.New(expetecdErrorMsg))

	// Assert
	brokerEvents, ok := memoryBroker.Events[sut.cfg.EventTopic]
	assert.True(t, ok)

	sendedEvent := brokerEvents[len(brokerEvents)-1]

	assert.Equal(t, expectedEvent, sendedEvent)
}

func TestConsumeSendAllFilesToStorage(t *testing.T) {
	// Prepare
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	pattern := filepath.Join(folder, "*.json")
	sut := newSut(pattern)

	memoryStorage, _ := sut.storage.(*storage.MemoryStorage)
	assert.Empty(t, memoryStorage.GetAllFiles())

	testFile1, err := createTempFile(folder, "test_file_1.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile1.FilePath)

	testFile2, err := createTempFile(folder, "test_file_2.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile2.FilePath)

	// Arrange

	go sut.Start()
	time.Sleep(time.Second * 2)

	// Assert
	storedFiles := memoryStorage.GetAllFiles()
	assert.Len(t, storedFiles, 2)
}
