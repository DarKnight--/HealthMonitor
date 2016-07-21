package ram

/*
#include "sys/types.h"
#include "sys/sysinfo.h"
*/
import "C"

type (
	// Config holds all the necessary parameters required by the module
	Config struct {
		Profile          string
		RAMWarningLimit  int
		RecheckThreshold int
	}
	//MemoryConst holds the constant data of the memory
	MemoryConst struct {
		TotalVirtual  int
		TotalPhysical int
	}
	//MemoryStat holds the current stats of the memory
	MemoryStat struct {
		FreeSwap     int
		FreePhysical int
	}
)

//LoadMemoryStats saves the current memory status in the MemoryStat struct
func (conf Config) LoadMemoryStats(stat *MemoryStat) {
	var memInfo _Ctype_struct_sysinfo
	C.sysinfo(&memInfo)
	memUnit := int(memInfo.mem_unit)
	stat.FreePhysical = int(memInfo.freeram) * memUnit
	stat.FreeSwap = int(memInfo.freeswap) * memUnit
}

//InitMemoryConst saves the memory constants in the MemoryStat struct
func (conf *Config) InitMemoryConst(consts *MemoryConst) {
	var memInfo _Ctype_struct_sysinfo
	C.sysinfo(&memInfo)
	memUnit := int(memInfo.mem_unit)
	consts.TotalPhysical = int(memInfo.totalram) * memUnit
	consts.TotalVirtual = int(memInfo.totalswap) * memUnit
}
