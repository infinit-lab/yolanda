package utils

import (
	"encoding/json"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"strconv"
	"strings"
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

		sections := strings.Split(adapter.Ip, ".")
		if len(sections) < 4 {
			t.Fatal("Failed to Split ip")
		}
		n, err := strconv.Atoi(sections[3])
		if err != nil {
			t.Fatal("Failed to Atoi. error: ", err)
		}
		sections[3] = strconv.Itoa(n + 1)

		adapter.Ip = ""
		for i, s := range sections {
			adapter.Ip += s
			if i != len(sections)-1 {
				adapter.Ip += "."
			}
		}
		data, _ = json.Marshal(adapter)
		logutils.TraceF("%d. %s", i+1, string(data))
		/*
			err = SetAdapter(adapter)
			if err != nil {
				t.Error("Failed to SetAdapter. error: ", err)
			}
		*/
	}
}
