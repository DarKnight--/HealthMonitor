/* live package is to check conectvity to address*/
package live

import (
	"os/exec"
	"syscall"
	"strconv"
	"time"
	
	"github.com/DarKnight--/HealthMonitor/config"
	"github.com/DarKnight--/HealthMonitor/utils"
)


// The functio is used to send ping request.
func ping(ip string)(bool){
	command := exec.Command("/bin/sh", "-c", "sudo ping " + ip + " -c 1 -W " + strconv.Itoa(config.Params.LagThreshold))
	var waitStatus syscall.WaitStatus
	
	if err := command.Run(); err != nil {
		utils.Perror(err)
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
		  waitStatus = exitError.Sys().(syscall.WaitStatus)
		}
	} else {
		// Command was successful
		waitStatus = command.ProcessState.Sys().(syscall.WaitStatus)
	}
	
	if waitStatus.ExitStatus() == 0{
		return true
	} else{
		return false
	}
}

func CheckConnection(output chan string) {
	var ip = "8.8.8.8"
	for {
		out := ping(ip)
		if out {
			output <-"You are connected to internet"
		} else{
			output <-"check your connection"
		}
		time.Sleep(time.Duration(config.Params.PingInterval) * time.Second)
	}
}

func CheckTarget(targets []string, output chan string) {
	for {
		time.Sleep(time.Duration(config.Params.PingInterval) * time.Second)
		for _, target := range targets{
			output <-pingTarget(target)
		}
	}
}

func pingTarget(ip string)(string){
	out := ping(ip)
	if out {
		return ip + " is up"
	} else{
		return ip + " is down"
	}
}