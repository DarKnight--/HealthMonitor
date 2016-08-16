package owtf

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/owtf/health_monitor/setup"
	"github.com/owtf/health_monitor/utils"
)

// OWTF is the driver function of the owtf module.
// It continuosly monitors the status of the OWTF whether it is scanning any target or not
func OWTF(status <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var lastStatus = true

	for {
		select {
		case <-status:
			return
		case <-time.After(time.Second):
			monitorOwtf(&lastStatus)
		}
	}
}

func monitorOwtf(lastStatus *bool) {
	var (
		workers []struct {
			Busy   bool `json:"busy"`
			Paused bool `json:"paused"`
		}
		owtfStatus bool
	)
	response, err := http.Get(workerPath)
	if !(err == nil && response.StatusCode/100 == 2) {
		//OWTF is down
		return
	}
	var dataByte []byte
	dataByte, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return
	}

	err = json.Unmarshal(dataByte, &workers)
	if err != nil {
		utils.ModuleError(setup.MainLogFile, "Unable parse json obtained", err.Error())
		return
	}
	owtfStatus = false
	//TODO check for free the workers
	for _, worker := range workers {
		if worker.Busy && !worker.Paused {
			owtfStatus = true
			break
		}
	}

	if owtfStatus != *lastStatus {
		if owtfStatus {
			startModules()
		} else {
			pauseModules()
		}
	}
	*lastStatus = owtfStatus
	time.Sleep(time.Second)
}

func pauseModules() {
	utils.SendModuleStatus("target", false) //turn off target module
}

func startModules() {
	utils.SendModuleStatus("target", true)
}
