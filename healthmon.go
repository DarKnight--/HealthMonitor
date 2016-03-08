package main

import (
	"io/ioutil"
	"log"
	"net/http"
	_ "strings"
	"encoding/json"
	
	"github.com/DarKnight--/HealthMonitor/config"
	"github.com/DarKnight--/HealthMonitor/utils"
	"github.com/DarKnight--/HealthMonitor/live"
)

var (
	c = make(chan string)
	targets config.TargetData
)

func getTarget() {
	var err error
	
	response, err := http.Get(config.Params.URL + "api/targets/search/")
	defer response.Body.Close()
	utils.Perror(err)
	
	data1, err := ioutil.ReadAll(response.Body)
	utils.Perror(err)
	
	err = json.Unmarshal([]byte(data1), &targets)
	utils.Perror(err)
	
}

func runModules(){
	var target_ips []string
	for _, target := range targets.Data{
		target_ips = append(target_ips, target.Host_ip)
	}
	log.Println(target_ips)
	go live.CheckTarget(target_ips, c)
	go live.CheckConnection(c)
}

func main() {
	//var err error
	getTarget()
	log.Println(targets)
	runModules()
	for i := range c {
		log.Println(i)
	}
}


