package joincitectdbf

import (
	"fmt"
	"github.com/Valentin-Kaiser/go-dbase/dbase"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// FindAndJoinDbfFiles loops through the master DBF table and reads the NAME column.
// For each folder in the NAME column, it looks for a variable.dbf file in the folder.
// If the variable.dbf file is found, it appends the data to the data slice.
// Once all the folders have been processed, it returns the data slice.

func FindAndJoinDbfFiles(masterDbfPath string, masterDbfTable *dbase.File, log *logrus.Logger, dbfFilename string) [][]string {

	rootDir := filepath.Dir(masterDbfPath)
	defer masterDbfTable.Close()
	var data [][]string
	headerWritten := false

	for !masterDbfTable.EOF() {
		row, err := masterDbfTable.Next()
		if err != nil {
			log.WithField("File", masterDbfPath).WithError(err).Warning("Error reading row in DBF")
			continue
		}

		folderNameField := row.FieldByName("NAME")
		if folderNameField == nil {
			log.WithField("File", masterDbfPath).WithError(err).Fatal("Unable to find NAME column")
			os.Exit(1)
		}
		folderName := folderNameField.GetValue().(string)
		variableDbfPath := filepath.Join(rootDir, folderName, dbfFilename)
		if _, err = os.Stat(variableDbfPath); os.IsNotExist(err) {
			log.WithField("File", folderName).WithError(err).Info("variable.dbf not found in folder")
			continue
		}

		variableTable, err := dbase.OpenTable(&dbase.Config{
			Filename:   variableDbfPath,
			TrimSpaces: true,
			Untested:   true, // Citect uses DBF Version 3, which is untested, we are just reading which is backward compatible with version 3.
			ReadOnly:   true,
		})
		log.WithField("Path", variableDbfPath).WithField("Row Count", variableTable.RowsCount()).Info("Opened File")

		if err != nil {
			log.WithField("File", variableDbfPath).WithError(err).Warning("Unable to open file")
			continue
		}

		if !headerWritten {
			var headers []string
			for _, column := range variableTable.Columns() {
				headers = append(headers, column.Name())
			}
			data = append(data, headers)
			headerWritten = true
		}

		for !variableTable.EOF() {
			variableRow, err := variableTable.Next()
			if err != nil {
				log.WithField("File", variableDbfPath).WithError(err).Error("Error reading dbf row")
				break
			}

			var record []string
			for _, field := range variableRow.Fields() {
				record = append(record, fmt.Sprintf("%v", field.GetValue()))
			}
			data = append(data, record)
		}
		variableTable.Close()
	}
	masterDbfTable.Close()

	return data
}
