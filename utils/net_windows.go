package utils

/*
#cgo CFLAGS: -DWINDOWS -O3
#cgo LDFLAGS: -liphlpapi
#include "net.h"
*/
import "C"
import (
	"github.com/infinit-lab/yolanda/logutils"
	"net"
	"strings"
	"unsafe"
)

func GetNetworkInfo() ([]*Adapter, error) {
	var adapters []*Adapter
	count := C.GetAdaptersCount()
	for i := C.int(0); i < count; i++ {
		var adapter C.struct__T_ADAPTER
		ret := C.GetAdapter(i, &adapter)
		if ret == 0 {
			a := new(Adapter)
			a.Index = int(adapter.index)
			a.Name = CArrayToGoString(unsafe.Pointer(&adapter.name[0]), nameLength)
			a.Description = CArrayToGoString(unsafe.Pointer(&adapter.description[0]), nameLength)
			a.Mac = CArrayToGoString(unsafe.Pointer(&adapter.mac[0]), addrLength)
			a.Ip = CArrayToGoString(unsafe.Pointer(&adapter.ip[0]), addrLength)
			a.Mask = CArrayToGoString(unsafe.Pointer(&adapter.mask[0]), addrLength)
			a.Gateway = CArrayToGoString(unsafe.Pointer(&adapter.gateway[0]), addrLength)
			if strings.Contains(a.Description, "Virtual") {
				continue
			}
			adapters = append(adapters, a)
		}
	}
	for _, adapter := range adapters {
		i, err := net.InterfaceByIndex(adapter.Index)
		if err != nil {
			logutils.Error("Failed to InterfaceByIndex. error: ", err)
			return nil, err
		}
		adapter.Name = i.Name
	}
	return adapters, nil
}
