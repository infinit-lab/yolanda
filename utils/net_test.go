package utils

import (
	"encoding/json"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"testing"
)

func TestGetNetworkInfo(t *testing.T) {
	os.Args = append(os.Args, "log.level=3")
	config.Exec()
	adapters, err := GetNetworkInfo()
	if err != nil {
		t.Fatal("Failed to GetNetworkInfo. error: ", err)
	}
	for i, adapter := range adapters {
		data, _ := json.Marshal(adapter)
		logutils.TraceF("%d. %s", i+1, string(data))
	}
}
