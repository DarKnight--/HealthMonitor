package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type (
	Target struct {
		Id         int    `json:"id"`
		Target_url string `json:"target_url"`
	}
)

func GetTarget() ([]Target, error) {
	var targets []Target
	var objmap map[string]*json.RawMessage
	var err error
	//get json data from OWTF target endnode
	var response *http.Response
	response, err = http.Get("http://127.0.0.1:8009/" + "api/targets/search/")
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
