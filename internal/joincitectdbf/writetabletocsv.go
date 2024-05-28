package joincitectdbf

import (
	"encoding/csv"
	"github.com/sirupsen/logrus"
	"os"
)

// WriteToCSV writes the data to a CSV file
func WriteToCSV(data [][]string, outputFile string, log *logrus.Logger) {

	csvFile, err := os.Create(outputFile)
	log.WithField("Output File", outputFile).Info("Creating CSV file")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvWriter := csv.NewWriter(csvFile)
	for _, row := range data {
		err = csvWriter.Write(row) // Changed from := to =
		if err != nil {
			log.Fatalf("failed writing row to CSV: %s", err)
		}
	}

	csvWriter.Flush()

	err = csvWriter.Error()
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Successfully written to file: %s", outputFile)
}
