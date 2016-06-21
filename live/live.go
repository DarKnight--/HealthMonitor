package live

import (
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"syscall"
)

// Config holds all the necessary parameters required by the module
type Config struct {
	Profile          string
	HeadURL          string
	RecheckThreshold int // time in milliseconds
	PingThreshold    int // time in milliseconds
	HeadThreshold    int // time in milliseconds
	PingAddress      string
	PingProtocol     string
}

var connection http.Client

// CheckByHEAD will check the internet connectivity by sending a head request
func (l Config) CheckByHEAD() error {
	resp, err := connection.Head(l.HeadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// CheckByDNS check the internet connectivity by resolving the host
// TODO check for dnslookup time, if fooled by local dns server
func (l Config) CheckByDNS() error {
	_, err := net.LookupHost("google.com")
	if err != nil {
		return err
	}
	return nil
}

// Ping check the connectivity by sending ICMP packet to the target
func (l Config) Ping() error {
	var err error
	command := exec.Command("/bin/sh", "-c", "sudo ping "+l.PingAddress+
		" -c 1 -W "+strconv.Itoa(l.PingThreshold/1000))
	var waitStatus syscall.WaitStatus
	if err = command.Run(); err != nil {
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			return err
		}
	} else {
		// Command was successful
		waitStatus = command.ProcessState.Sys().(syscall.WaitStatus)
	}

	if waitStatus.ExitStatus() == 0 {
		return nil
	}
	return err
}
