package utils

import (
	"log"
	"testing"
)

func TestGetMemoryStatus(t *testing.T) {
	load, total, avail, err := GetMemoryStatus()
	if err != nil {
		t.Error("Failed to GetMemoryStatus. error: ", err)
	}
	log.Printf("MemoryLoad is %v, TotalPhys is %v, AvailPhys is %v", load, total, avail)
}
