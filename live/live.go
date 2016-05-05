package live

import (
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"HealthMonitor/utils"
)

var (
	DAFAULT = 0 // 0-> GET, 1-> HEAD, 2-> Ping test
)

type live struct {
	headURL          string
	recheckThreshold int
	pingThreshold    int
	headThreshold    int
	pingProtocol     string
	pingAddress      string
	connection       http.Client
}

// TODO add concurrency and return parameters accordingly

func Live() *live {
	l := new(live)
	l.headURL = "https://google.com"
	l.headThreshold = 4
	l.connection = http.Client{
		Timeout: time.Duration(l.headThreshold) * time.Second,
	}
	// setup the parameters
	return l
}

func (l live) checkByHEAD() {
	resp, err := l.connection.Head(l.headURL)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
}

// TODO check for dnslookup time, if fooled by local dns server
func (l live) checkByDNS() {
	ipList, err := net.LookupHost("google.com")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ipList)
}

/* The function is used to send ping request.*/
func (l live) ping() bool {
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
	} else {
		return false
	}
}

/*
func main() {
	var l *live = Live()
	l.checkByHEAD()
	l.checkByDNS()
}
*/
