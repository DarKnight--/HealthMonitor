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
	CPUStat struct {
		CPUUsage float32
	}
)

func (stat *CPUStat) Init() error {
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

func (stat *CPUStat) CPUUsage() error {
	var totalUser, totalUserLow, totalSys, totalIdle, total uint64
	var percent float32
	stat.CPUUsage = percent

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
	stat.CPUUsage = percent
	return nil
}
