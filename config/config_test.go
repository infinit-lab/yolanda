package config

import (
	"os"
	"testing"
)

func TestExec(t *testing.T) {
	os.Args = append(os.Args, "param1.param2=true")
	Exec()
}

func TestGetString(t *testing.T) {
	str := GetString("param1.param2")
	if str != "true" {
		t.Errorf("Failed to GetString expect true, actual %s", str)
	}
	str = GetString("param1.param3")
	if str != "123456" {
		t.Errorf("Failed to GetString expect 123456, actual %s", str)
	}
	str = GetString("param1.param4.param5")
	if str != "12345.6" {
		t.Errorf("Failed to GetString expect 12345.6, actual %s", str)
	}
}

func TestGetInt(t *testing.T) {
	i := GetInt("param1.param3")
	if i != 123456 {
		t.Errorf("Failed to GetInt expect 123456, actual %d", i)
	}
}

func TestGetFloat64(t *testing.T) {
	f := GetFloat64("param1.param4.param5")
	if f != 12345.6 {
		t.Errorf("Failed to GetFloat64 expect 12345.6, actual %f", f)
	}
}

func TestGetBool(t *testing.T) {
	b := GetBool("param1.param2")
	if b != true {
		t.Errorf("Failed to GetBool expect true, actual %t", b)
	}
}
