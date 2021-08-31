package mlogger

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

type fileMlogger struct {
	maxSize  int
	filePath string
	lev      Level
	fileObj  *os.File
	msgChan  chan *msgData
}

type msgData struct {
	msg string
}

func NewFileMlogger(Lev string, maxSize int, filePath string) mlogger {
	level, err := parseLev(Lev)
	if err != nil {
		panic(err)
	}

	if maxSize <= 0 {
		//bytes，默认大小1MB
		maxSize = 1024 * 1024
	}

	fileFullPath := path.Join(filePath, strconv.Itoa(int(time.Now().Unix()))+".log")
	fileObj, err := os.OpenFile(fileFullPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	f := &fileMlogger{
		maxSize:  maxSize,
		lev:      level,
		filePath: filePath,
		fileObj:  fileObj,
		msgChan:  make(chan *msgData, 50000),
	}

	go f.writeToFile()
	return f
}

func (f *fileMlogger) outPut(lev Level, format string, a ...interface{}) {
	if lev > f.lev {
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName, fileName, line := getMsgCallInfo(3)
		msg := fmt.Sprintf("[%s] [%s] [%s | %s | %d] %v \n", now, reverseParesLev(lev), funcName, fileName, line, fmt.Sprintf(format, a...))
		md := &msgData{
			msg: msg,
		}
		select {
		case f.msgChan <- md:
		default:
			//防止阻塞
		}
	}
}

func (f *fileMlogger) writeToFile() {
	defer f.fileObj.Close()
	for {
		select {
		case newData := <-f.msgChan:
			//文件切割
			f.splitFile()
			//写入文件
			_, err := fmt.Fprintf(f.fileObj, newData.msg)
			if err != nil {
				fmt.Println(err)
			}

		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (f *fileMlogger) splitFile() {
	fi, err := f.fileObj.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	if fi.Size() >= int64(f.maxSize) {
		NewfileFullPath := path.Join(f.filePath, strconv.Itoa(int(time.Now().Unix()))+".log")
		NewFileObj, err := os.OpenFile(NewfileFullPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.fileObj.Close()
		f.fileObj = NewFileObj
	}
}

func (f *fileMlogger) Debug(format string, a ...interface{}) {
	f.outPut(DEBUG, format, a...)
}
func (f *fileMlogger) Info(format string, a ...interface{}) {
	f.outPut(INFO, format, a...)
}
func (f *fileMlogger) Warn(format string, a ...interface{}) {
	f.outPut(WARN, format, a...)
}
func (f *fileMlogger) Error(format string, a ...interface{}) {
	f.outPut(ERROR, format, a...)
}
func (f *fileMlogger) Fatal(format string, a ...interface{}) {
	f.outPut(FATAL, format, a...)
}
