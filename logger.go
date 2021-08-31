package mlogger

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type mlogger interface {
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
}

type Level int

const (
	ANY Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

func parseLev(lev string) (Level, error) {
	levLower := strings.ToLower(lev)
	switch levLower {
	case "any":
		return ANY, nil
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warn":
		return ERROR, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		//默认LOG
		return 0, errors.New("level not exists")
	}
}

func reverseParesLev(lev Level) string {
	switch lev {
	case DEBUG:
		return "D"
	case INFO:
		return "I"
	case WARN:
		return "W"
	case ERROR:
		return "E"
	case FATAL:
		return "F"
	default:
		return "null"
	}
}

func LevelColorStyle(lev Level) string {
	switch lev {
	case DEBUG:
		return "\033[1;32;40m%s\033[0m"
	case INFO:
		return "\033[1;36;40m%s\033[0m"
	case WARN:
		return "\033[1;33;40m%s\033[0m"
	case ERROR:
		return "\033[1;31;40m%s\033[0m"
	case FATAL:
		return "\033[1;35;40m%s\033[0m"
	default:
		return "%s"
	}
}

func getMsgCallInfo(skip int) (funcName, fileName string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime.Caller failed")
		return
	}
	fullSlice := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	funcName = fullSlice[len(fullSlice)-1]
	fileName = filepath.Base(file)
	return
}
