package logger

import (
	"log"
	"os"
)

type application struct {
	debugLog *log.Logger
	infoLog  *log.Logger
	errorLog *log.Logger
}

var AppLog *application

// we initialise the logger when importing this package
func init() {
	initLogger()
}

func initLogger() {
	AppLog = &application{
		debugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime),
		infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime),
	}
}
