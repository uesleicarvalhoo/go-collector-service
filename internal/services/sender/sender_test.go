package sender

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/collector"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/streamer"

	"github.com/uesleicarvalhoo/go-collector-service/pkg/broker"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/storage"
)

var tmpDir string

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

	lc, err := collector.NewLocalCollector()
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fp, []byte{}, 0o644)
	if err != nil {
		panic(err)
	}

	return models.NewFile(filepath.Base(fp), fp, lc), nil
}

func newSut() *Sender {
	sut, _ := newSutWithBroker()
	return sut
}

func newSutWithBroker() (*Sender, *broker.MemoryBroker) {
	memoryBroker := broker.NewMemoryBroker()
	streamerService, err := streamer.NewStreamer(memoryBroker, broker.CreateTopicInput{Name: "collector.files"})
	if err != nil {
		panic(err)
	}

	return NewSender(streamerService, storage.NewMemoryStorage()), memoryBroker
}

func TestPublishFileSendFileToStorage(t *testing.T) {
	// Prepare
	sut := newSut()
	memoryStorage := sut.storage.(*storage.MemoryStorage)

	// Arrange
	file, err := createDefaultDirTempFile("test_publish_file_send_file_to_storage.json")
	assert.Nil(t, err)

	fileKey := file.Name

	// Action
	_, err = sut.PublishFile(context.TODO(), fileKey, file)
	assert.Nil(t, err)

	// Arrange
	assert.True(t, memoryStorage.FileExists(fileKey))
}

func TestPublishFileSendEventToStreamer(t *testing.T) {
	// Prepare
	sut, memoryBroker := newSutWithBroker()

	// Arrange
	file, err := createDefaultDirTempFile("test_publish_file_send_file_to_storage.json")
	assert.Nil(t, err)

	fileKey := file.Name

	expectedEvent := models.Event{
		Topic: "collector.files",
		Key:   "published",
		Data:  map[string]string{"file_key": fileKey},
	}

	// Action
	_, err = sut.PublishFile(context.TODO(), fileKey, file)
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

	collectorService, err := collector.NewLocalCollector(fmt.Sprintf("%s/*", folder))
	assert.Nil(t, err)

	memoryStorage := sut.storage.(*storage.MemoryStorage)

	assert.Empty(t, memoryStorage.GetAllFiles())

	createTempFile(folder, "test_file_1.json")
	createTempFile(folder, "test_file_2.json")

	// Arrange
	go sut.Consume(collectorService)
	time.Sleep(time.Second * 1)

	// Assert
	storedFiles := memoryStorage.GetAllFiles()
	assert.Len(t, storedFiles, 2)
}
