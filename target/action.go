package target

import (
	"github.com/owtf/health_monitor/owtf"
	"github.com/owtf/health_monitor/utils"
)

func downAction(targetID int) {
	utils.ModuleLogs(logFile, "Sending pause signal to all owtf workers")
	err := owtf.PauseWorkerByTarget(targetID)
	if err != nil {
		utils.ModuleError(logFile, "Unable to pause the workers", err.Error())
	}
}

func upAction(targetID int) {
	err := owtf.ResumeWorkerByTarget(targetID)
	if err != nil {
		utils.ModuleError(logFile, "Unable to resume the worker", err.Error())
	}
}
