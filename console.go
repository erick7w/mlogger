package mlogger

import (
	"fmt"
	"time"
)

type consoleMlogger struct {
	lev Level
}

func NewConsoleMlogger(lev string) mlogger {
	Level, err := parseLev(lev)
	if err != nil {
		panic(err)
	}
	return &consoleMlogger{
		lev: Level,
	}
}

func (c *consoleMlogger) outPut(lev Level, format string, a ...interface{}) {
	if lev > c.lev {
		msg := fmt.Sprintf(format, a...)
		funcName, fileName, line := getMsgCallInfo(3)
		now := time.Now().Format("2006-01-02 15:04:05")
		levStr := fmt.Sprintf(LevelColorStyle(lev), reverseParesLev(lev))
		fmt.Printf("[%s] [%s] [%s | %s | %d] %v \n", now, levStr, funcName, fileName, line, msg)
	}
}

func (c *consoleMlogger) Debug(format string, a ...interface{}) {
	c.outPut(DEBUG, format, a...)
}
func (c *consoleMlogger) Info(format string, a ...interface{}) {
	c.outPut(INFO, format, a...)
}
func (c *consoleMlogger) Warn(format string, a ...interface{}) {
	c.outPut(WARN, format, a...)
}
func (c *consoleMlogger) Error(format string, a ...interface{}) {
	c.outPut(ERROR, format, a...)
}
func (c *consoleMlogger) Fatal(format string, a ...interface{}) {
	c.outPut(FATAL, format, a...)
}
