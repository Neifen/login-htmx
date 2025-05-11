package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	INFO  = createLogger("INFO")
	WARN  = createLogger("WARNING")
	ERROR = createLogger("ERROR")
)

func createLogger(level string) *log.Logger {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()

	logger := log.New(os.Stdout, (level + ": "), log.Ldate|log.Ltime)
	logger.SetOutput(io.MultiWriter(f, os.Stdout))

	return logger
}

func Info(message string) {
	INFO.Println(message)
}

func Warn(message string) {
	WARN.Println(message)
}

func ErrorMsg(message string) {
	ERROR.Println(message)
}

func Error(err error) {
	ERROR.Println(err)
}
