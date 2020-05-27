package utils

import (
	"log"
	"testing"
)

func TestGetCpuUseRate(t *testing.T) {
	rate, err := GetCpuUseRate()
	if err != nil {
		t.Error("Failed to get cpu use rate")
	}
	log.Print("GetCpuUseRate: ", rate)
	_, _ = GetCpuID()
}
