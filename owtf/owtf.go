package owtf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/owtf/health_monitor/setup"
	"github.com/owtf/health_monitor/utils"
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
	if err := CheckOWTF(); err == nil {
		return getRequest(workerPath + strconv.Itoa(worker) + "/pause")
	}
	return nil
}

// PauseWorkerByTarget send the signal to owtf to pause the worker working on the
// target with specified id
func PauseWorkerByTarget(id int) error {
	workerID, paused := getWorkerByTarget(id)
	if workerID == -1 {
		return errors.New("Unable to get the worker with target id = " + strconv.Itoa(id))
	}
	if paused {
		return nil
	}
	return PauseWorker(workerID)
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

// ResumeWorkerByTarget send the signal to owtf to resume the worker working on the
// target with specified id
func ResumeWorkerByTarget(id int) error {
	workerID, paused := getWorkerByTarget(id)
	if workerID == -1 {
		return errors.New("Unable to get the worker with target id = " + strconv.Itoa(id))
	}
	if paused {
		return ResumeWorker(workerID)
	}
	return nil
}

// PauseOWTF will pause the OWTF and write the logs to the specified file
func PauseOWTF(logFile *os.File) {
	utils.ModuleLogs(logFile, "Sending pause signal to all owtf workers")
	err := PauseAllWorker()
	if err != nil {
		utils.ModuleError(logFile, "Unable to pause all the workers", err.Error())
	}
}

// ResumeOWTF will resume the OWTF and write the logs to the specified file
func ResumeOWTF(logFile *os.File) {
	utils.ModuleLogs(logFile, "Sending resume signal to all owtf workers")
	err := ResumeAllWorker()
	if err != nil {
		utils.ModuleError(logFile, "Unable to resume all the workers", err.Error())
	}
}

func getWorkerByTarget(id int) (int, bool) {
	var (
		workers []struct {
			ID     int  `json:"id"`
			Paused bool `json:"paused"`
			Work   []struct {
				ID int `json:"id"`
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
		if len(worker.Work) != 0 && worker.Work[0].ID == id {
			return worker.ID, worker.Paused
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
