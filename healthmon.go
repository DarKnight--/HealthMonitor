package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_ "strings"

	"github.com/DarKnight--/HealthMonitor/config"
	"github.com/DarKnight--/HealthMonitor/live"
	"github.com/DarKnight--/HealthMonitor/utils"
	"github.com/DarKnight--/HealthMonitor/webui"
)

var (
	c       = make(chan string)
	targets config.TargetData
)

/*Function to used to get json data of targets from OWTF*/
func getTarget() bool {
	//get json data from OWTF target endnode
	response, err := http.Get(config.Params.URL + "api/targets/search/")
	if err != nil {
		utils.Perror(err.Error())
		return false
	}
	defer response.Body.Close()
	// Converting data recieved from http request to byte format
	data1, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Perror(err.Error())
		return false
	}

	// Converting json byte to targets data structure
	err = json.Unmarshal([]byte(data1), &targets)
	if err != nil {
		utils.Perror(err.Error())
		return false
	}
	return true
}

func runModules() {
	var target_ips []string
	
	ret := getTarget()

	if ret == false {
		utils.Perror("Unable to fetch target list from OWTF.")
		utils.Perror("OWTF is not running or data recieved is not of correct form")
		utils.Perror("Skipping target connectivity checks")
	} else {
		for _, target := range targets.Data {
			target_ips = append(target_ips, target.Host_ip)
		}
		log.Println(target_ips)
		go live.CheckTarget(target_ips, c)
	}
	go live.CheckConnection(c)
	go webui.RunServer("8015")
}

func main() {
	runModules()
	for i := range c {
		log.Println(i)
	}
}
