package ram

import (
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"health_monitor/setup"
	"health_monitor/utils"
)

type (
	Status struct {
		Normal bool
	}

	Info struct {
		Status Status
	}
)

var (
	ramInfo Info
	logFile *os.File
	conf    *Config
)

func Ram(status chan utils.Status, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		logFileName = path.Join(setup.ConfigVars.HomeDir, "disk.log")
		err         error
	)

	logFile = utils.OpenLogFile(logFileName)
	defer logFile.Close()

	utils.ModuleLogs(logFile, "Running with "+conf.Profile+" profile")
	for {
		select {
		case signal := <-status:
			utils.ModuleLogs(logFile, "Recieved signal to turn off. Signing off")
			return
		case <-time.After(time.Millisecond * time.Duration(conf.RecheckThreshold)):
			checkRam()
			runtime.Gosched()
		}
	}
}

func checkRam() {
}

func Init() {
	conf = LoadConfig()
}
