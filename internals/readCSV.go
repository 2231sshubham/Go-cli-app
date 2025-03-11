package internals

import (
	"code/utils"
	"encoding/csv"
	"errors"
	"log"
	"os"
)

type FileReader interface {
	Open(name string) (*os.File, error)
}

type OSFileReader struct{}

func (OSFileReader) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func ReadCSVContent(filePath string, urlChan chan<- string, reader FileReader) error {
	defer close(urlChan)
	csvFile, err := reader.Open(filePath)
	if err != nil {
		return err
	}

	defer utils.CloseFile(csvFile)

	csvReader := csv.NewReader(csvFile)

	_, err = csvReader.Read()
	if err != nil {
		return errors.New("invalid csv file")
	}
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("Skipping invalid record due to error: %v", err)
			continue
		}

		// Ensure record is not empty
		if len(record) == 0 || record[0] == "" {
			log.Println("Skipping empty record")
			continue
		}
		urlChan <- record[0]
	}
	return nil
}
