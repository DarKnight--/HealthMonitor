package live

import (
	"health_monitor/owtf"
	"health_monitor/utils"
)

func downAction() {
	utils.ModuleLogs(logFile, "Sending pause signal to all owtf workers")
	err := owtf.PauseAllWorker()
	if err != nil {
		utils.ModuleError(logFile, "Unable to pause all the workers", err.Error())
	}
}

func upAction() {
	utils.ModuleLogs(logFile, "Sending resume signal to all owtf workers")
	err := owtf.ResumeAllWorker()
	if err != nil {
		utils.ModuleError(logFile, "Unable to resume all the workers", err.Error())
	}
}
