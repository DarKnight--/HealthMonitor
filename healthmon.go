package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"health_monitor/api"
	"health_monitor/cpu"
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

	utils.LiveEmergency = make(chan bool)
	defer close(utils.LiveEmergency)

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
	utils.ControlChan = make(chan utils.Status)
	for {
		data := <-utils.ControlChan
		switch data.Module {
		case "live":
			break
		case "target":
			break
		case "disk":
			controlModuleHelper(data.Run, &setup.ModulesStatus.Disk, data.Module,
				disk.Disk, chans[2], wg)
		case "ram":
			controlModuleHelper(data.Run, &setup.ModulesStatus.RAM, data.Module,
				ram.RAM, chans[3], wg)
		case "cpu":
			controlModuleHelper(data.Run, &setup.ModulesStatus.CPU, data.Module,
				cpu.CPU, chans[4], wg)
		}
	}
}

func controlModuleHelper(run bool, moduleStatus *bool, moduleName string,
	module func(<-chan bool, *sync.WaitGroup), channel chan bool, wg *sync.WaitGroup) {
	if run && !*moduleStatus {
		*moduleStatus = true
		wg.Add(1)
		utils.ModuleLogs(setup.MainLogFile, "Started "+moduleName+"module")
		go module(channel, wg)
	} else if *moduleStatus {
		*moduleStatus = false
		utils.ModuleLogs(setup.MainLogFile, "Stopped "+moduleName+"module")
		channel <- true
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
	if setup.ModulesStatus.CPU {
		wg.Add(1)
		utils.ModuleLogs(setup.MainLogFile, "Started cpu module")
		go cpu.CPU(chans[4], wg)
	}
}

//Init initialises all the modules of the monitor
func Init() {
	live.Init()
	disk.Init()
	ram.Init()
	cpu.Init()
}

func tearDown(exitChan chan os.Signal, wg *sync.WaitGroup) {
	<-exitChan
	utils.ModuleLogs(setup.MainLogFile, "Shutdown signal received.")
	setup.Database.Close()
	setup.SaveStatus()
	utils.ModuleLogs(setup.MainLogFile, "Saved all config data. Stopping running modules")

	var module string
	for _, module = range utils.Modules {
		api.ChangeModuleStatus(module, false)
	}

	setup.MainLogFile.Close()
	setup.DBLogFile.Close()
	wg.Done()
	wg.Wait()
	os.Exit(0)
}
