package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"health_monitor/setup"
)

type (
	// Target contains list of targets in the database recieved from OWTF api
	Target struct {
		ID        int    `json:"id"`
		TargetURL string `json:"target_url"`
	}
)

// GetTarget function calls to the OWTF api to recieve all the targets in the
// OWTF database
func GetTarget() ([]Target, error) {
	const path = "/api/targets/search/"
	var (
		targets []Target
		objmap  map[string]*json.RawMessage
	)
	// get all the tagrget json data from OWTF target endnode
	response, err := http.Get(setup.ConfigVars.OWTFAddress + path)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Converting data recieved from http request to byte format
	dataByte, err := ioutil.ReadAll(response.Body)
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
	const path = "/api/worklist/search?target_url="
	var (
		data struct {
			RecordsFiltered int `json:"records_filtered"`
		}
	)

	response, err := http.Get(setup.ConfigVars.OWTFAddress + path + target)
	if err != nil {
		// TODO check for error and if OWTF is down shutdown monitor gracefully
		return false, err
	}
	defer response.Body.Close()
	var dataByte []byte
	dataByte, err = ioutil.ReadAll(response.Body)
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
