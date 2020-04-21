// +build linux

package utils

import (
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	memTotal     = "MemTotal:"
	memAvailable = "MemAvailable:"
)

func GetMemoryStatus() (uint32, uint64, uint64, error) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, err
	}
	lines := strings.Split(string(contents), "\n")

	var total, idle uint64
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == memTotal || fields[0] == memAvailable {
			val, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return 0, 0, 0, err
			}
			if fields[0] == memTotal {
				total = val * 1024
			} else {
				idle = val * 1024
				break
			}
		}
	}
	var load uint32
	if total != 0 && idle != 0 {
		load = uint32(100 * (total - idle) / total)
		return load, total, idle, nil
	}
	return 0, 0, 0, err
}
