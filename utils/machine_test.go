package utils

import (
	"log"
	"testing"
)

func TestGetMachineFingerprint(t *testing.T) {
	fingerprint, err := GetMachineFingerprint()
	if err != nil {
		t.Fatal("Failed to GetMachineFingerprint. error: ", err)
	}
	log.Println("Machine Fingerprint is ", fingerprint)
}
