package target

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/owtf/health_monitor/notify"
	"github.com/owtf/health_monitor/owtf"
	"github.com/owtf/health_monitor/setup"
	"github.com/owtf/health_monitor/utils"
)

type (
	// Status holds the status of the owtf target after the scan
	Status struct {
		Scanned bool
		Normal  bool
	}
)

var (
	targetHash map[string]string
	targetInfo map[string]Status
	lastStatus map[string]bool
	logFile    *os.File
	conf       *Config
)

// Target is the driver function of the module
func Target(status <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var logFileName = path.Join(setup.ConfigVars.HomeDir, "target.log")

	logFile = utils.OpenLogFile(logFileName)
	defer logFile.Close()

	err := owtf.CheckOWTF()
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "Owtf is not running, Signing off")
		setup.InternalModuleState.Target = false
		return
	}
	targetInfo = make(map[string]Status)
	targetHash = make(map[string]string)
	lastStatus = make(map[string]bool)

	utils.ModuleLogs(logFile, "Running with "+conf.Profile+" profile")
	checkTarget()
	for {
		select {
		case <-status:
			utils.ModuleLogs(logFile, "Received signal to turn off. Signing off")
			return
		case <-time.After(time.Millisecond * time.Duration(conf.RecheckThreshold)):
			checkTarget()
			runtime.Gosched()
		}
	}
}

func checkTarget() {
	targets, err := owtf.GetTarget()
	if err != nil {
		utils.ModuleError(logFile, "Unable to get list of targets", err.Error())
	}

	for _, target := range targets {
		status, err := owtf.CheckTarget(target.TargetURL)
		if err != nil {
			utils.ModuleError(logFile, "Unable to check target status", err.Error())
		} else if status {
			if _, isPresent := lastStatus[target.TargetURL]; isPresent {
				lastStatus[target.TargetURL] = targetInfo[target.TargetURL].Normal
			} else {
				lastStatus[target.TargetURL] = true
			}
			hash, ok := targetHash[target.TargetURL]
			if !ok {
				hash = loadTarget(target.TargetURL)
				if hash == "" {
					hash, err = generateHash(target.TargetURL)
					if err != nil {
						utils.ModuleError(logFile, "Unable to get hash for the target",
							"Hash is not in the database, tried for first time")
						continue
					}
					saveTarget(target.TargetURL, hash)
					continue
				}
				targetHash[target.TargetURL] = hash
			}
			result, err := conf.CheckStatus(target.TargetURL, hash)
			if err != nil {
				utils.ModuleError(logFile, "Error occurred during matching hash score for target "+
					target.TargetURL, err.Error())
				continue
			}
			if result {
				targetInfo[target.TargetURL] = Status{Scanned: true, Normal: true}
				utils.ModuleLogs(logFile, fmt.Sprintf("Target %s is up",
					target.TargetURL))
				if lastStatus[target.TargetURL] == false {
					upAction(target.ID)
				}
			} else {
				targetInfo[target.TargetURL] = Status{Scanned: true, Normal: false}
				utils.ModuleLogs(logFile, fmt.Sprintf("Target %s is down",
					target.TargetURL))
				if lastStatus[target.TargetURL] {
					downAction(target.ID)
					notify.SendDesktopAlert("OWTF - Health Monitor", fmt.Sprintf("Target %s seems to be down. Recheck the target.", target.TargetURL), notify.Critical, "")
				}
			}
			continue
		}
		targetInfo[target.TargetURL] = Status{Scanned: false, Normal: true}
		lastStatus[target.TargetURL] = true
	}
}

func generateHash(target string) (string, error) {
	response, err := http.Get(target)
	if err != nil {
		if setup.UserModuleState.Live {
			utils.LiveEmergency <- true
		}
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode/100 != 2 {
		return "", errors.New("Status code returned by target is " + response.Status)
	}
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	hash, err := HashString(body)
	if err != nil {
		return "", err
	}
	return hash, nil
}

/*
GetStatus function is getter function for the targetInfo to send status
of target monitor
*/
func GetStatus() map[string]Status {
	return targetInfo
}

//GetConfJSON returns the json byte array of the module's config
func GetConfJSON() []byte {
	data, err := json.Marshal(LoadConfig())
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the conf struct")
	}
	return data
}

// GetStatusJSON function returns the json string of the targetInfo struct
func GetStatusJSON() []byte {
	data, err := json.Marshal(targetInfo)
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "[!] Check the targetInfo struct")
	}
	return data
}

//Init is the initialization function of the module
func Init() {
	conf = LoadConfig()
	if conf == nil {
		utils.CheckConf(logFile, setup.MainLogFile, "target", &setup.UserModuleState.Profile, setup.Target)
	}
}
