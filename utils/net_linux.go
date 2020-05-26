package utils

/*
#cgo CFLAGS: -DLINUX -O3
#include "net.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
	"net"
	"io/ioutil"
	"strings"
	"strconv"
	"github.com/infinit-lab/yolanda/logutils"
)

func getGateway(name string) (string, error) {
	content, err := ioutil.ReadFile("/proc/net/route")
	if err != nil {
		logutils.Error("Failed to ReadFile. error: ", err)
		return "", err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		sections := strings.Split(line, "\t")
		if len(sections) < 4 {
			continue
		}
		if sections[0] == name && sections[3] == "0003" {
			gateway, err := strconv.ParseInt(sections[2], 16, 64)
			if err != nil {
				logutils.Error("Failed to ParseInt. error: ", err)
				return "", err
			}
			return fmt.Sprintf("%d.%d.%d.%d", byte(gateway), byte(gateway >> 8),
				byte(gateway >> 16), byte(gateway >> 24)), nil
		}
	}
	return "", nil
}

func GetNetworkInfo() ([]*Adapter, error) {
	var adapters []*Adapter
	count := C.GetAdaptersCount()
	for i := C.int(0); i < count; i++ {
		var adapter C.struct__T_ADAPTER
		ret := C.GetAdapter(i, &adapter)
		if ret != 0 {
			continue
		}
		a := new(Adapter)
		a.Name = CArrayToGoString(unsafe.Pointer(&adapter.name[0]), nameLength)
		a.Mac = CArrayToGoString(unsafe.Pointer(&adapter.mac[0]), addrLength)
		a.Ip = CArrayToGoString(unsafe.Pointer(&adapter.ip[0]), addrLength)
		a.Mask = CArrayToGoString(unsafe.Pointer(&adapter.mask[0]), addrLength)
		if a.Name == "lo" {
			continue
		}
		adapters = append(adapters, a)
	}
	for _, adapter := range adapters {
		i, err := net.InterfaceByName(adapter.Name)
		if err != nil {
			logutils.Error("Failed to InterfaceByName. error: ", err)
			return nil, err
		}
		adapter.Index = i.Index
		adapter.Gateway, err = getGateway(adapter.Name)
		if err != nil {
			logutils.Error("Failed to getGateway. error: ", err)
			return nil, err
		}
	}

	return adapters, nil
}
