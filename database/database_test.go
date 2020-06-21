package database

import (
	"encoding/json"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"testing"
	"time"
)

type testData struct {
	Code string `db:"code,omitupdate"`
	Param1 int `db:"param1"`
	Param2 string `db:"param2"`
	Param3 bool `db:"param3"`
	Param4 float32 `db:"param4"`
}

func TestSingleTableGet(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	os.Args = append(os.Args, "log.level=3")
	config.Exec()
	d := new (testData)
	err := SingleTableGet("123", "code", d, "`test`")
	if err != nil {
		t.Error("Failed to SingleTableGet. error: ", err)
		return
	}
	data, _ := json.Marshal(d)
	logutils.Trace(string(data))
}

func TestSingleTableCreate(t *testing.T) {
	d := new (testData)
	d.Code = "121"
	d.Param1 = 15
	d.Param2 = "Hello"
	d.Param3 = false
	d.Param4 = 3.14
	err := SingleTableCreate(d, "`test`")
	if err != nil {
		t.Error("Failed to SingleTableCreate. error: ", err)
		return
	}
}

func TestSingleTableUpdate(t *testing.T) {
	d := new (testData)
	d.Code = "121"
	d.Param1 = 17
	d.Param2 = "Update Hello"
	d.Param3 = true
	d.Param4 = 5.17
	err := SingleTableUpdate("121", "code", d, "`test`")
	if err != nil {
		t.Error("Failed to SingleTableUpdate. error: ", err)
		return
	}
}

func TestSingleTableDelete(t *testing.T) {
	err := SingleTableDelete("121", "code", "`test`")
	if err != nil {
		t.Error("Failed to SingleTableDelete. error: ", err)
		return
	}
}
