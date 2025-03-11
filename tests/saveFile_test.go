package tests

import (
	"code/internals"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type MockFileStorage struct {
	SaveErr error
}

func (m MockFileStorage) Save(filePath string, data []byte) error {
	if m.SaveErr != nil {
		return m.SaveErr
	}
	return os.WriteFile(filePath, data, 0644)
}

func TestOSFileStorage_Save_Success(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile.txt")
	storage := internals.OSFileStorage{}

	data := []byte("Hello, world!")
	err := storage.Save(filePath, data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify file exists and content is correct
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}
	if string(content) != string(data) {
		t.Errorf("Expected content %q, got %q", data, content)
	}
}

func TestOSFileStorage_Save_Fail_Create(t *testing.T) {
	storage := internals.OSFileStorage{}

	// Trying to save to an invalid path
	err := storage.Save("/invalid_path/test.txt", []byte("data"))
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestOSFileStorage_Save_Fail_Write(t *testing.T) {
	mockStorage := MockFileStorage{SaveErr: errors.New("write error")}

	err := mockStorage.Save("fakefile.txt", []byte("test data"))
	if err == nil || err.Error() != "write error" {
		t.Fatalf("Expected write error, got: %v", err)
	}
}

func TestPersistFileWorker_Success(t *testing.T) {
	tempDir := t.TempDir()
	resultChan := make(chan []byte, 2)
	mockStorage := internals.OSFileStorage{}

	// Sending mock data
	resultChan <- []byte("First file")
	resultChan <- []byte("Second file")
	close(resultChan)

	go internals.PersistFileWorker(tempDir, resultChan, mockStorage)

	// Wait briefly for worker to finish
	time.Sleep(100 * time.Millisecond)

	files, _ := filepath.Glob(filepath.Join(tempDir, "download_*.txt"))
	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}
}

func TestPersistFileWorker_SaveFail(t *testing.T) {
	tempDir := t.TempDir()
	resultChan := make(chan []byte, 1)
	mockStorage := MockFileStorage{SaveErr: errors.New("failed to save")}

	resultChan <- []byte("Test file")
	close(resultChan)

	go internals.PersistFileWorker(tempDir, resultChan, mockStorage)

	time.Sleep(100 * time.Millisecond)

	// No files should be created
	files, _ := filepath.Glob(filepath.Join(tempDir, "download_*.txt"))
	if len(files) != 0 {
		t.Errorf("Expected 0 files due to save failure, but found %d", len(files))
	}
}

func TestPersistFileWorker_EmptyChannel(t *testing.T) {
	tempDir := t.TempDir()
	resultChan := make(chan []byte)
	close(resultChan)

	go internals.PersistFileWorker(tempDir, resultChan, internals.OSFileStorage{})

	time.Sleep(100 * time.Millisecond)

	files, _ := filepath.Glob(filepath.Join(tempDir, "download_*.txt"))
	if len(files) != 0 {
		t.Errorf("Expected 0 files, but found %d", len(files))
	}
}
