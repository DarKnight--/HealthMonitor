package main

import (
	"flag"
	"fmt"
	"sync"

	"health_monitor/api"
	"health_monitor/disk"
	"health_monitor/live"
	"health_monitor/setup"
	"health_monitor/utils"
	"health_monitor/webui"
)

// Flags holds the health_monitor command line arguments
type Flags struct {
	NoWebUI *bool
	NoCLI   *bool
	Quite   *bool
}

func main() {
	var (
		wg    sync.WaitGroup
		flags Flags
		chans [5]chan bool //Number of modules
	)

	for i := range chans {
		chans[i] = make(chan bool)
	}

	controlModule(chans, &wg)

	flags.NoWebUI = flag.Bool("nowebui", false, "Disables the web ui")
	flags.NoCLI = flag.Bool("nocli", false, "Disables cli")
	flags.Quite = flag.Bool("quite", false, "Disables all notifications except email")

	flag.Parse()

	go webui.RunServer(setup.ConfigVars.Port)
	if (*flags.NoCLI == true) || (*flags.NoWebUI == false) {
		fmt.Printf("[*] Server is up and running at 127.0.0.1:%s\n", setup.ConfigVars.Port)
	}
}

func controlModule(chans [5]chan bool, wg *sync.WaitGroup) {
	api.ControlChan = make(chan utils.Status)
	for {
		data := <-api.ControlChan
		switch data.Module {
		case "live":
			if data.Run {
				wg.Add(1)
				go live.Live(chans[0], wg)
			} else {
				chans[0] <- true
			}
		case "target":
			break
		case "disk":
			if data.Run {
				wg.Add(1)
				go disk.Disk(chans[2], wg)
			} else {
				chans[2] <- true
			}
		}
	}
}
