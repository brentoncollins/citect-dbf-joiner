package joincitectdbf

import (
	"io/fs"
	"os"
	"sort"
)

// sortedSubdirs returns a list of subdirectories sorted by modification time (newest first).
func sortedSubdirs(dir string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	subDirs := make([]fs.DirEntry, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			subDirs = append(subDirs, entry)
		}
	}

	// Sort subDirs by modification time in descending order.
	sort.Slice(subDirs, func(i, j int) bool {
		infoI, _ := subDirs[i].Info()
		infoJ, _ := subDirs[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	return subDirs, nil
}
