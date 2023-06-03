package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	info   *log.Logger
	error_ *log.Logger
	init_  = false
)

func InitLogger() {
	curTime := time.Now()
	logFile, err := os.OpenFile("log"+curTime.Format("20060102150405"), os.O_CREATE, 0666)

	if err != nil {
		fmt.Print("Logger initialization failed. Using stdout as the log file. ")
		fmt.Println(err)
		logFile = os.Stdout
		return
	}
	info = log.New(logFile, "Info:", log.LstdFlags)
	error_ = log.New(logFile, "Error:", log.LstdFlags)
	init_ = true
}

func Info(text string, a ...any) {
	if !init_ {
		panic("U forgor InitLogger 💀")
	}

	info.Printf(text, a)
}

func Error(text string, a ...any) {
	if !init_ {
		panic("U forgor InitLogger 💀")
	}

	error_.Printf(text, a)
}
