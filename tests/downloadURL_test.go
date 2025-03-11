package tests

import (
	"code/internals"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// MockHTTPClient is a mock implementation of HTTPClient
type MockHTTPClient struct {
	ResponseBody string
	ResponseCode int
	Err          error
}

// Get simulates an HTTP request and returns a mock response
func (m MockHTTPClient) Get(url string) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &http.Response{
		StatusCode: m.ResponseCode,
		Body:       io.NopCloser(strings.NewReader(m.ResponseBody)),
	}, nil
}

func TestDownloadURL(t *testing.T) {
	tests := []struct {
		name         string
		mockClient   MockHTTPClient
		expectedErr  bool
		expectedData string
	}{
		{
			name: "Valid URL returns content",
			mockClient: MockHTTPClient{
				ResponseBody: "Hello, world!",
				ResponseCode: http.StatusOK,
			},
			expectedErr:  false,
			expectedData: "Hello, world!",
		},
		{
			name: "HTTP error from client",
			mockClient: MockHTTPClient{
				Err: errors.New("failed to connect"),
			},
			expectedErr: true,
		},
		{
			name: "Non-200 status code",
			mockClient: MockHTTPClient{
				ResponseBody: "Not found",
				ResponseCode: http.StatusNotFound,
			},
			expectedErr: true,
		},
		{
			name: "Empty content response",
			mockClient: MockHTTPClient{
				ResponseBody: "",
				ResponseCode: http.StatusOK,
			},
			expectedErr:  false, // It will be handled in worker
			expectedData: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := internals.DownloadURL("http://test.com", tc.mockClient)

			if tc.expectedErr {
				if result.Error == nil {
					t.Errorf("Expected error but got none")
				}
				return // Stop test if an error was expected
			}

			if result.Error != nil {
				t.Errorf("Unexpected error: %v", result.Error)
			}

			if string(result.Content) != tc.expectedData {
				t.Errorf("Expected content '%s', got '%s'", tc.expectedData, string(result.Content))
			}
		})
	}
}

func TestDownloadURLWorker(t *testing.T) {
	urlChan := make(chan string, 2)
	resultChan := make(chan []byte, 2)

	mockClient := MockHTTPClient{
		ResponseBody: "mock response",
		ResponseCode: http.StatusOK,
	}

	// Sending URLs to channel
	urlChan <- "http://example.com"
	urlChan <- "http://test.com"
	close(urlChan)

	go internals.DownloadURLWorker(urlChan, resultChan, 2, mockClient)

	var results []string
	timeout := time.After(2 * time.Second)
	done := make(chan bool)

	go func() {
		for res := range resultChan {
			results = append(results, string(res))
		}
		done <- true
	}()

	select {
	case <-done:
		expectedCount := 2
		if len(results) != expectedCount {
			t.Errorf("Expected %d responses, got %d", expectedCount, len(results))
		}
		for _, res := range results {
			if res != "mock response" {
				t.Errorf("Unexpected content: %s", res)
			}
		}
	case <-timeout:
		t.Fatal("Test timed out")
	}
}

func TestDownloadURLWorker_InvalidResponse(t *testing.T) {
	urlChan := make(chan string, 1) // Only one URL needed
	resultChan := make(chan []byte, 1)
	done := make(chan bool)

	mockClient := MockHTTPClient{
		ResponseBody: "",
		ResponseCode: http.StatusInternalServerError,
	}

	urlChan <- "http://badresponse.com"
	close(urlChan)

	go func() {
		internals.DownloadURLWorker(urlChan, resultChan, 1, mockClient)
		done <- true
	}()

	select {
	case <-done:
		select {
		case res := <-resultChan:
			t.Errorf("Expected no response, but got: %s", string(res))
		default:
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}
}
