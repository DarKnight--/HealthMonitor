package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

var (
	mutex sync.Mutex
	//ControlChan is the channel to send stop or start signal to main function
	ControlChan chan Status
	//Modules is the list of the modules currently implemented..
	Modules = []string{"live", "target", "disk", "ram", "cpu"}
	//LiveEmergency is the channel to call live module any time
	LiveEmergency chan bool
)

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

//OpenLogFile is the utility function to open log file
func OpenLogFile(logFileName string) *os.File {
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		PLogError(err)
	}
	return logFile
}

//SendModuleStatus is the utility function to change the state of module
func SendModuleStatus(module string, status bool) {
	signal := Status{Module: module, Run: status}
	ControlChan <- signal
}

//CheckConf is the utility function to check the config variable loaded from the
//database and if fails, then switch to default
func CheckConf(moduleLogFile *os.File, masterLogFile *os.File, module string,
	profile *string, setupFunc func()) {
	ModuleError(moduleLogFile, "Unable to find config for profile "+
		*profile, "Setting up environment")
	if *profile == "default" {
		setupFunc()
	} else {
		*profile = "default"
		ModuleError(masterLogFile, fmt.Sprintf("Unable to load profile: %s for %s module",
			*profile, module), "Restating monitor with default value")
		RestartAllModules()
	}
}

//RestartAllModules will restart all the modules.
func RestartAllModules() {
	var module string
	for _, module = range Modules {
		SendModuleStatus(module, false)
		SendModuleStatus(module, true)
	}
}
