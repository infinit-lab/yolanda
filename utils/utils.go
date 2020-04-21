package utils

import (
	"runtime"
	"syscall"
)

func init() {
	if runtime.GOOS == "windows" {
		kernel := syscall.NewLazyDLL("kernel32.dll")
		if kernel != nil {
			getSystemTimes = kernel.NewProc("GetSystemTimes")
			globalMemoryStatusEx = kernel.NewProc("GlobalMemoryStatusEx")
		}
	}
}
