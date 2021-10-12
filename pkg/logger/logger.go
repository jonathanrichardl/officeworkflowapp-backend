package logger

import (
	"log"
	"os"
)

type LoggerInstance struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}

func NewLogger() *LoggerInstance {

	InfoLogger := log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger := log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return &LoggerInstance{InfoLogger: InfoLogger, WarningLogger: WarningLogger, ErrorLogger: ErrorLogger}

}
