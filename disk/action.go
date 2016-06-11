package disk

import (
	"os"
	"reflect"

	"health_monitor/setup"
	"health_monitor/utils"
)

const (
	DebianAPTPath = "/var/cache/apt/archives"
)

var (
	kali = []byte{75, 97, 108, 105, 10}
)

func removeAptCache() {
	if reflect.DeepEqual(setup.OSVarient, kali) {
		err := os.RemoveAll(DebianAPTPath)
		if err != nil {
			utils.ModuleError(logFile, "Unable to free apt cache", err.Error())
		}
		utils.ModuleLogs(logFile, "Deleted apt cache successfully.")
	}
}
