package live

import (
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"syscall"

	"health_monitor/utils"
)

// Config holds all the necessary parameters required by the module
type Config struct {
	headURL          string
	recheckThreshold int // time in milliseconds
	pingThreshold    int // time in milliseconds
	headThreshold    int // time in milliseconds
	pingProtocol     string
	pingAddress      string
	connection       http.Client
}

func (l Config) checkByHEAD() bool {
	resp, err := l.connection.Head(l.headURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	return true
}

// TODO check for dnslookup time, if fooled by local dns server
func (l Config) checkByDNS() bool {
	ipList, err := net.LookupHost("google.com")
	if err != nil {
		log.Fatal(err)
		return false
	}
	log.Println(ipList)
	return true
}

/* The function is used to send ping request.*/
func (l Config) ping() bool {
	command := exec.Command("/bin/sh", "-c", "sudo ping "+l.pingAddress+
		" -c 1 -W "+strconv.Itoa(l.pingThreshold))
	var waitStatus syscall.WaitStatus

	if err := command.Run(); err != nil {
		utils.Perror(err.Error())
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
		}
	} else {
		// Command was successful
		waitStatus = command.ProcessState.Sys().(syscall.WaitStatus)
	}

	if waitStatus.ExitStatus() == 0 {
		return true
	}
	return false
}
