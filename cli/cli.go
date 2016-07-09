package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"health_monitor/utils"
)

var (
	cliFunctions map[string]func([]string) error
)

func init() {
	cliFunctions = make(map[string]func([]string) error)
	cliFunctions["exit"] = exit
	cliFunctions["disable"] = disableModule
	cliFunctions["enable"] = enableModule
	cliFunctions["help"] = help
}

func Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			utils.Perror(err.Error())
		}
		command := strings.Split(text[:len(text)-1], " ")
		if len(command) == 0 {
			continue
		}
		manageCommand(command)
	}
}

func manageCommand(command []string) {
	if function, ok := cliFunctions[command[0]]; ok {
		if err := function(command[1:]); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Command not found. Use 'help' to know more.")
	}
}
