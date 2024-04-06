package joincitectdbf

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// findMasterDBF searches for a master.dbf file in the given directory.
// If not found, it checks the latest subdirectories in order ordered on modification time.
// Reason for searching for the master.dbf file the current folder first, then subdirectories:
//The user can pass the input directory that contains the master.dbf file in the root directory.
//.
//└── User
//├── folder1
//│   └── variable.dbf
//├── folder2
//│   └── variable.dbf
//├── folder3
//│   └── variable.dbf
//└── master.dbf
//
//Or user can pass the folder with the latest compiles, if the utility does not find the master.dbf file in the root directory mentioned above.
//It will attempt to find it in the latest subdirectory based on modification time.
//It will only join in the DBF's in the subfolders of the first found master.dbf file.
//.
//└── Compiles
//│
//├── Compile1------------------------------Date: 01/01/2021
//│   ├── folder1
//│   │   └── variable.dbf
//│   ├── folder2
//│   │   └── variable.dbf
//│   ├── folder3
//│   │   └── variable.dbf
//│   └── master.dbf
//│
//└── Compile2------------------------------Date: 02/01/2021
//├── folder1
//│ 	└── variable.dbf
//├── folder2
//│ 	└── variable.dbf
//├── folder3
//│ 	└── variable.dbf
//└── master.dbf

func findMasterDBF(dir string, log *logrus.Logger) (string, error) {
	masterPath := filepath.Join(dir, "master.dbf")
	if _, err := os.Stat(masterPath); err == nil {
		// MASTER.DBF found in the current directory.
		return masterPath, nil
	}

	// Read subdirectories, sort by modification time in descending order.
	subDirs, err := sortedSubdirs(dir)
	if err != nil {
		return "", err
	}

	// Iterate through sorted subdirectories to find master.dbf.
	for _, subDir := range subDirs {
		log.WithField("Directory", subDir.Name()).Info("Checking subdirectory")
		subDirPath := filepath.Join(dir, subDir.Name())
		masterPath = filepath.Join(subDirPath, "master.dbf")
		if _, err = os.Stat(masterPath); err == nil {
			// master.dbf found in a subdirectory.
			return masterPath, nil
		} else {
			log.WithField("Directory", masterPath).Warning("master.dbf not found")
		}
	}

	return "", fmt.Errorf("master.dbf not found in %s or its immediate subdirectories", dir)
}
