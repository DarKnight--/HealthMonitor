package cli

import (
	"errors"
	"fmt"
	"syscall"

	"health_monitor/api"
	"health_monitor/setup"
	"health_monitor/utils"

	"github.com/fatih/color"
)

const (
	helpString = `Usage: command <arguments>
List of commmands:
help			: To view this message
enable <moduleName>	: To enable a module
disable <moduleName>	: To disable a module
status			: To check status of all modules
exit			: To turn off the monitor`
)

func disableModule(argument []string) error {
	if len(argument) > 1 {
		return errors.New("Wrong command, use disable <moduleName>")
	}
	return toggleModule(argument[0], false)
}

func exit(argument []string) error {
	color.Blue("Shutting down monitor gacefully")
	utils.ExitChan <- syscall.SIGINT
	return nil
}

func enableModule(argument []string) error {
	if len(argument) > 1 {
		return errors.New("Wrong command, use" + color.CyanString("enable <moduleName>"))
	}
	return toggleModule(argument[0], true)
}

func help(argument []string) error {
	color.Green("CLI for the OWTF - Health Monitor.")
	color.Cyan(helpString)
	return nil
}

func status(argument []string) error {
	color.New(color.Underline, color.Italic, color.FgCyan, color.Bold).Println("OWTF - Health Monitor Module's Status")
	liveShortStatus()
	diskShortStatus()
	cpuShortStatus()
	ramShortStatus()
	targetShortStatus()
	return nil
}

func liveShortStatus() {
	fmt.Printf("%-35s", "Internet Connectivity (live)")
	moduleWorkingStatus(setup.ModulesStatus.Live, api.LiveStatus().Normal)
}

func diskShortStatus() {
	normal := true
	fmt.Printf("%-35s", "Disk (disk)")
	for _, value := range api.DiskStatus() {
		if value.Status.Inode != 1 || value.Status.Space != 1 {
			normal = false
			break
		}
	}
	moduleWorkingStatus(setup.ModulesStatus.Disk, normal)
}

func cpuShortStatus() {
	fmt.Printf("%-35s", "CPU (cpu)")
	moduleWorkingStatus(setup.ModulesStatus.CPU, api.CPUStatus().Status.Normal)
}

func ramShortStatus() {
	fmt.Printf("%-35s", "RAM (ram)")
	moduleWorkingStatus(setup.ModulesStatus.RAM, api.RAMStatus().Status.Normal)
}

func targetShortStatus() {
	fmt.Printf("%-35s", "OWTF's Targets (target)")
	normal := true

	for _, value := range api.TargetStatus() {
		if value.Scanned {
			if value.Normal == false {
				normal = false
			}
		}
	}

	moduleWorkingStatus(setup.ModulesStatus.Target, normal)
}

func moduleWorkingStatus(status bool, workingStatus bool) {
	fmt.Print(":\t")
	if status {
		if workingStatus {
			color.Green("On")
		} else {
			color.Red("On")
		}
	} else {
		color.Cyan("Off")
	}
}

func toggleModule(module string, state bool) error {
	if doesModuleExists(module) {
		utils.SendModuleStatus(module, state)
		return nil
	} else {
		return errors.New("Specified module not found, allowed modules " + color.New(color.FgCyan).SprintFunc()(utils.Modules))
	}
}

func doesModuleExists(module string) bool {
	for _, workingModule := range utils.Modules {
		if workingModule == module {
			return true
		}
	}
	return false
}
