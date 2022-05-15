package fileserver

import (
	"context"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tmpDir string

func init() {
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	tmpDir = folder
}

func newSut() *LocalFileServer {
	collector, err := NewLocalFileServer(Config{})
	if err != nil {
		panic(err)
	}

	return collector
}

func createTempFile(fileName string) (string, error) {
	fp := filepath.Join(tmpDir, fileName)

	err := ioutil.WriteFile(fp, []byte{}, fs.ModePerm)
	if err != nil {
		return "", err
	}

	return fp, nil
}

func TestGlob(t *testing.T) {
	// Prepare
	sut := newSut()

	// Arrange
	expectedFile, err := createTempFile("test_json_file.json")
	assert.Nil(t, err)

	ignoredFile, err := createTempFile("test_ignore_this_file.txt")
	assert.Nil(t, err)

	// Action
	pattern := filepath.Join(tmpDir, "*.json")
	collectedFiles, err := sut.Glob(context.TODO(), pattern)
	assert.Nil(t, err)

	// Assert
	expectedFileIsCollected := false
	ignoredFileIsCollected := false

	for _, file := range collectedFiles {
		if file == expectedFile {
			expectedFileIsCollected = true
		}
		if file == ignoredFile {
			ignoredFileIsCollected = true
		}
	}

	assert.Truef(t, expectedFileIsCollected, "File '%s' not collected", expectedFile)
	assert.Falsef(t, ignoredFileIsCollected, "File '%s' is collected", ignoredFile)
}

func TestGlobShouldReturnOnlyFiles(t *testing.T) {
	// Prepare
	sut := newSut()

	// Arrange
	expectedFile, err := createTempFile("test_glob_should_return_only_files_collected_file.json")
	assert.Nil(t, err)

	ignoredDir := filepath.Join(tmpDir, "ignored_dir")
	err = os.Mkdir(ignoredDir, os.ModePerm)
	assert.Nil(t, err)

	// Action
	pattern := filepath.Join(tmpDir, "*.json")
	collectedFiles, err := sut.Glob(context.TODO(), pattern)
	assert.Nil(t, err)

	// Assert
	expectedFileIsCollected := false
	dirIsIgnored := true

	for _, file := range collectedFiles {
		if file == expectedFile {
			expectedFileIsCollected = true
		}
		if file == ignoredDir {
			dirIsIgnored = false
		}
	}

	assert.Truef(t, expectedFileIsCollected, "File '%s' not collected", expectedFile)
	assert.Truef(t, dirIsIgnored, "Directory '%s' is not ignored", ignoredDir)
}

func TestGlobShouldReturnOnlyAbsPath(t *testing.T) {
	// Arrange
	sut := newSut()

	// Action
	collectedFiles, err := sut.Glob(context.TODO(), "./*")
	assert.Nil(t, err)

	for _, file := range collectedFiles {
		assert.True(t, filepath.IsAbs(file))
	}
}

func TestRemoveFileDeleteFile(t *testing.T) {
	// Prepare
	sut := newSut()
	fileName := "test_remove_file_delete_from_storage.json"

	// Arrange
	file, err := createTempFile(fileName)
	assert.Nil(t, err)
	assert.FileExists(t, file)

	// Action
	err = sut.Remove(context.TODO(), file)
	assert.Nil(t, err)

	// Assert
	assert.NoFileExists(t, file)
}

func TestOpenReturnValidContentReader(t *testing.T) {
	// Prepare
	sut := newSut()
	fileData := "testing!"
	fileName := "test_file_get_reader_return_valid_content_reader.json"

	// Arrange
	filePath := filepath.Join(tmpDir, fileName)

	err := ioutil.WriteFile(filePath, []byte(fileData), 0o644)
	assert.Nil(t, err)

	// Action
	reader, err := sut.Open(context.TODO(), filePath)
	assert.Nil(t, err)

	// Assert
	data, err := ioutil.ReadAll(reader)
	assert.Nil(t, err)

	assert.Equal(t, []byte(fileData), data)
}

func TestMoveShouldCreateFolders(t *testing.T) {
	// Arrange
	sut := newSut()

	file, err := createTempFile("test_move_file_should_create_folders.txt")
	assert.Nil(t, err)
	assert.FileExists(t, file)

	dir, name := filepath.Split(file)

	// Action
	newpath := filepath.Join(dir, "sent", name)
	err = sut.Move(context.TODO(), file, newpath)
	assert.Nil(t, err)

	// Assert
	assert.NoFileExists(t, file)
	assert.FileExists(t, newpath)
}

func TestAcquireLockShouldSaveLockerToLockedFiles(t *testing.T) {
	// Arrange
	sut := newSut()

	file, err := createTempFile("test_acquire_lock_should_save_locker_on_locked_files.txt")
	assert.Nil(t, err)
	assert.FileExists(t, file)

	// Action
	err = sut.AcquireLock(context.TODO(), file)
	assert.Nil(t, err)

	// Assert
	_, ok := sut.lockedFiles[file]
	assert.True(t, ok)
}

func TestAcquireLockShouldReturnErrorWhenTryLockLockedFile(t *testing.T) {
	// Arrange
	sut := newSut()

	file, err := createTempFile("test_acquire_lock_should_return_error_when_try_lock_locked_file.txt")
	assert.Nil(t, err)
	assert.FileExists(t, file)

	err = sut.AcquireLock(context.TODO(), file)
	assert.Nil(t, err)

	// Action
	err = sut.AcquireLock(context.TODO(), file)

	// Assert
	assert.NotNil(t, err)
}

func TestReleaseLockShouldRemoveLockerFromLockedFiles(t *testing.T) {
	// Arrange
	sut := newSut()

	file, err := createTempFile("test_acquire_lock_should_remove_locker_from_locked_files.txt")
	assert.Nil(t, err)
	assert.FileExists(t, file)

	err = sut.AcquireLock(context.TODO(), file)
	assert.Nil(t, err)

	// Action
	err = sut.ReleaseLock(context.TODO(), file)
	assert.Nil(t, err)

	// Assert
	_, ok := sut.lockedFiles[file]
	assert.False(t, ok)
}

func TestAcquireLockShouldReturnErrorWhenFileNoExists(t *testing.T) {
	// Prepare
	sut := newSut()

	// Arrange
	inexistingFile := "this_file_does_not_exist.txt"
	assert.NoFileExists(t, inexistingFile)

	// Action
	err := sut.AcquireLock(context.TODO(), inexistingFile)

	// Assert
	assert.NotNil(t, err)
}
