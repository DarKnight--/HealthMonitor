package main

import (
	"fmt"
	"os"
	"os/exec"
)

func install() {
	if os.Getegid() != 0 {
		fmt.Println("Please run the program with sudo privileges")
		os.Exit(1)
	}
	fmt.Println("Installing ssdeep package...")
	command := exec.Command("/bin/sh", "-c", "apt-get install ssdeep")
	if err := command.Run(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Installed successfully.")
	fmt.Println("All required packages have been installed.")
	fmt.Println("Please rerun the executable without intall option")
	os.Exit(0)
}
