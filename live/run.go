package live

import (
	"time"

	"health_monitor/utils"
)

// Live is the driver function of this module for monitor
func Live(status chan utils.Status) {
	var live Config
	Default := live.checkByHEAD
	if live.checkByDNS() {
		Default = live.checkByDNS
	}
	if live.ping() {
		Default = live.ping
	}

	select {
	case signal := <-status:
		if signal.Module == 1 && signal.Run == false {
			return
		}

	case <-time.After(time.Second * time.Duration(live.recheckThreshold)):
		Default() // TODO logging option to complete
	}
}
