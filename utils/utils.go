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
		fmt.Println("Error in opening log file")
		fmt.Println(err.Error())
		fmt.Println("OWTF Health Monitor will now exit. Run the setup script to" +
			"set up the log and configuration filess")
		os.Exit(1)
	}
}

// PFileError is used to print the error when monitor do not have sufficient
// file permission
func PFileError(fileName string) {
	log.Println("Unable to modify or create %s", fileName)
	log.Println("Please check the permission associated with the file")
}

// GetPath returns the absolute path
func GetPath(paths string) string {
	if strings.HasPrefix(paths, "/") {
		return paths
	}
	return path.Join(os.Getenv("HOME"), paths)
}

func ModuleLogs(filename *os.File, status string) {
	mutex.Lock()
	log.SetOutput(filename)
	log.Println(status)
	mutex.Unlock()
}

func ModuleError(filename *os.File, err string, description string) {
	mutex.Lock()
	log.SetOutput(filename)
	log.Println("[!] Error occured : " + err)
	log.Print(description)
	mutex.Unlock()
}
