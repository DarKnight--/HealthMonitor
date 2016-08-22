package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

var (
	mutex, owtfMutex sync.Mutex

	//ControlChan is the channel to send stop or start signal to main function
	ControlChan chan Status

	//LiveEmergency is the channel to call live module any time
	LiveEmergency chan bool

	//ExitChan is the channel to send signal to exit monitor gracefully
	ExitChan chan os.Signal

	// RestartModules is the channel to send signal to restart all the modules
	// In this case Run field will direct to load config variable also
	RestartModules chan Status

	//Modules is the list of the modules currently implemented..
	Modules = []string{"live", "target", "disk", "ram", "cpu"}

	// This variable will work like semaphore. If any module is dependent on owtf is turned on it will
	// increase the count. So owtf module will only get shutdown signal if this variable is 0
	owtfModuleDependence = 0
)

/*
Status struct is sent on the channel to controlModule function by different modules
to change the status of the modules
*/
type Status struct {
	Module string
	Run    bool
}

// Perror is used to print the errors.
func Perror(out string) {
	log.Println("[!] Error: " + out)
}

// PLogError is used to print the error when log file is inaccessible
func PLogError(err error) {
	if err != nil {
		fmt.Println("Error in opening log file")
		fmt.Println(err.Error())
		fmt.Println("OWTF Health Monitor will now exit. Run the setup script to" +
			"set up the log and configuration files")
		os.Exit(1)
	}
}

/*
PFileError is used to print the error when monitor do not have sufficient
file permission.
*/
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

// ModuleLogs is used to write the logs of the module in the filename file
func ModuleLogs(filename *os.File, status string) {
	mutex.Lock() //The lock is required to prevent mismatch of logfile during race condition
	log.SetOutput(filename)
	log.Println(status)
	mutex.Unlock()
}

// ModuleError is used to log the errors of the module in the filename file
func ModuleError(filename *os.File, err string, description string) {
	mutex.Lock()
	log.SetOutput(filename)
	log.Println("[!] Error occurred : " + err)
	log.Print(description)
	mutex.Unlock()
}

// OpenLogFile is the utility function to open log file
func OpenLogFile(logFileName string) *os.File {
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		PLogError(err)
	}
	return logFile
}

/*
SendModuleStatus is the utility function to change the state of module. It will
send the status to the main function for required changes.
*/
func SendModuleStatus(module string, status bool) {
	signal := Status{Module: module, Run: status}
	ControlChan <- signal
}

// SendStatusToAllModules send the status to all the modules start/stop
func SendStatusToAllModules(status bool) {
	for _, module := range Modules {
		SendModuleStatus(module, status)
	}
}

/*
CheckConf is the utility function to check the config variable loaded from the
database and if fails, then switch to default.
*/
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
		RestartModules <- Status{Module:"all", Run:true}
	}
}

/*
CheckInstalledPackage will return true if a package with specified commandName
is installed.
*/
func CheckInstalledPackage(commandName string) bool {
	command := exec.Command(commandName, "--help")
	if command.Run() != nil {
		return false
	}
	return true
}

/*
AddOWTFModuleDependence will increase the count of modules depending on owtf
module. In case owtf is not running it will send signal to start owtf module,
if the dependence count positive.
*/
func AddOWTFModuleDependence() {
	owtfMutex.Lock()
	if owtfModuleDependence == 0 {
		SendModuleStatus("owtf", true)
	}
	owtfModuleDependence++
	owtfMutex.Unlock()
}

/*
RemoveOWTFModuleDependence will decrease the count of modules depending on owtf
module. In case owtf is running it will send signal to stop owtf module,
if the dependence count non positive.
*/
func RemoveOWTFModuleDependence() {
	owtfMutex.Lock()
	owtfModuleDependence--
	if owtfModuleDependence == 0 {
		SendModuleStatus("owtf", false)
	}
	owtfMutex.Unlock()
}
