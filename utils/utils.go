// Package contains various utility function used in the HealthMonitor
package utils

import (
	"fmt"
	"log"
	"os"
)

/* It is used to print the errors.*/
func Perror(out string) {
	log.Println("[!] Error: " + out)
}

func PLogError(err error){
	if err != nil{
		fmt.Println("Unable to open log file, path to the log file not found.")
		fmt.Println(`OWTF Health Monitor will now exit. Run the setup script to
					set up the log and configuration filess`)
		os.Exit(1)
	}
}

func PDBFileError(err error){
	if err != nil{
		fmt.Println("File error: %v\n", err)
		fmt.Println("Configuration file is corrupt. Please run the setup script"+
					"to correct the error. OWTF Health Monitor will now exit.")
		os.Exit(1)
	}
}