package internals

import (
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type DefaultHTTPClient struct{}

func (DefaultHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

type DownloadedContent struct {
	Content []byte
	Error   error
}

func DownloadURL(url string, client HTTPClient) DownloadedContent {
	resp, err := client.Get(url)
	if err != nil {
		return DownloadedContent{Error: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DownloadedContent{Error: errors.New("unexpected status code: " + resp.Status)}
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return DownloadedContent{Error: err}
	}
	return DownloadedContent{Content: content}
}

func DownloadURLWorker(urlChan <-chan string, resultChan chan<- []byte, maxWorkers int, client HTTPClient) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)

	for url := range urlChan {
		wg.Add(1)
		go func(url string) {
			sem <- struct{}{}
			defer wg.Done()
			defer func() { <-sem }()

			content := DownloadURL(url, client)
			if content.Error != nil {
				log.Printf("Error downloading %s: %v", url, content.Error)
				return
			}
			if len(content.Content) < 1 {
				log.Printf("Empty content fount at %s returning.", url)
				return
			}
			resultChan <- content.Content
		}(url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()
}
