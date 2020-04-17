package logutils

import (
	"github.com/infinit-lab/yolanda/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type Level int

const (
	LevelError   = 0
	LevelWarning = 1
	LevelInfo    = 2
	LevelTrace   = 3
)

var (
	trace   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
)

func init() {
	log.Println("Initializing log...")
	path := config.GetString("log.path")
	if path != "" {
		logName := time.Now().Format("2006-01-02")
		logName = path + "/" + logName + ".log"
		file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil && file != nil {
			trace = log.New(io.MultiWriter(file, os.Stdout), "TRACE: ", log.Ldate|log.Ltime)
			info = log.New(io.MultiWriter(file, os.Stdout), "INFO: ", log.Ldate|log.Ltime)
			warning = log.New(io.MultiWriter(file, os.Stdout), "WARNING: ", log.Ldate|log.Ltime)
			error = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime)
			return
		}
	}
	log.Println("Failed to open file or path is empty")
	trace = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime)
	info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime)
	error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
}

func getLevel() int {
	return config.GetInt("log.level")
}

func getCaller() string {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		_, fileName := filepath.Split(file)
		return fileName + ":" + strconv.Itoa(line)
	}
	return ""
}

func logLn(level int, v ...interface{}) {
	var value []interface{}
	value = append(value, getCaller())
	value = append(value, v...)
	if getLevel() >= level {
		switch level {
		case LevelTrace:
			trace.Println(value...)
			break
		case LevelInfo:
			info.Println(value...)
			break
		case LevelWarning:
			warning.Println(value...)
			break
		case LevelError:
			error.Println(value...)
			break
		default:
			break
		}
	}
}

func logF(level int, format string, args ...interface{}) {
	var value []interface{}
	value = append(value, getCaller())
	value = append(value, args...)
	format = "%s " + format
	if getLevel() >= level {
		switch level {
		case LevelTrace:
			trace.Printf(format, value...)
			break
		case LevelInfo:
			info.Printf(format, value...)
			break
		case LevelWarning:
			warning.Printf(format, value...)
			break
		case LevelError:
			error.Printf(format, value...)
			break
		default:
			break
		}
	}
}

func Trace(v ...interface{}) {
	logLn(LevelTrace, v...)
}

func TraceF(format string, args ...interface{}) {
	logF(LevelTrace, format, args...)
}

func Info(v ...interface{}) {
	logLn(LevelInfo, v...)
}

func InfoF(format string, args ...interface{}) {
	logF(LevelInfo, format, args...)
}

func Warning(v ...interface{}) {
	logLn(LevelWarning, v...)
}

func WarningF(format string, args ...interface{}) {
	logF(LevelWarning, format, args...)
}

func Error(v ...interface{}) {
	logLn(LevelError, v...)
}

func ErrorF(format string, args ...interface{}) {
	logF(LevelError, format, args...)
}
