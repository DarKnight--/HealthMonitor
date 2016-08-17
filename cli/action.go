package cli

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/owtf/health_monitor/api"
	"github.com/owtf/health_monitor/setup"
	"github.com/owtf/health_monitor/utils"

	"github.com/fatih/color"
)

const (
	helpString = `Usage: command <arguments>
List of commmands:
help			: To view this message
enable <moduleName>	: To enable a module
disable <moduleName>	: To disable a module
status			: To check status of all modules
status <moduleName>	: To check status of particular module
profile <current/all>	: Use current to get current used profile and all to list all the profiles in database
load <profileName	: To load a particular profile
owtf <resume/pause>	: To send signal to OWASP-OWTF to pause or resume all the workers
disk clean <>		: To clean trash or package manager cache use 'trash' or 'pm_cache'. To do basic cleanup use root or home directory path
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
	if len(argument) == 0 {
		printHeading("OWTF - Health Monitor Module's Status")
		liveShortStatus()
		diskShortStatus()
		cpuShortStatus()
		ramShortStatus()
		targetShortStatus()
	} else if len(argument) == 1 {
		switch argument[0] {
		case "live":
			liveDetailStatus()
		case "disk":
			diskDetailStatus()
		case "cpu":
			cpuDetailStatus()
		case "ram":
			ramDetailStatus()
		case "target":
			targetDetailStatus()
		default:
			color.Red("Module not found")
			color.New(color.FgCyan).Println("Allowed modules: ", utils.Modules)
		}
	} else {
		color.Red("Wrong command used")
		fmt.Println("Usage:")
		fmt.Println("status			to show brief status of all modules")
		fmt.Println("status <moduleName>	to show status of a particular module")
	}
	return nil
}

func liveShortStatus() {
	fmt.Printf("%-35s", "Internet Connectivity (live)")
	moduleWorkingStatus(setup.InternalModuleState.Live, api.LiveStatus().Normal)
}

func liveDetailStatus() {
	printHeading("Status of internet connectivity (live) module")
	fmt.Printf("%-35s", "Internet Connectivity (live) :\t")
	if setup.InternalModuleState.Live {
		moduleStatus := api.LiveStatus()
		if moduleStatus.Normal {
			color.Green("On")
			fmt.Println("You are connected to the internet.")
		} else {
			color.Red("On")
			fmt.Println("You are not connected to the internet.")
		}
	} else {
		color.Cyan("Off")
		fmt.Println("Please turn on the module to monitor internet connectivity status")
	}
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
	moduleWorkingStatus(setup.InternalModuleState.Disk, normal)
}

func diskDetailStatus() {
	printHeading("Status of disk module")
	diskShortStatus()
	if setup.InternalModuleState.Disk {
		printDiskTable()
		fmt.Println()
		printInodeTable()
	}
}

func printDiskTable() {
	fmt.Println("Description of disk blocks available in the system:")
	printLine()
	color.New(color.FgWhite, color.Bold, color.Underline).Printf("| %-30s | %-15s | %-15s |  %%  |\n",
		"Filesystem", "Free Blocks", "Total Blocks")
	if setup.InternalModuleState.Disk {
		for key, value := range api.DiskStatus() {
			colorFunc := color.New(color.FgWhite)
			if value.Status.Space != 1 {
				colorFunc.Add(color.FgRed)
			}
			colorFunc.Printf("| %-30s | %-15d | %-15d | %d%% |\n", key,
				value.Stats.FreeBlocks, value.Const.TotalBlocks,
				percent(value.Stats.FreeBlocks, value.Const.TotalBlocks))
		}
	}
	printLine()
}

func printInodeTable() {
	fmt.Println("Description of disk blocks available in the system:")
	printLine()
	color.New(color.FgWhite, color.Bold, color.Underline).Printf("| %-30s | %-15s | %-15s |  %%  |\n",
		"Filesystem", "Free Inodes", "Total Inodes")
	if setup.InternalModuleState.Disk {
		for key, value := range api.DiskStatus() {
			colorFunc := color.New(color.FgWhite)
			if value.Status.Inode != 1 {
				colorFunc.Add(color.FgRed)
			}
			colorFunc.Printf("| %-30s | %-15d | %-15d | %d%% |\n", key,
				value.Stats.FreeInodes, value.Const.TotalInodes,
				percent(value.Stats.FreeInodes, value.Const.TotalInodes))
		}
	}
	printLine()
}

func printLine() {
	color.New(color.FgWhite, color.Bold, color.Underline).Printf("%76s\n", " ")
}

func percent(value int, total int) int {
	if total == 0 {
		return 0
	}
	return (value * 100) / total
}

func cpuShortStatus() {
	fmt.Printf("%-35s", "CPU (cpu)")
	moduleWorkingStatus(setup.InternalModuleState.CPU, api.CPUStatus().Status.Normal)
}

func cpuDetailStatus() {
	printHeading("Status of CPU module")
	cpuShortStatus()
	if setup.InternalModuleState.CPU {
		moduleStatus := api.CPUStatus()
		colorFunc := color.New(color.FgWhite)
		if moduleStatus.Status.Normal == false {
			colorFunc.Add(color.FgRed)
		}
		colorFunc.Printf("CPU usage is %d%%\n", moduleStatus.Stats.CPUUsage)
	}
}

func ramShortStatus() {
	fmt.Printf("%-35s", "RAM (ram)")
	moduleWorkingStatus(setup.InternalModuleState.RAM, api.RAMStatus().Status.Normal)
}

func ramDetailStatus() {
	printHeading("Status of RAM module")
	ramShortStatus()
	if setup.InternalModuleState.RAM {
		moduleStatus := api.RAMStatus()
		colorFunc := color.New(color.FgWhite)
		if moduleStatus.Status.Normal == false {
			colorFunc.Add(color.FgRed)
		}
		colorFunc.Printf("RAM is %d%% free\n", percent(moduleStatus.Stats.FreePhysical,
			moduleStatus.Consts.TotalPhysical))

		colorFunc.Printf("Swap area is %d%% free\n", percent(moduleStatus.Stats.FreeSwap,
			moduleStatus.Consts.TotalSwap))
	}
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

	moduleWorkingStatus(setup.InternalModuleState.Target, normal)
}

func targetDetailStatus() {
	color.Cyan("Detailed status of all the OWTF's target")
	targetShortStatus()
	if setup.InternalModuleState.Target {
		for key, value := range api.TargetStatus() {
			fmt.Printf("%-45s :\t", key)
			if value.Scanned {
				if value.Normal {
					color.Green("Connected")
				} else {
					color.Red("Not connected")
				}
			} else {
				color.Cyan("Not under scan")
			}
		}
	}
}

func printHeading(heading string) {
	color.New(color.Underline, color.Italic, color.FgCyan, color.Bold).Println(heading)
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
		api.ChangeModuleStatus(module, state)
		return nil
	}
	return errors.New("Specified module not found, allowed modules " + color.New(color.FgCyan).SprintFunc()(utils.Modules))
}

func doesModuleExists(module string) bool {
	for _, workingModule := range utils.Modules {
		if workingModule == module {
			return true
		}
	}
	return false
}

func loadProfile(argument []string) error {
	var err error
	if len(argument) == 1 {
		if err = api.LoadNewProfile(argument[0]); err == nil {
			color.Cyan("Successfully loaded '%s' profile", argument[0])
			return nil
		}
		return err
	}
	return errors.New("Wrong command, use load <profileName>")
}

func manageOWTF(argument []string) error {
	var err error
	if len(argument) == 1 {
		if argument[0] == "resume" {
			if err = api.ResumeOWTF(); err == nil {
				color.Cyan("Successfully sent resume signal to all the workers.")
				return nil
			}
			return err
		} else if argument[0] == "pause" {
			if err = api.PauseOWTF(); err == nil {
				color.Cyan("Successfully sent pause signal to all the workers.")
				return nil
			}
			return err
		}
	}
	return errors.New("Wrong command, use owtf <resume/pause>")
}

func manageDisk(argument []string) error {
	var err error = nil
	if len(argument) == 2 {
		if argument[0] == "clean" {
			switch argument[1] {
			case "/":
				api.BasicDiskCleanup(argument[1])
			case os.Getenv("HOME"):
				api.BasicDiskCleanup(argument[1])
			case "trash":
				err = api.CleanTrashFolder()
			case "pm_cache":
				err = api.DeletePackageManagerCache()
			default:
				err = fmt.Errorf("Argument '%s' is incorrect", argument[1])
			}

			if err == nil {
				color.Cyan("Successfully cleaned %s", argument[1])
				return nil
			}
			return err
		}
	}
	return fmt.Errorf("Wrong command, use disk clean </, %s, trash or pm_cache", os.Getenv("HOME"))
}

func manageProfile(argument []string) error {
	if len(argument) == 1 {
		switch argument[0]{
		case "current":
			fmt.Print("Current Profile: ")
			color.Cyan(setup.UserModuleState.Profile)
		case "all":
			fmt.Print("Saved profiles: ")
			color.Cyan(fmt.Sprintln(setup.GetAllProfiles()))
		case "default":
			return errors.New("Argument not supported. Use <current/all>")
		}
		return nil
	}
	return errors.New("Wrong command, use profile <current/all>")
}