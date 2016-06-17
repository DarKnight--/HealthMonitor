package ram

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

//MemoryByPID return the physical memory reserved by the process with given pid
func MemoryByPID(pid int) (uint64, error) {

	f, err := os.Open(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return 0, err
	}
	defer f.Close()

	res := uint64(0)
	pfx := []byte("VmRSS:")
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Bytes()
		if bytes.HasPrefix(line, pfx) {
			_, err := fmt.Sscanf(string(line[6:]), "%d", &res)
			if err != nil {
				return 0, err
			}
			break
		}
	}
	if err := r.Err(); err != nil {
		return 0, err
	}

	return res, nil
}
