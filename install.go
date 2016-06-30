package main

import (
	"fmt"
	"os"
	"os/exec"

	"health_monitor/utils"
)

const (
	tempPath = "/tmp/monitor"
)

func install() {
	if os.Getegid() != 0 {
		utils.Perror("Please run the program with sudo privileges")
		os.Exit(1)
	}

	complete := true
	fmt.Println("[-]Installing ssdeep package...")
	command := exec.Command("/bin/sh", "-c", `wget http://downloads.sourceforge.net/project/ssdeep/ssdeep-2.13/ssdeep-2.13.tar.gz; \
		tar zxvf ssdeep-2.13.tar.gz; cd ssdeep-2.13; ./configure; make; make install `)
	command.Stdout = os.Stdout
	command.Dir = tempPath
	if err := command.Run(); err != nil {
		utils.Perror("Unable to install ssdeep")
		complete = false
	}
	printFinalMsg(complete)
	os.Exit(0)
}

func printFinalMsg(complete bool) {
	if complete {
		fmt.Println("Installed successfully.")
		fmt.Println("All required packages have been installed.")
		fmt.Println("Please rerun the executable without install option")
	} else {
		utils.Perror("Some packages were not installed")
	}
}
