package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	info   *log.Logger
	error_ *log.Logger
}

var loggerInstance *Logger

func New() *Logger {
	if loggerInstance != nil {
		return loggerInstance
	}

	curTime := time.Now()
	logFile, err := os.OpenFile("log"+curTime.Format("20060102150405"), os.O_CREATE, 0666)

	if err != nil {
		fmt.Print("Logger initialization failed. Using stdout as the logger file. ")
		fmt.Println(err)
		logFile = os.Stdout
	}
	info := log.New(logFile, "Info:", log.LstdFlags)
	error_ := log.New(logFile, "Error:", log.LstdFlags)
	loggerInstance = &Logger{info, error_}
	return loggerInstance
}

func (lg *Logger) Info(text string, a ...any) {
	if len(a) == 0 {
		lg.info.Print(text)
	} else {
		lg.info.Printf(text, a)
	}
}

func (lg *Logger) Error(text string, a ...any) {
	if len(a) == 0 {
		lg.error_.Print(text)
	} else {
		lg.error_.Printf(text, a...)
	}
}
