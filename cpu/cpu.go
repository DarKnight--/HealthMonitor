package cpu

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

var (
	lastTotalUser    = uint64(0)
	lastTotalUserLow = uint64(0)
	lastTotalSys     = uint64(0)
	lastTotalIdle    = uint64(0)
)

type (
	// Config holds all the necessary parameters required by the module
	Config struct {
		Profile          string
		CPUWarningLimit  int
		RecheckThreshold int
	}
	// Stat holds the data about the % usage of cpu
	Stat struct {
		CPUUsage int
	}
)

//Init method is used to initialiase the global variables
func (conf Config) Init() error {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return err
	}
	defer f.Close()

	pfx := []byte("cpu")
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Bytes()
		if bytes.HasPrefix(line, pfx) {
			_, err := fmt.Sscanf(string(line[3:]), "%d %d %d %d", &lastTotalUser, &lastTotalUserLow,
				&lastTotalSys, &lastTotalIdle)
			if err != nil {
				return err
			}
			break
		}
	}
	if err := r.Err(); err != nil {
		return err
	}
	return nil
}

//CPUUsage sets the current cpu usage to the passed variable pointer
func (conf *Config) CPUUsage(stat *Stat) error {
	var totalUser, totalUserLow, totalSys, totalIdle, total uint64
	var percent float32
	stat.CPUUsage = int(percent)

	f, err := os.Open("/proc/stat")
	if err != nil {
		return err
	}
	defer f.Close()

	pfx := []byte("cpu")
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Bytes()
		if bytes.HasPrefix(line, pfx) {
			_, err := fmt.Sscanf(string(line[3:]), "%d %d %d %d", &totalUser, &totalUserLow,
				&totalSys, &totalIdle)
			if err != nil {
				return err
			}
			break
		}
	}
	if err := r.Err(); err != nil {
		return err
	}

	if totalUser < lastTotalUser || totalUserLow < lastTotalUserLow ||
		totalSys < lastTotalSys || totalIdle < lastTotalIdle {
		//Overflow detection. Just skip this value.
		percent = -1.0
	} else {
		total = (totalUser - lastTotalUser) + (totalUserLow - lastTotalUserLow) +
			(totalSys - lastTotalSys)
		percent = float32(total)
		total += totalIdle - lastTotalIdle
		percent /= float32(total)
		percent *= 100
	}

	lastTotalUser = totalUser
	lastTotalUserLow = totalUserLow
	lastTotalSys = totalSys
	lastTotalIdle = totalIdle
	stat.CPUUsage = int(percent)
	return nil
}
