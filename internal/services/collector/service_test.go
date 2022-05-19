package collector

import (
	"context"
	"io/fs"
	"io/ioutil"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/fileserver"
)

func createTempFile(dir, fileName string) (models.File, error) {
	if dir == "" {
		tmpDir, err := ioutil.TempDir("", "*")
		if err != nil {
			panic(err)
		}

		dir = tmpDir
	}

	fp := path.Join(dir, fileName)

	err := ioutil.WriteFile(fp, []byte{}, fs.ModePerm)
	if err != nil {
		panic(err)
	}

	server := newFileServer()

	info, err := server.Stat(context.TODO(), fp)
	if err != nil {
		panic(err)
	}

	return models.NewFile(fileName, fp, fileName, info.Size(), info.ModTime(), server)
}

func newFileServer() services.FileServer {
	server, err := fileserver.NewLocalFileServer(fileserver.Config{})
	if err != nil {
		panic(err)
	}
	return server
}

func newSut(patterns ...string) *Collector {
	server := newFileServer()

	cfg := Config{
		MatchPatterns: patterns,
	}

	collector, err := New(1, cfg, server, &sync.WaitGroup{}, &sync.WaitGroup{})
	if err != nil {
		panic(err)
	}

	return collector
}

func TestCollectFilesShouldReturnOneFilesWhenMaxCollectBatchSizeIsOne(t *testing.T) {
	// Prepare
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	pattern := path.Join(folder, "*.json")
	sut := newSut(pattern)

	// Arrange
	testFile1, err := createTempFile(folder, "test_file_1.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile1.FilePath)

	testFile2, err := createTempFile(folder, "test_file_2.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile2.FilePath)

	sut.cfg.MaxCollectBatchSize = 0

	// Action
	fileChannel := make(chan models.File, 2)
	sut.collectGroup.Add(1)
	sut.collectFilesWithPattern(context.TODO(), fileChannel, pattern)

	f1 := <-fileChannel
	f2 := <-fileChannel

	// Assert
	assert.Equal(t, testFile1, f1)
	assert.Equal(t, testFile2, f2)
}

func TestCollectFilesShouldReturnAllFilesWhenMaxCollectBatchSizeIsZero(t *testing.T) {
	// Prepare
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	pattern := path.Join(folder, "*.json")
	sut := newSut(pattern)

	// Arrange
	testFile1, err := createTempFile(folder, "test_file_1.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile1.FilePath)

	testFile2, err := createTempFile(folder, "test_file_2.json")
	assert.Nil(t, err)
	assert.FileExists(t, testFile2.FilePath)

	sut.cfg.MaxCollectBatchSize = 1

	fileChannel := make(chan models.File, 1)

	// Action
	sut.collectGroup.Add(1)
	sut.collectFilesWithPattern(context.TODO(), fileChannel, pattern)
	f1 := <-fileChannel

	// Assert
	assert.Equal(t, testFile1, f1)
}
