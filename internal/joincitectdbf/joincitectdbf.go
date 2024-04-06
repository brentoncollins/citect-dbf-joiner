package joincitectdbf

import (
	"github.com/Valentin-Kaiser/go-dbase/dbase"
	"github.com/sirupsen/logrus"
	"os"
)

// GetMadterDbfTable calls masterDbfPath to search for master.dbf in the given directory.
// If not found, it checks the latest subdirectory based on modification time.
// If the master.dbf file is found, it opens the file and returns the Master DBF Table and the path to the master.dbf file.
// If the master.dbf file is not found, it logs an error and exits the application.
func GetMadterDbfTable(inputDir string, log *logrus.Logger) (*dbase.File, string) {

	// Attempt to find the master.dbf file
	masterDbfPath, err := findMasterDBF(inputDir, log)
	if err != nil {
		log.WithError(err).Error("Exiting")
		os.Exit(1)
	} else {
		log.WithField("File", masterDbfPath).Info("master.dbf found, processing...")
	}

	table, err := dbase.OpenTable(&dbase.Config{
		Filename:   masterDbfPath,
		TrimSpaces: true,
		Untested:   true,
		ReadOnly:   true,
	})
	if err != nil {
		panic(err)
	}

	return table, masterDbfPath

}
