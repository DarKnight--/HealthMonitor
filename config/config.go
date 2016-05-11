package config

import (
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

var (
	// ConfigVars will hold necessary variables loaded from config file
	ConfigVars struct {
		HomeDir     string
		DBFile      string
		OWTFAddress string
	}
)

func init() {
	var configFile = path.Join(os.Getenv("HOME"), ".owtfMonitor", "config",
		"config.toml") // The necessary config file required by health_monitor

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("The config file is missing. Creating one with default settings")
		os.Exit(1) // TODO remove it with a function for creating a config file
	}

	_, err := toml.DecodeFile(configFile, &ConfigVars) // Read the config file
	if err != nil {
		fmt.Println(err)
		fmt.Println("The config file is corupt. Do you want a remove all files" +
			"and setup health_monitor again (y/N)?")
		os.Exit(1) // TODO remove it with appropriate function
	}
}
