// Package contains various utility function used in the HealthMonitor
package utils

import (
	"log"
)

/* It is used to print the errors.*/
func Perror(out string) {
	log.Println("[!] Error: " + out)
}
