package main

import (
	"flag"
	"fmt"
	"github.com/brentoncollins/citect-dbf-joiner/internal/joincitectdbf"
	"github.com/brentoncollins/citect-dbf-joiner/internal/log"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-input-dir string\n\tThe input directory that either contains the master.dbf or multiple compile folders (required)\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-dbf string\n\tComma separated list of DBF's you want to join, each one will have a separate output. (default \"variable,digalm,equip\")\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-output-dir string\n\tThe output CSV file path (default \"Current working directory\")\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-log-path string\n\tThe output log path (default \"Current working directory\")\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-master-dbf-dir\n\tSet this flag to use the master DBF file in the input-dir to determine the folders that contain the DBF files.\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-search-compile-folders-for-master-dbf\n\tSet this flag to find all sub-folders in input-dir, then search each sub-folder by date descending for a master DBF to determine the folders that contain the DBF files within that folder\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-project-compile-dir\n\tSet this flag to find all sub-folders in input-dir, and get DBF's from each folder.\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-search-latest-compile-dirs\n\tSet this flag to find all sub-folders in input-dir, identify the latest sub-folder by date descending, each sub-folder within the latest sub-folder will be searched for a DBF.\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Example:\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\t%s -inputdir \"%s\" -outputdir \"%s\" -logpath \"%s\" -dbf \"%s\"\n", os.Args[0], "C:\\ProgramData\\AVEVA Plant SCADA 2023\\User", "C:\\Temp\\", "C:\\Temp\\application.log", "variable,digalm,equip")
}

func main() {

	// Get the command line arguments
	inputDir := flag.String("input-dir", "", "The input directory that either contains the master.dbf or multiple compile folders")
	outputDir := flag.String("output-dir", "", "Output CSV file path")
	logFile := flag.String("log-path", "citect-dbf-joiner.log", "Output log file")
	dbfTypes := flag.String("dbf", "variable,digalm,equip", "List of DBF types")
	masterDbfDir := flag.Bool("master-dbf-dir", false, "Set this flag to use the master DBF file in the input-dir to determine the folders that contain the DBF files.")
	SearchCompileFolderForDbf := flag.Bool("search-compile-folders-for-master-dbf", false, "Set this flag to find all sub-folders in input-dir, then search each sub-folder by date descending for a master DBF to determine the folders that contain the DBF files within that folder")
	projectCompileDir := flag.Bool("project-compile-dir", false, "Set this flag to find all sub-folders in input-dir, and get DBF's from each folder.")
	searchLatestCompileDirs := flag.Bool("search-latest-compile-dirs", false, "Set this flag to find all sub-folders in input-dir, identify the latest sub-folder by date descending, each sub-folder within the latest sub-folder will be searched for a DBF.")

	flag.Parse()

	// Make sure only one flag is true
	countTrue := 0
	if *masterDbfDir {
		countTrue++
	}
	if *SearchCompileFolderForDbf {
		countTrue++
	}
	if *projectCompileDir {
		countTrue++
	}
	if *searchLatestCompileDirs {
		countTrue++
	}

	// Check if exactly one flag is true
	if countTrue != 1 {
		fmt.Println("Exactly one of the flags -master-dbf-dir, -search-compile-folders-for-master-dbf, -project-compile-dir, -search-latest-compile-dirs must used.")
		os.Exit(1)
	}

	// Get the list of dbf types to join
	dbfList := strings.Split(*dbfTypes, ",")

	// If the input directory is not supplied, exit the application and print the usage to the terminal.
	if *inputDir == "" {
		fmt.Println("The input directory not supplied.")
		usage()
		os.Exit(1)
	}

	// Check if the input directory exists
	if _, err := os.Stat(*inputDir); os.IsNotExist(err) {
		fmt.Println("The input directory does not exist.")
		usage()
		os.Exit(1)
	}

	var outputDirectory string
	var err error

	if *outputDir == "" {
		// If output directory is not provided, use the current working directory
		outputDirectory, err = os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			os.Exit(1)
		}
	} else {
		// If output directory is provided, check if it exists
		if _, err = os.Stat(*outputDir); os.IsNotExist(err) {
			fmt.Println("The output directory does not exist.")
			os.Exit(1)
		}
		// If it exists, set it as the output directory
		outputDirectory = *outputDir
	}

	// Get the logger
	logger := log.Logger(*logFile)

	for _, dbfType := range dbfList {
		logger.Infof("Processing %s.dbf files", dbfType)

		dbfFilename := fmt.Sprintf("%s.dbf", dbfType)
		csvFilename := fmt.Sprintf("%s.csv", dbfType)

		var folders []string
		var masterDbfPath string

		if *masterDbfDir {
			masterDbfPath, err = joincitectdbf.FindMasterDBF(*inputDir, false, logger)
			if err != nil {
				logger.WithError(err).Error("Exiting")
				return
			}
			// Get the folders from the master dbf
			folders, _ = joincitectdbf.GetMasterDBFTableFoldersAsSlice(masterDbfPath, logger)

		} else if *SearchCompileFolderForDbf {
			// Search for the master dbf in the input directory sub-folders
			masterDbfPath, err = joincitectdbf.FindMasterDBF(*inputDir, true, logger)
			if err != nil {
				logger.WithError(err).Error("Exiting")
				return
			}
			// Get the folders from the master dbf
			folders, _ = joincitectdbf.GetMasterDBFTableFoldersAsSlice(masterDbfPath, logger)
		} else if *projectCompileDir {
			// Get the folders from the input directory
			folders, _ = joincitectdbf.GetFolders(*inputDir, true)

		} else if *searchLatestCompileDirs {
			// Ger the folders from the input directory's newest sub-folder
			folders, _ = joincitectdbf.GetFolders(*inputDir, false)
		} else {
			fmt.Println("Exactly one of the flags -master-dbf-dir, -search-compile-folders-for-master-dbf, -project-compile-dir, -search-latest-compile-dirs must used.")
			os.Exit(1)
		}

		// Set the output directory for the CSV
		csvFullPath := filepath.Join(outputDirectory, csvFilename)

		// Get the variable table
		variableTable := joincitectdbf.FindAndJoinDbfFiles(folders, logger, dbfFilename)
		joincitectdbf.WriteToCSV(variableTable, csvFullPath, logger)

	}

}
