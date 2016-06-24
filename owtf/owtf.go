package owtf

import (
	"encoding/json"
	"errors"
	"time"

	"health_monitor/setup"
	"health_monitor/utils"

	"github.com/valyala/fasthttp"
)

type (
	// Target contains list of targets in the database recieved from OWTF api
	Target struct {
		ID        int    `json:"id"`
		TargetURL string `json:"target_url"`
	}
)

const (
	targetPath      = "/api/targets/search/"
	checkTargetPath = "/api/worklist/search?target_url="
	workerPath      = "http://127.0.0.1:8010/api/workers/"
)

func init() {
	go monitorOwtf()
}

// GetTarget function calls to the OWTF api to recieve all the targets in the
// OWTF database
func GetTarget() ([]Target, error) {
	var (
		targets []Target
		objmap  map[string]*json.RawMessage
	)
	// get all the tagrget json data from OWTF target endnode
	_, response, err := fasthttp.Get(nil, setup.ConfigVars.OWTFAddress+targetPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &objmap)

	// Converting json byte to targets data structure
	err = json.Unmarshal(*objmap["data"], &targets)
	if err != nil {
		return targets, err
	}
	return targets, nil
}

// CheckTarget checks whether the target in the database actually under the
// scan.
func CheckTarget(target string) (bool, error) {
	var (
		data []struct {
			Id int `json:"id"`
		}
	)

	_, response, err := fasthttp.Get(nil, "http://127.0.0.1:8009/api/targets?target_url=https://google.com")
	if err != nil {
		// TODO check for error and if OWTF is down shutdown monitor gracefully
		return false, err
	}

	err = json.Unmarshal(response, &data)
	if err != nil {
		return false, err
	}

	if data[0].Id > 0 {
		return true, nil
	}

	return false, nil
}

//CheckOWTF will check the running status of owtf based on the api. It will return
// error if OWTF is down or --nowebui option is used
func CheckOWTF() error {
	return getRequest(setup.ConfigVars.OWTFAddress)
}

//PauseWorker will pause the worker with specified worker value
func PauseWorker(worker int) error {
	return getRequest(workerPath + string(worker) + "/pause")
}

//PauseAllWorker will pause all the workers running by OWTF
func PauseAllWorker() error {
	return PauseWorker(0)
}

//ResumeWorker will resume the worker with specified worker value
func ResumeWorker(worker int) error {
	return getRequest(workerPath + string(worker) + "/resume")
}

//ResumeAllWorker will resume all the workers running by OWTF
func ResumeAllWorker() error {
	return ResumeWorker(0)
}

func getRequest(path string) error {
	status, _, err := fasthttp.Get(nil, path)
	if err == nil {
		if status == 200 {
		} else {
			return errors.New("Response code is " + string(status))
		}
		return nil
	}
	return err
}

func monitorOwtf() {
	for true {
		time.Sleep(time.Second)
		if setup.ModulesStatus.Profile != "" {
			break
		}
	}
	time.Sleep(time.Second)

	var (
		workers []struct {
			Busy   bool `json:"busy"`
			Paused bool `json:"busy"`
		}
		owtfStatus bool
		lastStatus bool
	)

	for true {
		status, response, err := fasthttp.Get(nil, workerPath)
		if !(err == nil && status/100 == 2) {
			//OWTF is down
			continue
		}

		err = json.Unmarshal(response, &workers)
		if err != nil {
			utils.ModuleError(setup.MainLogFile, "Unable parse json obtained", err.Error())
			continue
		}
		owtfStatus = false
		lastStatus = true
		//TODO check for free the workers
		for _, worker := range workers {
			if worker.Busy && !worker.Paused {
				owtfStatus = true
				break
			}
		}
		if owtfStatus != lastStatus {
			if owtfStatus {
				startModules()
			} else {
				pauseModules()
			}
		}
		time.Sleep(time.Second)
	}
}

func pauseModules() {
	utils.SendModuleStatus("target", false) //turn of target module
}

func startModules() {
	utils.SendModuleStatus("target", true)
}
