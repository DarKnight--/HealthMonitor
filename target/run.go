package target

import (
	"encoding/json"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"health_monitor/owtf"
	"health_monitor/setup"
	"health_monitor/utils"

	"github.com/valyala/fasthttp"
)

type (
	TargetStatus struct {
		Scanned bool
		Normal  bool
	}
)

var (
	targetHash map[string]string
	targetInfo map[string]TargetStatus
	logFile    *os.File
	conf       *Config
)

func Target(status <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var logFileName = path.Join(setup.ConfigVars.HomeDir, "target.log")

	logFile = utils.OpenLogFile(logFileName)
	defer logFile.Close()

	targetInfo = make(map[string]TargetStatus)

	utils.ModuleLogs(logFile, "Running with "+conf.Profile+" profile")
	checkTarget()
	for {
		select {
		case <-status:
			utils.ModuleLogs(logFile, "Recieved signal to turn off. Signing off")
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
		//TODO check owtf status
	}

	for _, target := range targets {
		status, err := owtf.CheckTarget(target.TargetURL)
		if err != nil {
			utils.ModuleError(logFile, "Unable to check target status", err.Error())
		}
		if status {
			hash, ok := targetHash[target.TargetURL]
			if !ok {
				hash = generateHash(target.TargetURL)
				if len(hash) == 0 {
					//TODO Alert to do
					continue
				}
				targetHash[target.TargetURL] = hash
				// Save this hash to database
				continue
			}
			result := compareTargetHash(target.TargetURL, hash)
			if result {
				targetInfo[target.TargetURL] = TargetStatus{Scanned: true, Normal: true}
			} else {
				targetInfo[target.TargetURL] = TargetStatus{Scanned: true, Normal: false}
				//TODO action for target
			}
		}
	}
}

func generateHash(target string) string {
	status, response, err := fasthttp.Get(nil, "http://localhost:8009")
	if err != nil {
		utils.LiveEmergency <- true
		return ""
	}
	hash, err := HashString(response)
	if err != nil {
		return ""
	}
	return hash
}

func compareTargetHash(target string, hash string) bool {
	newHash := generateHash(target)
	if len(newHash) == 0 {
		// TODO Alert
		return false
	}
	result := CompareHash(hash, newHash)
	if result == -1 {
		// TODO Alert
		return false
	} else if result < conf.FuzzyThreshold {
		// TODO Alert for major change in target
		return false
	}
	return true
}

// GetStatus function is getter funtion for the targetInfo to send status
// of target monitor
func GetStatus() map[string]TargetStatus {
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

// GetStatusJSON function retuns the json string of the targetInfo struct
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
}
