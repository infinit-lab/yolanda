package utils

import (
	"errors"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type fileTime struct {
	lowDateTime uint32
	highDateTime uint32
}

var getSystemTimes *syscall.LazyProc

func systemTimes() (uint64, uint64, uint64, error) {
	if getSystemTimes == nil {
		return 0, 0, 0, errors.New("获取系统时间函数为空")
	}
	var idleTime, kernelTime, userTime fileTime
	ret, _, _ := getSystemTimes.Call(uintptr(unsafe.Pointer(&idleTime)), uintptr(unsafe.Pointer(&kernelTime)), uintptr(unsafe.Pointer(&userTime)))
	if ret == 0 {
		return 0, 0, 0, errors.New("获取系统时间失败")
	}
	idleTime64 := uint64(idleTime.highDateTime) << 32 | uint64(idleTime.lowDateTime)
	kernelTime64 := uint64(kernelTime.highDateTime)  << 32 | uint64(kernelTime.lowDateTime)
	userTime64 := uint64(userTime.highDateTime) << 32 | uint64(userTime.lowDateTime)
	return idleTime64, kernelTime64, userTime64, nil
}

func procStat() (idle, total uint64, err error) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					return 0, 0, err
				}
				total += val
				if i == 4 {
					idle = val
				}
			}
			return
		}
	}
	return 0, 0, errors.New("文件格式错误")
}

func GetCpuUseRate() (int, error) {
	if runtime.GOOS == "windows" {
		oldIdleTime, oldKernelTime, oldUserTime, err := systemTimes()
		if err != nil {
			return 0, err
		}
		time.Sleep(time.Second)
		idleTime, kernelTime, userTime, err := systemTimes()
		if err != nil {
			return 0, err
		}
		idleTimeDelta := idleTime - oldIdleTime
		kernelTimeDelta := kernelTime - oldKernelTime
		userTimeDelta := userTime - oldUserTime
		useRate := 100 - (idleTimeDelta * 100 / (kernelTimeDelta + userTimeDelta))
		return int(useRate), nil
	} else if runtime.GOOS == "linux" {
		oldIdle, oldTotal, err := procStat()
		if err != nil {
			return 0, err
		}
		time.Sleep(time.Second)
		idle, total, err := procStat()
		if err != nil {
			return 0, err
		}
		idleTicks := idle - oldIdle
		totalTicks := total - oldTotal
		useRate := 100 * (totalTicks - idleTicks) / totalTicks
		return int(useRate), nil
	}
	return 0, errors.New("操作系统不支持")
}
