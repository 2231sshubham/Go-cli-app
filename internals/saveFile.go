package internals

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileStorage interface {
	Save(filePath string, data []byte) error
}

type OSFileStorage struct{}

func (OSFileStorage) Save(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}
	return nil
}

func PersistFileWorker(path string, resultChan <-chan []byte, storage FileStorage) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	for data := range resultChan {
		filename := filepath.Join(path, fmt.Sprintf("download_%d.txt", time.Now().UnixNano()))
		if err := storage.Save(filename, data); err != nil {
			log.Printf("Error saving file: %v", err)
		}
	}
}
