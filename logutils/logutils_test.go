package logutils

import (
	"testing"
)

func TestTrace(t *testing.T) {
	Trace("This is a trace ", LevelTrace)
}

func TestTraceFf(t *testing.T) {
	TraceF("This is a trace %d", LevelTrace)
}

func TestInfo(t *testing.T) {
	Info("This is a info ", LevelInfo)
}

func TestInfoF(t *testing.T) {
	InfoF("This is a info %d", LevelInfo)
}

func TestWarning(t *testing.T) {
	Warning("This is a warning ", LevelWarning)
}

func TestWarningF(t *testing.T) {
	WarningF("This is a warning %d", LevelWarning)
}

func TestError(t *testing.T) {
	Error("This is a error ", LevelError)
}

func TestErrorF(t *testing.T) {
	ErrorF("This is a error %d", LevelError)
}
