package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

var mutex sync.Mutex

// Status struct is used by monitor to send different modules signal to abort
type Status struct {
	Module string
	Run    bool
}

// Perror is used to print the errors.
func Perror(out string) {
	log.Println("[!] Error: " + out)
}

// PLogError is used to print the error when log file is inaccessable
func PLogError(err error) {
	if err != nil {
		fmt.Println("Error in opening log file")
		fmt.Println(err.Error())
		fmt.Println("OWTF Health Monitor will now exit. Run the setup script to" +
			"set up the log and configuration files")
		os.Exit(1)
	}
}

// PFileError is used to print the error when monitor do not have sufficient
// file permission
func PFileError(fileName string) {
	log.Println(fmt.Sprintf("Unable to modify or create %s", fileName))
	log.Println("Please check the permission associated with the file")
}

// GetPath returns the absolute path
func GetPath(paths string) string {
	if strings.HasPrefix(paths, "/") {
		return paths
	}
	return path.Join(os.Getenv("HOME"), paths)
}

// ModuleLogs is used to write the logs of the module in the @filename file
func ModuleLogs(filename *os.File, status string) {
	mutex.Lock()
	log.SetOutput(filename)
	log.Println(status)
	mutex.Unlock()
}

// ModuleError is used to log the errors of the module in the @filename file
func ModuleError(filename *os.File, err string, description string) {
	mutex.Lock()
	log.SetOutput(filename)
	log.Println("[!] Error occured : " + err)
	log.Print(description)
	mutex.Unlock()
}

func OpenLogFile(logFileName string) *os.File {
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		PLogError(err)
	}
	return logFile
}
