// +build windows

package utils

import (
	"errors"
	"syscall"
	"unsafe"
)

var globalMemoryStatusEx *syscall.LazyProc

type memoryStatusEx struct {
	length               uint32
	memoryLoad           uint32
	totalPhys            uint64
	availPhys            uint64
	totalPageFile        uint64
	availPageFile        uint64
	totalVirtual         uint64
	availVirtual         uint64
	availExtendedVirtual uint64
}

func GetMemoryStatus() (uint32, uint64, uint64, error) {
	if globalMemoryStatusEx == nil {
		return 0, 0, 0, errors.New("加载获取内存状态失败")
	}
	var status memoryStatusEx
	status.length = uint32(unsafe.Sizeof(status))
	ret, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&status)))
	if ret == 0 {
		return 0, 0, 0, errors.New("获取内存状态失败")
	}
	return status.memoryLoad, status.totalPhys, status.availPhys, nil
}
