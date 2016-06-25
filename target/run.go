package target

import (
	"encoding/json"
	"errors"
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

	err := owtf.CheckOWTF()
	if err != nil {
		utils.ModuleError(logFile, err.Error(), "Owtf is not running, Signing off")
		setup.ModulesStatus.Target = false
		return
	}
	targetInfo = make(map[string]TargetStatus)
	targetHash = make(map[string]string)

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
		} else if status {
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
			}
		}
	}
}

func generateHash(target string) (string, error) {
	status, response, err := fasthttp.Get(nil, target)
	if err != nil {
		utils.LiveEmergency <- true
		return "", err
	}
	if status/100 != 2 {
		return "", errors.New("Status code returned by target is " + string(status))
	}
	hash, err := HashString(response)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func compareTargetHash(target string, hash string) bool {
	newHash, err := generateHash(target)
	if err != nil {
		utils.ModuleError(logFile, "Unable to generate hash for target "+target, err.Error())
		return false
	}
	result := CompareHash(hash, newHash)
	if result == -1 {
		utils.ModuleError(logFile, "Unable to compare hashes", "Please check the target")
		return false
	} else if result < conf.FuzzyThreshold {
		utils.ModuleError(logFile, "Target "+target+" is possible down or blocking OWTF", "Check the logs of OWTF and target")
		//TODO action for target
		return false
	}
	utils.ModuleLogs(logFile, "Target "+target+" is up")
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
	if conf == nil {
		utils.CheckConf(logFile, setup.MainLogFile, "target", &setup.ModulesStatus.Profile, setup.SetupTarget)
	}
}
