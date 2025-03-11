package tests

import (
	_ "bytes"
	"code/internals"
	"errors"
	"io"
	"os"
	"testing"
)

// MockFileReader is a mock implementation of FileReader for testing
type MockFileReader struct {
	Content string
	Err     error
}

func (m MockFileReader) Open(name string) (*os.File, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	// Create a temporary file with mock content
	tmpFile, err := os.CreateTemp("", "testfile-*.csv")
	if err != nil {
		return nil, err
	}
	_, _ = tmpFile.WriteString(m.Content)
	_, _ = tmpFile.Seek(0, io.SeekStart) // Reset file pointer
	return tmpFile, nil
}

func TestReadCSVContent(t *testing.T) {
	tests := []struct {
		name        string
		mockReader  MockFileReader
		expectedErr bool
		expectedURL []string
	}{
		{
			name: "Valid CSV with URLs",
			mockReader: MockFileReader{
				Content: "url\nhttp://example.com\nhttps://test.com",
			},
			expectedErr: false,
			expectedURL: []string{"http://example.com", "https://test.com"},
		},
		{
			name: "File open error",
			mockReader: MockFileReader{
				Err: errors.New("failed to open file"),
			},
			expectedErr: true,
		},
		{
			name: "Empty CSV file",
			mockReader: MockFileReader{
				Content: "",
			},
			expectedErr: true,
		},
		//{
		//	name: "CSV with empty records",
		//	mockReader: MockFileReader{
		//		Content: "url\n\nhttp://example.com\n",
		//	},
		//	expectedErr: false,
		//	expectedURL: []string{"http://example.com"},
		//},
		//{
		//	name: "CSV with invalid records",
		//	mockReader: MockFileReader{
		//		Content: "url\ninvalid-url\nhttp://example.com\n",
		//	},
		//	expectedErr: false,
		//	expectedURL: []string{"invalid-url", "http://example.com"},
		//},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			urlChan := make(chan string, 10)
			err := internals.ReadCSVContent("test.csv", urlChan, tc.mockReader)

			if tc.expectedErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return // Stop test if an error was expected
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			var urls []string
			for url := range urlChan {
				urls = append(urls, url)
			}

			if len(urls) != len(tc.expectedURL) {
				t.Errorf("Expected %d URLs, got %d", len(tc.expectedURL), len(urls))
			}

			for i, expectedURL := range tc.expectedURL {
				if urls[i] != expectedURL {
					t.Errorf("Expected URL %s, got %s", expectedURL, urls[i])
				}
			}
		})
	}
}
