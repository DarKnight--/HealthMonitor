package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"health_monitor/config"
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
		err     error
	)
	// get all the tagrget json data from OWTF target endnode
	var response *http.Response
	response, err = http.Get(config.ConfigVars.OWTFAddress + path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	// Converting data recieved from http request to byte format
	var dataByte []byte
	dataByte, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal(dataByte, &objmap)

	// Converting json byte to targets data structure
	err = json.Unmarshal(*objmap["data"], &targets)
	if err != nil {
		log.Println(err)
		return targets, err
	}
	return targets, nil
}

// CheckTarget checks whether the target in the database actually under the
// scan.
func CheckTarget(target string) bool {
	const path = "/api/worklist/search?target_url="
	var (
		err      error
		response *http.Response
		data     struct {
			RecordsFiltered int `json:"records_filtered"`
		}
	)

	response, err = http.Get(config.ConfigVars.OWTFAddress + path + target)
	if err != nil {
		log.Println(err)
		// TODO check for error and if OWTF is down shutdown monitor gracefully
		return false
	}
	defer response.Body.Close()
	var dataByte []byte
	dataByte, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	err = json.Unmarshal(dataByte, &data)
	if err != nil {
		log.Println("Error occured during decoding")
		return false
	}

	if data.RecordsFiltered > 0 {
		return true
	}

	return false
}
