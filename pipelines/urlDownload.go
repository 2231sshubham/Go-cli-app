package pipelines

import (
	"context"
	"log"
	"sync"
	"time"

	"code/internals"
)

const (
	maxWorkerCount = 50
)

func RunURLDownloadPipeline(file string, downloadDir string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	urlChan := make(chan string, 100)
	resultChan := make(chan []byte, 50)
	var wg sync.WaitGroup

	fileReader := internals.OSFileReader{}
	httpClient := internals.DefaultHTTPClient{}
	fileStorage := internals.OSFileStorage{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := internals.ReadCSVContent(file, urlChan, fileReader); err != nil {
			log.Fatalf("Error reading CSV: %v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		internals.DownloadURLWorker(urlChan, resultChan, maxWorkerCount, httpClient)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		internals.PersistFileWorker(downloadDir, resultChan, fileStorage)
	}()

	go func() {
		<-ctx.Done()
		log.Println("Pipeline timeout reached, stopping...")
	}()

	wg.Wait()
	log.Println("Pipeline finished.")
}
