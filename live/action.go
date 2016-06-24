package live

import (
	"health_monitor/owtf"
)

func downAction() {
	owtf.PauseAllWorker()
}

func upAction() {
	owtf.ResumeAllWorker()
}
