package main

import (
	"CitectDBFJoiner/internal/joincitectdbf"
	"CitectDBFJoiner/internal/log"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-inputdir string\n\tTHe input directory that either contains the master.dbf or multiple compile folders (required)\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-dbf string\n\tComma seperated list of DBF's you want to join, each one will have a seperate output. (default \"variable,digalm,equip\")\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-outputdir string\n\tThe output CSV file path (default \"Current working directory\")\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "-logpath string\n\tThe output log path (default \"Current working directory\")\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Example:\n")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\t%s -inputdir \"%s\" -outputdir \"%s\" -logpath \"%s\" -dbf \"%s\"\n", os.Args[0], "C:\\ProgramData\\AVEVA Plant SCADA 2023\\User", "C:\\Temp\\", "C:\\Temp\\application.log", "variable,digalm,equip")
}

func main() {

	// Get the command line arguments
	inputDir := flag.String("inputdir", "", "THe input directory that either contains the master.dbf or multiple compile folders")
	outputDir := flag.String("outputdir", "", "Output CSV file path")
	logFile := flag.String("logpath", "application.log", "Output log file")
	dbfTypes := flag.String("dbf", "variable,digalm,equip", "List of DBF types")
	flag.Parse()

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

		dbfFilename := fmt.Sprintf("%s.dbf", dbfType)
		csvFilename := fmt.Sprintf("%s.csv", dbfType)

		logger.WithField("Joining DBF Type", dbfFilename)

		// Get the table and master DBF path
		table, masterDbfPath := joincitectdbf.GetMadterDbfTable(*inputDir, logger)

		// Set the output directory for the CSV
		csvFullPath := filepath.Join(outputDirectory, csvFilename)

		// Get the variable table
		variableTable := joincitectdbf.FindAndJoinDbfFiles(masterDbfPath, table, logger, dbfFilename)
		joincitectdbf.WriteToCSV(variableTable, csvFullPath, logger)

	}

}
