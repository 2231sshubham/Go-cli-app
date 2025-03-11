package main

import (
	"code/pipelines"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: go run main.go <csv filepath>")
	}
	filePath := os.Args[1]
	downloadPath := os.Args[2]
	log.Println("Starting pipeline...")
	pipelines.RunURLDownloadPipeline(filePath, downloadPath)
}
