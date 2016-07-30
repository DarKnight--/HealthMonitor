package ram

/*
#include "sys/types.h"
#include "sys/sysinfo.h"
*/
import "C"
import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type (
	// Config holds all the necessary parameters required by the module
	Config struct {
		Profile          string
		RAMWarningLimit  int
		RecheckThreshold int
	}
	//MemoryConst holds the constant data of the memory
	MemoryConst struct {
		TotalSwap     int
		TotalPhysical int
	}
	//MemoryStat holds the current stats of the memory
	MemoryStat struct {
		FreeSwap     int
		FreePhysical int
	}
)

//LoadMemoryStats saves the current memory status in the MemoryStat struct
func (conf Config) LoadMemoryStats(stat *MemoryStat) error {
	var memInfo _Ctype_struct_sysinfo
	C.sysinfo(&memInfo)
	memUnit := int(memInfo.mem_unit)
	stat.FreeSwap = int(memInfo.freeswap) * memUnit / 1024

	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return err
	}
	defer f.Close()

	pfx := []byte("MemAvailable:")
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Bytes()
		if bytes.HasPrefix(line, pfx) {
			// len("MemAvailable:") == 13
			_, err := fmt.Sscanf(string(line[13:]), "%d", &stat.FreePhysical)
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

//InitMemoryConst saves the memory constants in the MemoryStat struct
func (conf *Config) InitMemoryConst(consts *MemoryConst) error {
	var memInfo _Ctype_struct_sysinfo
	C.sysinfo(&memInfo)
	memUnit := int(memInfo.mem_unit)
	consts.TotalSwap = int(memInfo.totalswap) * memUnit / 1024

	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return err
	}
	defer f.Close()

	pfx := []byte("MemTotal:")
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Bytes()
		if bytes.HasPrefix(line, pfx) {
			// len("MemTotal: == 10
			_, err := fmt.Sscanf(string(line[10:]), "%d", &consts.TotalPhysical)
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
