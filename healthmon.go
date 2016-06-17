package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"health_monitor/api"
	"health_monitor/disk"
	"health_monitor/live"
	"health_monitor/ram"
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

	flags.NoWebUI = flag.Bool("nowebui", false, "Disables the web ui")
	flags.NoCLI = flag.Bool("nocli", false, "Disables cli")
	flags.Quite = flag.Bool("quite", false, "Disables all notifications except email")

	flag.Parse()

	go webui.RunServer(setup.ConfigVars.Port)
	if (*flags.NoCLI == true) || (*flags.NoWebUI == false) {
		fmt.Printf("[*] Server is up and running at http://127.0.0.1:%s\n", setup.ConfigVars.Port)
	}

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(1)
	go tearDown(exitChan, &wg)
	Init()
	utils.ModuleLogs(setup.MainLogFile, "Running modules from last saved profile")
	runModules(chans, &wg)
	controlModule(chans, &wg)
	setup.SaveStatus()
}

func controlModule(chans [5]chan bool, wg *sync.WaitGroup) {
	api.ControlChan = make(chan utils.Status)
	for {
		data := <-api.ControlChan
		switch data.Module {
		case "live":
			if data.Run && !setup.ModulesStatus.Live {
				setup.ModulesStatus.Live = true
				wg.Add(1)
				utils.ModuleLogs(setup.MainLogFile, "Started live module")
				go live.Live(chans[0], wg)
			} else if setup.ModulesStatus.Live {
				setup.ModulesStatus.Live = false
				utils.ModuleLogs(setup.MainLogFile, "Stopped live module")
				chans[0] <- true
			}
		case "target":
			break
		case "disk":
			if data.Run && !setup.ModulesStatus.Disk {
				setup.ModulesStatus.Disk = true
				wg.Add(1)
				utils.ModuleLogs(setup.MainLogFile, "Started disk module")
				go disk.Disk(chans[2], wg)
			} else if setup.ModulesStatus.Disk {
				setup.ModulesStatus.Disk = false
				utils.ModuleLogs(setup.MainLogFile, "Stopped disk module")
				chans[2] <- true
			}
		case "ram":
			if data.Run && !setup.ModulesStatus.RAM {
				setup.ModulesStatus.RAM = true
				wg.Add(1)
				utils.ModuleLogs(setup.MainLogFile, "Started ram module")
				go ram.RAM(chans[3], wg)
			} else if setup.ModulesStatus.RAM {
				setup.ModulesStatus.RAM = false
				utils.ModuleLogs(setup.MainLogFile, "Stopped ram module")
				chans[3] <- true
			}
		}
	}
}

func runModules(chans [5]chan bool, wg *sync.WaitGroup) {
	if setup.ModulesStatus.Live {
		wg.Add(1)
		utils.ModuleLogs(setup.MainLogFile, "Started live module")
		go live.Live(chans[0], wg)
	}
	if setup.ModulesStatus.Target {
		utils.ModuleLogs(setup.MainLogFile, "Started target module")
	}
	if setup.ModulesStatus.Disk {
		wg.Add(1)
		utils.ModuleLogs(setup.MainLogFile, "Started disk module")
		go disk.Disk(chans[2], wg)
	}
	if setup.ModulesStatus.RAM {
		wg.Add(1)
		utils.ModuleLogs(setup.MainLogFile, "Started ram module")
		go ram.RAM(chans[3], wg)
	}
}

//Init initialises all the modules of the monitor
func Init() {
	live.Init()
	disk.Init()
	ram.Init()
}

func tearDown(exitChan chan os.Signal, wg *sync.WaitGroup) {
	<-exitChan
	utils.ModuleLogs(setup.MainLogFile, "Shutdown signal received.")
	setup.Database.Close()
	setup.SaveStatus()
	utils.ModuleLogs(setup.MainLogFile, "Saved all config data. Stopping running modules")

	var module string
	for module = range api.ConfFunc {
		api.ChangeModuleStatus(module, false)
	}

	setup.MainLogFile.Close()
	setup.DBLogFile.Close()
	wg.Done()
	wg.Wait()
	os.Exit(0)
}
