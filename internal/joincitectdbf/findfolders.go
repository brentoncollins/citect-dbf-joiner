package joincitectdbf

import (
	"fmt"
	"github.com/Valentin-Kaiser/go-dbase/dbase"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// GetMasterDBFTableFoldersAsSlice reads the master.dbf file and returns a slice of folder paths.
func GetMasterDBFTableFoldersAsSlice(masterDbfPath string, log *logrus.Logger) ([]string, error) {
	log.WithField("File", masterDbfPath).Info("Opening file")
	masterDbfTable, err := dbase.OpenTable(&dbase.Config{
		Filename:   masterDbfPath,
		TrimSpaces: true,
		Untested:   true,
		ReadOnly:   true,
	})
	if err != nil {
		log.WithField("File", masterDbfPath).WithError(err).Fatal("Unable to open file")
		os.Exit(1)
	}
	defer masterDbfTable.Close()

	var folders []string
	rootDir := filepath.Dir(masterDbfPath)

	for !masterDbfTable.EOF() {
		row, err := masterDbfTable.Next()

		if err != nil {

			continue
		}

		folderNameField := row.FieldByName("NAME")
		if folderNameField == nil {

			os.Exit(1)
		}
		folderName := folderNameField.GetValue().(string)
		folderPath := filepath.Join(rootDir, folderName)

		if _, err = os.Stat(folderPath); !os.IsNotExist(err) {
			log.WithField("Folder", folderPath).Info("Found folder")
			folders = append(folders, folderPath)
		}
	}

	return folders, err
}

// GetFolders returns a slice of folder paths in the provided directory.
func GetFolders(dir string, fromCurrent bool) ([]string, error) {
	subDirs, err := sortedSubDirs(dir)

	if err != nil {
		return nil, err
	}

	if !fromCurrent && len(subDirs) > 0 {
		// Get the latest subdirectory
		latestSubDir := subDirs[0]
		dir = filepath.Join(dir, latestSubDir.Name())
		subDirs, err = sortedSubDirs(dir)
		if err != nil {
			return nil, err
		}
	}

	folders := make([]string, len(subDirs))
	for i, subDir := range subDirs {

		folders[i] = filepath.Join(dir, subDir.Name())
	}

	return folders, nil
}

// FindAndJoinDbfFiles reads the variable.dbf files in the provided folders and returns the data as a slice of slices.
func FindAndJoinDbfFiles(folders []string, log *logrus.Logger, dbfFilename string) [][]string {
	var data [][]string
	headerWritten := false

	for _, folderPath := range folders {
		variableDbfPath := filepath.Join(folderPath, dbfFilename)
		if _, err := os.Stat(variableDbfPath); os.IsNotExist(err) {
			log.WithField("File", folderPath).WithError(err).Info("variable.dbf not found in folder")
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

	return data
}
