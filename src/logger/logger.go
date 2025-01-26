package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	logFile     *os.File
)

func init() {
	// Create logs directory
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// Create daily log file
	logFileName := fmt.Sprintf("logs/app_%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	logFile = file

	// Create multi-writer for both file and console
	infoWriter := io.MultiWriter(os.Stdout, logFile)
	errorWriter := io.MultiWriter(os.Stderr, logFile)
	debugWriter := io.MultiWriter(os.Stdout, logFile)

	// Initialize loggers with multi-writer
	infoLogger = log.New(infoWriter, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime)
	debugLogger = log.New(debugWriter, "DEBUG: ", log.Ldate|log.Ltime)
}

func getCallerInfo() string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("[%s:%d]", filepath.Base(file), line)
}

func Info(format string, v ...interface{}) {
	infoLogger.Printf("%s %s", getCallerInfo(), fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	errorLogger.Printf("%s %s", getCallerInfo(), fmt.Sprintf(format, v...))
}

func Debug(format string, v ...interface{}) {
	debugLogger.Printf("%s %s", getCallerInfo(), fmt.Sprintf(format, v...))
}

func Fatal(format string, v ...interface{}) {
	errorLogger.Fatalf("%s %s", getCallerInfo(), fmt.Sprintf(format, v...))
}

func Panic(format string, v ...interface{}) {
	errorLogger.Panicf("%s %s", getCallerInfo(), fmt.Sprintf(format, v...))
}

func Print(format string, v ...interface{}) {
	infoLogger.Printf("%s %s", getCallerInfo(), fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	errorLogger.Printf("%s %s", getCallerInfo(), fmt.Sprintf(format, v...))
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
