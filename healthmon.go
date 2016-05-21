package main

import (
	"fmt"
	_ "strings"
	"sync"

	"health_monitor/config"
	"health_monitor/disk"
	"health_monitor/live"
	"health_monitor/utils"
)

func main() {
	var wg sync.WaitGroup
	config.Database.Ping()
	signal := make(chan utils.Status)
	wg.Add(1)
	go live.Live(signal, &wg)
	wg.Add(1)
	go disk.Disk(signal, &wg)
	fmt.Println("hey")
	wg.Wait()
}
