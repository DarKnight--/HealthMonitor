package cli

import (
	"errors"
	"fmt"
	"syscall"

	"health_monitor/utils"
)

const (
	helpString = `CLI for the OWTF - Health Monitor.
Usage: command <arguments>
List of commmands:
help	: To view this message`
)

func disableModule(argument []string) error {
	if len(argument) > 1 {
		return errors.New("Wrong command, use disable <moduleName>")
	}
	return toggleModule(argument[0], false)
}

func exit(argument []string) error {
	utils.ExitChan <- syscall.SIGINT
	return nil
}

func enableModule(argument []string) error {
	if len(argument) > 1 {
		return errors.New("Wrong command, use enable <moduleName>")
	}
	return toggleModule(argument[0], true)
}

func help(argument []string) error {
	fmt.Println(helpString)
	return nil
}

func toggleModule(module string, state bool) error {
	if doesModuleExists(module) {
		utils.SendModuleStatus(module, false)
		return nil
	} else {
		return errors.New("Specified module not found, allowed modules " + fmt.Sprint(utils.Modules))
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
