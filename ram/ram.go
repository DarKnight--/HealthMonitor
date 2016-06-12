package ram

/*
#include "sys/types.h"
#include "sys/sysinfo.h"
*/
import "C"

type (
	MemoryConst struct {
		TotalVirtual  int
		TotalPhysical int
	}

	MemoryStat struct {
		FreeVirtual  int
		FreePhysical int
	}
)

func (stat *MemoryStat) LoadMemoryStats() {
	var memInfo _Ctype_struct_sysinfo
	C.sysinfo(&memInfo)
	memUnit := int(memInfo.mem_unit)
	stat.FreePhysical = int(memInfo.freeram) * memUnit
	stat.FreeVirtual = int(memInfo.freeswap)*memUnit + stat.FreePhysical
}

func (consts *MemoryConst) InitMemoryConst() {
	var memInfo _Ctype_struct_sysinfo
	C.sysinfo(&memInfo)
	memUnit := int(memInfo.mem_unit)

	consts.TotalPhysical = int(memInfo.totalram) * memUnit
	consts.TotalVirtual = int(memInfo.totalswap)*memUnit + consts.TotalPhysical
}
