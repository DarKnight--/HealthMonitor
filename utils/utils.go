package utils

import (
	"log"
)

func Perror(err error) {
	if err != nil {
		log.Println("[!] Error: " + err.Error())
	}
}
