package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"
)

func CloseFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatalf("Error closing file, %v", err)
	}
}

func GenerateFilename() string {
	// Generate 8 random bytes
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		panic(fmt.Sprintf("Failed to generate random bytes: %v", err))
	}

	// Convert bytes to hex string
	randomStr := hex.EncodeToString(bytes)

	// Append timestamp to ensure uniqueness
	timestamp := time.Now().UnixNano()

	// Construct the filename
	return fmt.Sprintf("%d_%s%s", timestamp, randomStr)
}
