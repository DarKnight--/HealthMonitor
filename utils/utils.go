package utils

import (
	"fmt"
	"log"
	"os"
)

type Status struct {
	Module int
	Run    bool
}

// Perror is used to print the errors.
func Perror(out string) {
	log.Println("[!] Error: " + out)
}

// PLogError is used to print the error when log file is inaccessable
func PLogError(err error) {
	if err != nil {
		fmt.Println("Unable to open log file, path to the log file not found.")
		fmt.Println("OWTF Health Monitor will now exit. Run the setup script to" +
			"set up the log and configuration filess")
		os.Exit(1)
	}
}

// PDBFileError is used to print the error when config is not correct
func PDBFileError(err error) {
	if err != nil {
		fmt.Println("File error: %v\n", err)
		fmt.Println("Configuration file is corrupt. Please run the setup script" +
			"to correct the error. OWTF Health Monitor will now exit.")
		os.Exit(1)
	}
}

// PFileError is used to print the error when monitor do not have sufficient
// file permission
func PFileError(fileName string) {
	fmt.Println("Unable to modify or create %s", fileName)
	fmt.Println("Please check the permission associated with the file")
}
