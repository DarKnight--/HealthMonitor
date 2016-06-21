package owtf

import (
	"encoding/json"
	"errors"

	"health_monitor/setup"

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
		data struct {
			RecordsFiltered int `json:"records_filtered"`
		}
	)

	_, response, err := fasthttp.Get(nil, setup.ConfigVars.OWTFAddress+checkTargetPath+target)
	if err != nil {
		// TODO check for error and if OWTF is down shutdown monitor gracefully
		return false, err
	}

	err = json.Unmarshal(response, &data)
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
