package sender

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/streamer"

	"github.com/uesleicarvalhoo/go-collector-service/pkg/broker"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/fileserver"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/storage"
)

var (
	tmpDir       string
	memoryBroker *broker.MemoryBroker
)

func init() {
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	tmpDir = folder
}

func createDefaultDirTempFile(fileName string) (models.File, error) {
	return createTempFile(tmpDir, fileName)
}

func createTempFile(dir, fileName string) (models.File, error) {
	fp := filepath.Join(dir, fileName)

	err := ioutil.WriteFile(fp, []byte{}, 0o644)
	if err != nil {
		panic(err)
	}

	return models.NewFile(fileName, fp, fileName)
}

func newSut() *Sender {
	brokerService := broker.NewMemoryBroker()
	streamerService, err := streamer.NewStreamer(brokerService, broker.CreateTopicInput{Name: "collector.files"})
	if err != nil {
		panic(err)
	}

	fs, err := fileserver.NewLocalFileServer(fileserver.Config{})
	if err != nil {
		panic(err)
	}

	memoryBroker = brokerService

	return NewSender(streamerService, storage.NewMemoryStorage(), fs)
}

func TestPublishFileSendFileToStorage(t *testing.T) {
	// Prepare
	sut := newSut()
	memoryStorage := sut.storage.(*storage.MemoryStorage)

	// Arrange
	file, err := createDefaultDirTempFile("test_publish_file_send_file_to_storage.json")
	assert.Nil(t, err)

	// Action
	_, err = sut.PublishFile(context.TODO(), file)
	assert.Nil(t, err)

	// Arrange
	assert.True(t, memoryStorage.FileExists(file.Key))
}

func TestPublishFileSendEventToStreamer(t *testing.T) {
	// Prepare
	sut := newSut()

	// Arrange
	file, err := createDefaultDirTempFile("test_publish_file_send_file_to_storage.json")
	assert.Nil(t, err)

	expectedEvent := models.Event{
		Topic: "collector.files",
		Key:   "published",
		Data:  map[string]string{"file_key": file.Key},
	}

	// Action
	_, err = sut.PublishFile(context.TODO(), file)
	assert.Nil(t, err)

	// Assert
	brokerEvents, ok := memoryBroker.Events["collector.files"]
	assert.True(t, ok)

	sendedEvent := brokerEvents[len(brokerEvents)-1]

	assert.Equal(t, sendedEvent, expectedEvent)
}

func TestConsumeSendAllFilesToStorage(t *testing.T) {
	// Prepare
	sut := newSut()

	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	memoryStorage := sut.storage.(*storage.MemoryStorage)
	assert.Empty(t, memoryStorage.GetAllFiles())

	testFile1, err := createTempFile(folder, "test_file_1.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile1.FilePath)

	testFile2, err := createTempFile(folder, "test_file_2.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile2.FilePath)

	// Arrange
	pattern := filepath.Join(folder, "*.json")

	sut.Start(pattern)
	time.Sleep(time.Second * 1)

	// Assert
	storedFiles := memoryStorage.GetAllFiles()
	assert.Len(t, storedFiles, 2)
}
