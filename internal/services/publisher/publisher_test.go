package publisher

import (
	"context"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
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

	server, err := fileserver.NewLocalFileServer(fileserver.Config{})
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(dir, fileName)

	err = ioutil.WriteFile(fp, []byte{}, fs.ModePerm)
	if err != nil {
		panic(err)
	}

	return models.NewFile(fileName, fp, fileName, server)
}

func newSut() *Publisher {
	eventChannel := make(chan models.Event, 10)
	waitGroup := &sync.WaitGroup{}

	return New(1, "files", storage.NewMemoryStorage(), eventChannel, waitGroup)
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

func TestProcessFileSendFileToStorage(t *testing.T) {
	// Prepare
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	sut := newSut()

	// Arrange
	memoryStorage, _ := sut.storage.(*storage.MemoryStorage)
	assert.Empty(t, memoryStorage.GetAllFiles())

	testFile1, err := createTempFile(folder, "test_file_1.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile1.FilePath)

	// Action
	sut.waitGroup.Add(1)
	sut.processFile(context.TODO(), testFile1)

	// Assert
	storedFiles := memoryStorage.GetAllFiles()
	assert.Len(t, storedFiles, 1)
}

func TestProcessFileShouldReturnErrorWhenFileNoExists(t *testing.T) {
	// Prepare
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	sut := newSut()

	// Arrange
	inexistingFile, err := createTempFile(folder, "test_inexist_file.json")
	assert.Nil(t, err)
	assert.FileExists(t, inexistingFile.FilePath)
	err = os.Remove(inexistingFile.FilePath)
	assert.NoFileExists(t, inexistingFile.FilePath)

	// Action
	sut.waitGroup.Add(1)
	err = sut.processFile(context.TODO(), inexistingFile)

	// Assert
	sut.waitGroup.Wait()
	assert.NotNil(t, err)
}

func TestHandleConsumeFiles(t *testing.T) {
	// Prepare
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	sut := newSut()
	fileChannel := make(chan models.File, 2)

	// Arrange
	memoryStorage, _ := sut.storage.(*storage.MemoryStorage)
	assert.Empty(t, memoryStorage.GetAllFiles())

	testFile1, err := createTempFile(folder, "test_file_1.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile1.FilePath)

	testFile2, err := createTempFile(folder, "test_file_2.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile2.FilePath)

	// Action
	fileChannel <- testFile1
	fileChannel <- testFile2
	sut.waitGroup.Add(2)

	sut.Handle(context.TODO(), fileChannel)
	sut.waitGroup.Wait()

	// Assert
	storedFiles := memoryStorage.GetAllFiles()
	assert.Len(t, storedFiles, 2)
}

func TestHandleShouldSendSuccessEventWhenProcessFileReturnNil(t *testing.T) {
	// Prepare
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	sut := newSut()
	fileChannel := make(chan models.File, 2)

	// Arrange
	successFile, err := createTempFile(folder, "test_success_file.json")
	assert.Nil(t, err)
	assert.FileExists(t, successFile.FilePath)

	errorFile, err := createTempFile(folder, "test_error_file.json")
	assert.Nil(t, err)
	assert.FileExists(t, successFile.FilePath)

	err = os.Remove(errorFile.FilePath)
	assert.Nil(t, err)
	assert.NoFileExists(t, errorFile.FilePath)

	// Action
	fileChannel <- successFile
	sut.waitGroup.Add(1)

	sut.Handle(context.Background(), fileChannel)
	sut.waitGroup.Wait()

	// Assert
	expectedEvent, err := models.NewEvent(sut.EventTopic, "success", map[string]string{"file_key": successFile.Key})
	assert.Nil(t, err)

	event := <-sut.eventChannel

	assert.Equal(t, expectedEvent, event)
}

func TestHandleShouldSendErrorEventWhenProcessFileReturnError(t *testing.T) {
	// Prepare
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	sut := newSut()
	fileChannel := make(chan models.File, 2)

	// Arrange
	errorFile, err := createTempFile(folder, "test_error_file.json")
	assert.Nil(t, err)
	assert.FileExists(t, errorFile.FilePath)

	err = os.Remove(errorFile.FilePath)
	assert.Nil(t, err)
	assert.NoFileExists(t, errorFile.FilePath)

	// Action
	fileChannel <- errorFile
	sut.waitGroup.Add(1)

	sut.Handle(context.TODO(), fileChannel)
	sut.waitGroup.Wait()

	// Assert
	event := <-sut.eventChannel
	eventData := event.Data.(map[string]string)
	eventError, _ := eventData["error"]
	eventFilePath, _ := eventData["file_path"]

	assert.Equal(t, "error", event.Key)
	assert.Equal(t, sut.EventTopic, event.Topic)
	assert.Equal(t, errorFile.FilePath, eventFilePath)
	assert.Contains(t, eventError, "no such file or directory")
}
