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
	keys := make(map[string]string)
	keys["code"] = "123"
	keys["param1"] = "3"
	err := SingleTableGet(keys, d, "`test`")
	if err != nil {
		t.Error("Failed to SingleTableGet. error: ", err)
		return
	}
	data, _ := json.Marshal(d)
	logutils.Trace(string(data))

	values, err := SingleTableGetList(keys, d, "`test`")
	if err != nil {
		t.Error("Failed to SingleTableGetList. error: ", err)
		return
	}
	data, _ = json.Marshal(values)
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

	keys := make(map[string]string)
	keys["code"] = "121"
	d = new (testData)
	err = SingleTableGet(keys, d, "`test`")
	if err != nil {
		t.Fatal("Failed to SingleTableGet. error: ", err)
	}
	data, _ := json.Marshal(d)
	logutils.Trace(string(data))
}

func TestSingleTableUpdate(t *testing.T) {
	d := new (testData)
	d.Code = "121"
	d.Param1 = 17
	d.Param2 = "Update Hello"
	d.Param3 = true
	d.Param4 = 5.17
	keys := make(map[string]string)
	keys["code"] = "121"
	err := SingleTableUpdate(keys, d, "`test`")
	if err != nil {
		t.Error("Failed to SingleTableUpdate. error: ", err)
		return
	}

	d = new (testData)
	err = SingleTableGet(keys, d, "`test`")
	if err != nil {
		t.Fatal("Failed to SingleTableGet. error: ", err)
	}
	data, _ := json.Marshal(d)
	logutils.Trace(string(data))
}

func TestSingleTableDelete(t *testing.T) {
	keys := make(map[string]string)
	keys["code"] = "121"
	err := SingleTableDelete(keys, "`test`")
	if err != nil {
		t.Error("Failed to SingleTableDelete. error: ", err)
		return
	}
	d := new (testData)
	err = SingleTableGet(keys, d, "`test`")
	if err == nil {
		t.Error("Should not be no error")
	}
	logutils.Trace(err)
}
