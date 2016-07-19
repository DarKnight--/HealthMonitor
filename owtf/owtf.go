package owtf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"health_monitor/setup"
	"health_monitor/utils"
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

func Init() {
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
	response, err := http.Get(setup.ConfigVars.OWTFAddress + targetPath)
	if err != nil {
		return nil, err
	}
	var dataByte []byte
	dataByte, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dataByte, &objmap)

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
		data struct {
			RecordsFiltered int `json:"records_filtered"`
		}
	)

	response, err := http.Get(setup.ConfigVars.OWTFAddress + checkTargetPath + target)
	if err != nil {
		return false, err
	}
	var dataByte []byte
	dataByte, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(dataByte, &data)
	if err != nil {
		return false, err
	}

	if data.RecordsFiltered > 0 {
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
	return getRequest(workerPath + strconv.Itoa(worker) + "/pause")
}

func PauseWorkerByTarget(id int) error {
	workerId, paused := getWorkerByTarget(id)
	if workerId == -1 {
		return errors.New("Unable to get the worker with target id = " + strconv.Itoa(id))
	}
	if paused {
		return nil
	}
	return PauseWorker(workerId)
}

//PauseAllWorker will pause all the workers running by OWTF
func PauseAllWorker() error {
	return PauseWorker(0)
}

//ResumeWorker will resume the worker with specified worker value
func ResumeWorker(worker int) error {
	return getRequest(workerPath + strconv.Itoa(worker) + "/resume")
}

//ResumeAllWorker will resume all the workers running by OWTF
func ResumeAllWorker() error {
	return ResumeWorker(0)
}

func ResumeWorkerByTarget(id int) error {
	workerId, paused := getWorkerByTarget(id)
	if workerId == -1 {
		return errors.New("Unable to get the worker with target id = " + strconv.Itoa(id))
	}
	if paused {
		return ResumeWorker(workerId)
	}
	return nil
}

func getWorkerByTarget(id int) (int, bool) {
	var (
		workers []struct {
			Id     int  `json:"id"`
			Paused bool `json:"paused`
			Work   []struct {
				Id int `json:"id"`
			} `json:"work"`
		}
	)

	response, err := http.Get(workerPath)
	if !(err == nil && response.StatusCode/100 == 2) {
		return -1, false
	}
	var dataByte []byte
	dataByte, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return -1, false
	}

	err = json.Unmarshal(dataByte, &workers)
	if err != nil {
		return -1, false
	}
	for _, worker := range workers {
		if len(worker.Work) != 0 && worker.Work[0].Id == id {
			return worker.Id, worker.Paused
		}
	}
	return -1, false
}

func getRequest(path string) error {
	response, err := http.Get(path)
	if err == nil {
		defer response.Body.Close()
		if response.StatusCode == 200 {
		} else {
			return errors.New("Response code is " + response.Status)
		}
		return nil
	}
	return err
}

func monitorOwtf() {
	var (
		workers []struct {
			Busy   bool `json:"busy"`
			Paused bool `json:"paused"`
		}
		owtfStatus bool
		lastStatus bool
	)
	lastStatus = true
	for true {
		response, err := http.Get(workerPath)
		if !(err == nil && response.StatusCode/100 == 2) {
			//OWTF is down
			continue
		}
		var dataByte []byte
		dataByte, err = ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			continue
		}

		err = json.Unmarshal(dataByte, &workers)
		if err != nil {
			utils.ModuleError(setup.MainLogFile, "Unable parse json obtained", err.Error())
			continue
		}
		owtfStatus = false
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
		lastStatus = owtfStatus
		time.Sleep(time.Second)
	}
}

func pauseModules() {
	utils.SendModuleStatus("target", false) //turn off target module
}

func startModules() {
	utils.SendModuleStatus("target", true)
}
