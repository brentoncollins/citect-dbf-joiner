package joincitectdbf

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// FindMasterDBF searches for master.dbf in the provided directory or subdirectories.
func FindMasterDBF(dir string, searchMasterDBF bool, log *logrus.Logger) (string, error) {
	masterPath := filepath.Join(dir, "master.dbf")
	if _, err := os.Stat(masterPath); err == nil {
		log.Infof("master.dbf found in %s", dir)
		return masterPath, nil
	} else if !searchMasterDBF {
		log.Infof("master.dbf not found in %s", dir)
		return "", fmt.Errorf("master.dbf not found in %s and searching subdirectory for master.dbf was not set", dir)
	}

	// Read subdirectories, sort by modification time in descending order.
	subDirs, err := sortedSubDirs(dir)
	log.Infof("Searching for master.dbf in %s", dir)
	if err != nil {
		return "", err
	}

	// Iterate through sorted subdirectories to find master.dbf.
	for _, subDir := range subDirs {
		log.WithField("Directory", subDir.Name()).Info("Checking subdirectory")
		subDirPath := filepath.Join(dir, subDir.Name())
		masterPath = filepath.Join(subDirPath, "master.dbf")
		if _, err = os.Stat(masterPath); err == nil {
			log.Infof("master.dbf found in %s", subDirPath)
			return masterPath, nil
		} else {
			log.Infof("master.dbf not found in %s", subDirPath)
			log.WithField("Directory", masterPath).Warning("master.dbf not found")
		}
	}

	return "", fmt.Errorf("master.dbf not found in %s or its immediate subdirectories", dir)
}
