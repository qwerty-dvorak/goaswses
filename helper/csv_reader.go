package helper

import (
	"encoding/csv"
	"os"
)

func ReadCsvFile(filePath string) ([]string, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	var emails []string
	for _, record := range records {
		emails = append(emails, record[0])
	}
	if err != nil {
		return nil, err
	}
	return emails, nil
}