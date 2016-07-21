package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "health_monitor/api"
	"health_monitor/cli"
	"health_monitor/cpu"
	"health_monitor/disk"
	"health_monitor/live"
	"health_monitor/owtf"
	"health_monitor/ram"
	"health_monitor/setup"
	"health_monitor/target"
	"health_monitor/utils"
	"health_monitor/webui"
)

// Flags holds the health_monitor command line arguments
type Flags struct {
	NoWebUI *bool
	NoCLI   *bool
	Quite   *bool
	Install *bool
}

func main() {
	defer func() {
		fmt.Println("In case the program exited due to dependency failure try running '-install' option")
	}()
	var (
		wg    sync.WaitGroup
		flags Flags
		chans [6]chan bool //Number of modules
	)

	for i := range chans {
		chans[i] = make(chan bool)
	}

	utils.LiveEmergency = make(chan bool)
	defer close(utils.LiveEmergency)

	flags.NoWebUI = flag.Bool("nowebui", false, "Disables the web ui")
	flags.NoCLI = flag.Bool("nocli", false, "Disables cli")
	flags.Quite = flag.Bool("quite", false, "Disables all notifications except email")
	flags.Install = flag.Bool("install", false, "Installs necessary dependencies")

	flag.Parse()

	if *flags.Install {
		install()
	}

	if (*flags.NoCLI == true) || (*flags.NoWebUI == false) {
		go webui.RunServer(setup.ConfigVars.Port)
		fmt.Printf("[*] Server is starting at http://127.0.0.1:%s\n", setup.ConfigVars.Port)
	}

	utils.ExitChan = make(chan os.Signal, 1)
	signal.Notify(utils.ExitChan, syscall.SIGINT, syscall.SIGTERM)
	// The buffer size should atleast be double the number of modules implemented
	utils.ControlChan = make(chan utils.Status, 12)
	wg.Add(1)
	go tearDown(&wg)
	Init()
	utils.ModuleLogs(setup.MainLogFile, "Running modules from last saved profile")
	runModules(chans, &wg)
	go controlModule(chans, &wg)
	if *flags.NoCLI == false {
		cli.Run()
	}
	wg.Wait()
}

func controlModule(chans [6]chan bool, wg *sync.WaitGroup) {
	for {
		data := <-utils.ControlChan
		switch data.Module {
		case "owtf":
			controlModuleHelper(data.Run, &setup.OWTFModuleStatus, data.Module,
				owtf.OWTF, chans[0], wg)
		case "live":
			controlModuleHelper(data.Run, &setup.InternalModuleState.Live, data.Module,
				live.Live, chans[1], wg)
		case "target":
			controlModuleHelper(data.Run, &setup.InternalModuleState.Target, data.Module,
				target.Target, chans[2], wg)
		case "disk":
			controlModuleHelper(data.Run, &setup.InternalModuleState.Disk, data.Module,
				disk.Disk, chans[3], wg)
		case "ram":
			controlModuleHelper(data.Run, &setup.InternalModuleState.RAM, data.Module,
				ram.RAM, chans[4], wg)
		case "cpu":
			controlModuleHelper(data.Run, &setup.InternalModuleState.CPU, data.Module,
				cpu.CPU, chans[5], wg)
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

func runModules(chans [6]chan bool, wg *sync.WaitGroup) {
	if setup.UserModuleState.Live {
		wg.Add(1)
		utils.ModuleLogs(setup.MainLogFile, "Started live module")
		setup.InternalModuleState.Live = true
		go live.Live(chans[1], wg)
	}
	if setup.UserModuleState.Target {
		wg.Add(1)
		utils.ModuleLogs(setup.MainLogFile, "Started target module")
		setup.InternalModuleState.Target = true
		utils.IncreaseOWTFDependenceCount()
		go target.Target(chans[2], wg)
	}
	if setup.UserModuleState.Disk {
		wg.Add(1)
		utils.ModuleLogs(setup.MainLogFile, "Started disk module")
		setup.InternalModuleState.Disk = true
		go disk.Disk(chans[3], wg)
	}
	if setup.UserModuleState.RAM {
		wg.Add(1)
		utils.ModuleLogs(setup.MainLogFile, "Started ram module")
		setup.InternalModuleState.RAM = true
		go ram.RAM(chans[4], wg)
	}
	if setup.UserModuleState.CPU {
		wg.Add(1)
		utils.ModuleLogs(setup.MainLogFile, "Started cpu module")
		setup.UserModuleState.CPU = true
		go cpu.CPU(chans[5], wg)
	}
}

//Init initialises all the modules of the monitor
func Init() {
	live.Init()
	target.Init()
	disk.Init()
	ram.Init()
	cpu.Init()
}

func tearDown(wg *sync.WaitGroup) {
	<-utils.ExitChan
	utils.ModuleLogs(setup.MainLogFile, "Shutdown signal received.")
	setup.Database.Close()
	setup.SaveStatus()
	utils.ModuleLogs(setup.MainLogFile, "Saved all config data. Stopping running modules")

	var module string
	for _, module = range utils.Modules {
		utils.SendModuleStatus(module, false)
	}

	utils.SendModuleStatus("owtf", false)
	setup.MainLogFile.Close()
	setup.DBLogFile.Close()
	wg.Done()
	wg.Wait()
	os.Exit(0)
}
